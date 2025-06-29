package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bakkerme/ai-news-auditability-service/internal/models"
	"github.com/dgraph-io/badger/v3"
)

var db *badger.DB
var defaultRunDataTTL time.Duration

const (
	dbPathPrefix   = "badger"
	runDataDir     = "rundata"
	benchmarkDir   = "benchmarks"
)

// InitDB initializes the BadgerDB database.
// It creates the database directory if it doesn't exist.
// It also sets up a goroutine for garbage collection.
func InitDB(basePath string, runDataTTLHours int) error {
	dbDir := filepath.Join(basePath, dbPathPrefix)
	if err := os.MkdirAll(dbDir, 0777); err != nil {
		return fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
	}

	opts := badger.DefaultOptions(dbDir)
	opts.Logger = nil // Disable Badger's default logger to avoid noise; we'll log errors ourselves.

	var err error
	db, err = badger.Open(opts)
	if err != nil {
		return fmt.Errorf("failed to open badger database: %w", err)
	}

	if runDataTTLHours > 0 {
		defaultRunDataTTL = time.Duration(runDataTTLHours) * time.Hour
		log.Printf("Run data TTL set to %v", defaultRunDataTTL)
	} else {
		log.Println("Run data TTL not set (or set to zero/negative), entries will not expire by default TTL.")
	}

	// Run GC periodically
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			} else if err != badger.ErrNoRewrite {
				log.Printf("Error during BadgerDB GC: %v", err)
			}
		}
	}()

	log.Println("BadgerDB initialized successfully at", dbDir)
	return nil
}

// CloseDB closes the BadgerDB database.
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing BadgerDB: %v", err)
		} else {
			log.Println("BadgerDB closed successfully.")
		}
	}
}

// getRunDBKeyPrefix returns the path for storing run data.
// For now, we store all runs under a common prefix.
// We could extend this to use subdirectories if needed (e.g., by date).
func getRunDBKeyPrefix() []byte {
	return []byte(runDataDir + "/")
}

// SaveRunData saves the provided RunData to BadgerDB.
// The runID is used as the key.
func SaveRunData(runID string, data models.PersistedRunData) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	key := []byte(filepath.Join(runDataDir, runID))

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal run data to JSON: %w", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, jsonData)
		if defaultRunDataTTL > 0 {
			entry = entry.WithTTL(defaultRunDataTTL)
		}
		return txn.SetEntry(entry)
	})
	if err != nil {
		return fmt.Errorf("failed to save run data (ID: %s) to BadgerDB: %w", runID, err)
	}
	log.Printf("Successfully saved run data with ID: %s", runID)
	return nil
}

// GetRunData retrieves RunData from BadgerDB by its ID.
func GetRunData(runID string) (*models.PersistedRunData, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	key := []byte(filepath.Join(runDataDir, runID))
	var runData models.PersistedRunData

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("run data with ID '%s' not found: %w", runID, err)
			}
			return fmt.Errorf("failed to get run data (ID: %s) from BadgerDB: %w", runID, err)
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("failed to copy value for run data (ID: %s): %w", runID, err)
		}

		if err := json.Unmarshal(val, &runData); err != nil {
			return fmt.Errorf("failed to unmarshal run data (ID: %s) from JSON: %w", runID, err)
		}
		return nil
	})

	if err != nil {
		return nil, err // Error already contains specific details
	}
	return &runData, nil
}

