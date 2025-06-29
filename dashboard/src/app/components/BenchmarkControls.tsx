'use client';

import { useState, useEffect } from 'react';

interface BenchmarkControlsProps {
  runId: string;
}

interface BenchmarkStatus {
  isRunning: boolean;
}

export default function BenchmarkControls({ runId }: BenchmarkControlsProps) {
  const [isRunning, setIsRunning] = useState(false);
  const [loading, setLoading] = useState(true);

  // Check benchmark status
  const checkBenchmarkStatus = async () => {
    try {
      const response = await fetch('http://localhost:8080/v1/benchmarks/status');
      if (response.ok) {
        const data: BenchmarkStatus = await response.json();
        setIsRunning(data.isRunning);
      }
    } catch (error) {
      console.error('Error checking benchmark status:', error);
    } finally {
      setLoading(false);
    }
  };

  // Check status on component mount
  useEffect(() => {
    checkBenchmarkStatus();
  }, []);

  const handleStartBenchmark = async () => {
    try {
      setLoading(true);
      const response = await fetch(`http://localhost:8080/v1/benchmarks/create/${runId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
      });
      
      if (response.ok) {
        setIsRunning(true);
        alert('Benchmark started successfully!');
      } else {
        const error = await response.text();
        alert(`Failed to start benchmark: ${error}`);
      }
    } catch (error) {
      alert(`Error starting benchmark: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex gap-3">
      <a
        href={`/benchmarks/${runId}`}
        className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors duration-200 text-sm font-medium"
      >
        View Benchmark
      </a>
      <button
        onClick={handleStartBenchmark}
        disabled={isRunning || loading}
        className={`px-4 py-2 text-white rounded-md transition-colors duration-200 text-sm font-medium flex items-center gap-2 ${
          isRunning || loading
            ? 'bg-gray-400 cursor-not-allowed'
            : 'bg-green-600 hover:bg-green-700'
        }`}
      >
        {(isRunning || loading) && (
          <svg
            className="animate-spin h-4 w-4 text-white"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            ></circle>
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
        )}
        {isRunning ? 'Benchmark Running...' : loading ? 'Checking...' : 'Start Benchmark'}
      </button>
    </div>
  );
}