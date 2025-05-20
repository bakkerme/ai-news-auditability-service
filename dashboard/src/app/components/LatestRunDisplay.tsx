'use client';

import { useEffect, useState } from 'react';
import { getLatestRun, RunMetadata } from '../lib/api';

export default function LatestRunDisplay() {
  const [run, setRun] = useState<RunMetadata | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true);
        const data = await getLatestRun();
        setRun(data);
        setError(null);
      } catch (err: any) {
        setError(err.message || 'Failed to load latest run data.');
        setRun(null);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  if (loading) {
    return <div className="p-4"><p className="text-gray-500">Loading latest run data...</p></div>;
  }

  if (error) {
    return (
      <div className="p-4 bg-red-100 border border-red-400 text-red-700 rounded">
        <p className="font-bold">Error:</p>
        <p>{error}</p>
      </div>
    );
  }

  if (!run) {
    return <div className="p-4"><p className="text-gray-500">No run data available.</p></div>;
  }

  return (
    <div className="p-4 border border-gray-300 rounded-lg shadow-sm bg-white">
      <h2 className="text-xl font-semibold mb-3 text-gray-700">Latest Run Details</h2>
      <div className="space-y-2">
        <p><strong className="text-gray-600">ID:</strong> {run.id}</p>
        <p><strong className="text-gray-600">Persona:</strong> {run.personaName}</p>
        <p><strong className="text-gray-600">Run Date:</strong> {new Date(run.runDate).toLocaleString()}</p>
        {run.overallModelUsed && <p><strong className="text-gray-600">Model Used:</strong> {run.overallModelUsed}</p>}
        {typeof run.totalItems === 'number' && <p><strong className="text-gray-600">Total Items:</strong> {run.totalItems}</p>}
        {typeof run.hasBenchmark === 'boolean' && (
          <p><strong className="text-gray-600">Benchmarked:</strong> {run.hasBenchmark ? 'Yes' : 'No'}</p>
        )}
      </div>
    </div>
  );
} 