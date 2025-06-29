import Layout from '../../components/Layout';
import MetadataCard from '../../components/MetadataCard';
import AccordionItem from '../../components/AccordionItem';
import TruncatableContent from '../../components/TruncatableContent';
import { BenchmarkResults, RunData } from '../../types';

async function getBenchmarkData(runId: string): Promise<BenchmarkResults> {
  try {
    const res = await fetch(`http://localhost:8080/v1/benchmarks/${runId}`, {
      cache: 'no-store',
    });

    if (!res.ok) {
      const errorBody = await res.text();
      console.error(`Error fetching benchmark data: ${res.status} ${res.statusText} - ${errorBody}`);
      throw new Error(`Failed to fetch benchmark data: ${res.status} ${res.statusText}. Body: ${errorBody}`);
    }
    const data: BenchmarkResults = await res.json();
    console.log('Benchmark data:', data);
    return data;
  } catch (error) {
    console.error('Exception during fetch benchmark data:', error);
    if (error instanceof Error && error.message.startsWith('Failed to fetch')) {
      throw error;
    }
    throw new Error(`An unexpected error occurred while fetching benchmark data: ${error instanceof Error ? error.message : String(error)}`);
  }
}

async function getRunData(runId: string): Promise<RunData> {
  try {
    const res = await fetch(`http://localhost:8080/v1/runs/${runId}`, {
      cache: 'no-store',
    });

    if (!res.ok) {
      const errorBody = await res.text();
      console.error(`Error fetching run data: ${res.status} ${res.statusText} - ${errorBody}`);
      throw new Error(`Failed to fetch run data: ${res.status} ${res.statusText}. Body: ${errorBody}`);
    }
    const data: RunData = await res.json();
    console.log('Run data:', data);
    return data;
  } catch (error) {
    console.error('Exception during fetch run data:', error);
    if (error instanceof Error && error.message.startsWith('Failed to fetch')) {
      throw error;
    }
    throw new Error(`An unexpected error occurred while fetching run data: ${error instanceof Error ? error.message : String(error)}`);
  }
}

function getQualityColor(rating: string): string {
  switch (rating) {
    case 'Excellent': return 'text-green-600 dark:text-green-400 bg-green-50 dark:bg-green-900/20';
    case 'Good': return 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20';
    case 'Fair': return 'text-yellow-600 dark:text-yellow-400 bg-yellow-50 dark:bg-yellow-900/20';
    case 'Poor': return 'text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20';
    default: return 'text-slate-600 dark:text-slate-400 bg-slate-50 dark:bg-slate-700';
  }
}

function getRelevanceColor(correct: boolean): string {
  return correct 
    ? 'text-green-600 dark:text-green-400 bg-green-50 dark:bg-green-900/20'
    : 'text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20';
}

interface BenchmarkPageProps {
  params: Promise<{
    runId: string;
  }>;
}

