import React from 'react';

interface LLMResponseDisplayProps {
  summary: string | undefined;
}

const LLMResponseDisplay: React.FC<LLMResponseDisplayProps> = ({ summary }) => {
  if (!summary) {
    return (
      <div className="bg-white shadow-lg rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-3 text-slate-700">LLM Output</h2>
        <p className="text-slate-500">No summary available.</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow-lg rounded-lg p-6">
      <h2 className="text-xl font-semibold mb-3 text-slate-700">LLM Output (Summary)</h2>
      <pre className="whitespace-pre-wrap text-sm bg-slate-100 text-slate-800 p-3 rounded-md shadow-sm">{summary}</pre>
    </div>
  );
};

export default LLMResponseDisplay; 