// ListRunMetadata retrieves a list of RunMetadata from BadgerDB.
// This is a basic implementation that iterates over keys with the runDataDir prefix.
// For production, consider pagination and more efficient querying/indexing if performance becomes an issue.
func ListRunMetadata(limit int) ([]models.RunMetadata, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var runs []models.RunMetadata
	keyPrefix := getRunDBKeyPrefix()

	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		count := 0
		for it.Seek(keyPrefix); it.ValidForPrefix(keyPrefix); it.Next() {
			if limit > 0 && count >= limit {
				break
			}
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var runData models.PersistedRunData
				// We need to unmarshal the full RunData to get metadata.
				// If performance is critical and RunData objects are large,
				// consider storing metadata separately or in a more queryable format.
				if err := json.Unmarshal(val, &runData); err != nil {
					// Log error but try to continue if possible, or return error to stop.
					log.Printf("error unmarshalling RunData for key %s: %v", string(item.Key()), err)
					return nil // Skip this item
				}

				// Construct RunMetadata from RunData
				// The RunID for metadata should be the key used in the DB (without prefix)
				dbKey := string(item.Key())
				runIDInDB := filepath.Base(dbKey)

				meta := models.RunMetadata{
					ID:          runIDInDB, // Use the actual DB key part
					RunDate:     runData.RunDate,
					PersonaName: runData.Persona.Name,
					// Model field was removed from RunMetadata, so not setting it.
				}
				runs = append(runs, meta)
				count++
				return nil
			})
			if err != nil {
				return fmt.Errorf("error processing item value: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list run metadata from BadgerDB: %w", err)
	}
	return runs, nil
}

// DeleteRunData deletes RunData from BadgerDB by its ID.
func DeleteRunData(runID string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	key := []byte(filepath.Join(runDataDir, runID))

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err == badger.ErrKeyNotFound {
			// Consider whether to return an error or not if key doesn't exist
			log.Printf("Attempted to delete non-existent run data with ID: %s", runID)
			return nil // Or return an error indicating not found
		}
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to delete run data (ID: %s) from BadgerDB: %w", runID, err)
	}
	log.Printf("Successfully deleted run data with ID (if it existed): %s", runID)
	return nil
}

// SaveBenchmarkResults saves benchmark results to BadgerDB
func SaveBenchmarkResults(benchmarkID string, results models.BenchmarkResults) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	key := []byte(fmt.Sprintf("%s/%s", benchmarkDir, benchmarkID))

	jsonData, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal benchmark results to JSON: %w", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, jsonData)
		if defaultRunDataTTL > 0 {
			entry = entry.WithTTL(defaultRunDataTTL)
		}
		return txn.SetEntry(entry)
	})
	if err != nil {
		return fmt.Errorf("failed to save benchmark results (ID: %s) to BadgerDB: %w", benchmarkID, err)
	}
	log.Printf("Successfully saved benchmark results with ID: %s", benchmarkID)
	return nil
}

// GetBenchmarkResults retrieves benchmark results from BadgerDB by run ID
// Note: This searches for benchmark results associated with a run ID
func GetBenchmarkResults(runID string) (*models.BenchmarkResults, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	keyPrefix := []byte(fmt.Sprintf("%s/", benchmarkDir))
	var results *models.BenchmarkResults

	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(keyPrefix); it.ValidForPrefix(keyPrefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var benchmarkResults models.BenchmarkResults
				if err := json.Unmarshal(val, &benchmarkResults); err != nil {
					log.Printf("error unmarshalling benchmark results for key %s: %v", string(item.Key()), err)
					return nil // Skip this item
				}

				// Check if this benchmark is for the requested run ID
				if benchmarkResults.RunID == runID {
					results = &benchmarkResults
					return nil
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("error processing benchmark item value: %w", err)
			}
			
			// If we found results, break out of the loop
			if results != nil {
				break
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get benchmark results for run ID %s: %w", runID, err)
	}

	if results == nil {
		return nil, fmt.Errorf("benchmark results for run ID '%s' not found", runID)
	}

	return results, nil
}

// GetBenchmarkResultsByBenchmarkID retrieves benchmark results by benchmark ID
func GetBenchmarkResultsByBenchmarkID(benchmarkID string) (*models.BenchmarkResults, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	key := []byte(fmt.Sprintf("%s/%s", benchmarkDir, benchmarkID))
	var results models.BenchmarkResults

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("benchmark results with ID '%s' not found: %w", benchmarkID, err)
			}
			return fmt.Errorf("failed to get benchmark results (ID: %s) from BadgerDB: %w", benchmarkID, err)
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("failed to copy value for benchmark results (ID: %s): %w", benchmarkID, err)
		}

		if err := json.Unmarshal(val, &results); err != nil {
			return fmt.Errorf("failed to unmarshal benchmark results (ID: %s) from JSON: %w", benchmarkID, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &results, nil
}