export default async function BenchmarkPage({ params }: BenchmarkPageProps) {
  const { runId } = await params;
  let benchmarkData: BenchmarkResults | null = null;
  let runData: RunData | null = null;
  let fetchError: string | null = null;

  try {
    [benchmarkData, runData] = await Promise.all([
      getBenchmarkData(runId),
      getRunData(runId)
    ]);
  } catch (error) {
    console.error("Error in BenchmarkPage while fetching data:", error);
    fetchError = error instanceof Error ? error.message : "An unknown error occurred.";
  }

  if (fetchError || !benchmarkData || !runData) {
    return (
      <Layout>
        <div className="container mx-auto p-4">
          <h1 className="text-2xl font-bold text-red-600 dark:text-red-400 mb-4">Error Loading Benchmark Data</h1>
          <p className="text-red-500 dark:text-red-300 text-center text-xl">
            Failed to load benchmark data for run ID: {runId}. Please check if the benchmark exists or try again later.
          </p>
          {fetchError && (
            <pre className="mt-4 p-4 bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200 border border-red-300 dark:border-red-700 rounded-md whitespace-pre-wrap">
              {fetchError}
            </pre>
          )}
        </div>
      </Layout>
    );
  }

  if (benchmarkData.failureReason) {
    return (
      <Layout>
        <div className="container mx-auto p-4">
          <h1 className="text-2xl font-bold text-red-600 dark:text-red-400 mb-4">Benchmark Failed</h1>
          <p className="text-red-500 dark:text-red-300 text-center text-xl mb-4">
            The benchmark for run ID {runId} failed to complete.
          </p>
          <div className="p-4 bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200 border border-red-300 dark:border-red-700 rounded-md">
            <h3 className="font-semibold mb-2">Failure Reason:</h3>
            <p>{benchmarkData.failureReason}</p>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <MetadataCard
        runId={runId}
        timestamp={runData.runDate}
        overallModel={runData.overallModelUsed}
        imageModel={runData.imageModelUsed}
        webContentModel={runData.webContentModelUsed}
        totalProcessingTime={runData.totalProcessingTime}
      />

      {/* Benchmark Results Summary */}
      <div className="mb-6 p-6 bg-white dark:bg-slate-800 shadow-lg rounded-lg border border-slate-200 dark:border-slate-700">
        <h2 className="text-2xl font-semibold mb-4 text-slate-800 dark:text-slate-100 border-b dark:border-slate-600 pb-2">Benchmark Results</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
          <div className="p-4 bg-slate-50 dark:bg-slate-700 rounded-lg">
            <h3 className="text-sm font-medium text-slate-600 dark:text-slate-400">Total Items</h3>
            <p className="text-2xl font-bold text-slate-800 dark:text-slate-200">{benchmarkData.totalItems}</p>
          </div>
          <div className="p-4 bg-slate-50 dark:bg-slate-700 rounded-lg">
            <h3 className="text-sm font-medium text-slate-600 dark:text-slate-400">Quality Score</h3>
            <p className="text-2xl font-bold text-slate-800 dark:text-slate-200">{benchmarkData.qualityScore.toFixed(1)}%</p>
          </div>
          <div className="p-4 bg-slate-50 dark:bg-slate-700 rounded-lg">
            <h3 className="text-sm font-medium text-slate-600 dark:text-slate-400">Relevance Accuracy</h3>
            <p className="text-2xl font-bold text-slate-800 dark:text-slate-200">{(benchmarkData.relevanceAccuracy * 100).toFixed(1)}%</p>
          </div>
        </div>
        
        <div className="space-y-4">
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Persona</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">{benchmarkData.personaName}</p>
          </div>
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Focus Areas</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">
              {benchmarkData.personaFocusAreas.join(', ')}
            </p>
          </div>
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Completed</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">
              {new Date(benchmarkData.timestamp).toLocaleString()}
            </p>
          </div>
        </div>
      </div>

      {/* Detailed Evaluations */}
      <div className="mb-6">
        <h2 className="text-2xl font-semibold mb-4 text-slate-800 dark:text-slate-100">Detailed Evaluations</h2>
        {Object.keys(benchmarkData.detailedEvaluations).length > 0 ? (
          Object.entries(benchmarkData.detailedEvaluations).map(([itemId, evaluation]) => {
            // Find the corresponding entry from run data
            const entry = runData.entrySummaries.find(e => e.results.id === itemId);
            const title = entry?.results.title || `Item ${itemId}`;
            
            return (
              <AccordionItem key={itemId} title={title}>
                <div className="space-y-4 p-2">
                  {/* Item Details */}
                  {entry && (
                    <div className="p-3 border border-slate-200 dark:border-slate-700 rounded-md bg-white dark:bg-slate-800 mb-3">
                      <h4 className="text-md font-semibold text-slate-700 dark:text-slate-300 mb-2">Item Details</h4>
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-2 text-sm">
                        <div className="bg-slate-50 dark:bg-slate-700 p-2 rounded text-slate-800 dark:text-slate-200">
                          <span className="font-medium text-slate-600 dark:text-slate-400">ID:</span> {itemId}
                        </div>
                        <div className="bg-slate-50 dark:bg-slate-700 p-2 rounded text-slate-800 dark:text-slate-200">
                          <span className="font-medium text-slate-600 dark:text-slate-400">Relevant:</span> {entry.results.isRelevant ? 'Yes' : 'No'}
                        </div>
                        <div className="bg-slate-50 dark:bg-slate-700 p-2 rounded text-slate-800 dark:text-slate-200">
                          <span className="font-medium text-slate-600 dark:text-slate-400">Processing Time:</span> {entry.processingTimeMs} ms
                        </div>
                        <div className="bg-slate-50 dark:bg-slate-700 p-2 rounded text-slate-800 dark:text-slate-200 col-span-full">
                          <span className="font-medium text-slate-600 dark:text-slate-400">Link:</span> {entry.results.link || "N/A"}
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Evaluation Results */}
                  <div className="p-4 border border-slate-200 dark:border-slate-700 rounded-md bg-slate-50 dark:bg-slate-800">
                    <h4 className="text-md font-semibold text-slate-700 dark:text-slate-300 mb-3">Evaluation Results</h4>
                    
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                      <div className={`p-3 rounded-md border ${getQualityColor(evaluation.qualityRating)}`}>
                        <h5 className="font-medium mb-1">Quality Rating</h5>
                        <p className="text-lg font-bold">{evaluation.qualityRating}</p>
                      </div>
                      <div className={`p-3 rounded-md border ${getRelevanceColor(evaluation.relevanceCorrect)}`}>
                        <h5 className="font-medium mb-1">Relevance Assessment</h5>
                        <p className="text-lg font-bold">{evaluation.relevanceCorrect ? 'Correct' : 'Incorrect'}</p>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <TruncatableContent title="Quality Explanation" content={evaluation.qualityExplanation} initialLineLimit={3} />
                      <TruncatableContent title="Relevance Explanation" content={evaluation.relevanceExplanation} initialLineLimit={3} />
                    </div>
                  </div>

                  {/* Original Content */}
                  {entry && (
                    <div className="space-y-3">
                      <TruncatableContent title="Original Summary" content={entry.results.summary || ""} initialLineLimit={3} />
                      {entry.results.commentSummary && (
                        <TruncatableContent title="Comment Summary" content={entry.results.commentSummary} initialLineLimit={3} />
                      )}
                      {entry.results.relevanceToCriteria && (
                        <TruncatableContent title="Relevance to Criteria" content={entry.results.relevanceToCriteria} initialLineLimit={3} />
                      )}
                      {entry.results.imageDescription && (
                        <TruncatableContent title="Image Description" content={entry.results.imageDescription} initialLineLimit={3} />
                      )}
                      {entry.results.webContentSummary && (
                        <TruncatableContent title="Web Content Summary" content={entry.results.webContentSummary} initialLineLimit={3} />
                      )}
                    </div>
                  )}
                </div>
              </AccordionItem>
            );
          })
        ) : (
          <p className="text-slate-600 dark:text-slate-400">No detailed evaluations available for this benchmark.</p>
        )}
      </div>
    </Layout>
  );
}
