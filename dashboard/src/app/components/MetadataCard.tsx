import React from 'react';

interface MetadataCardProps {
  runId: string;
  timestamp: string;
  overallModel?: string;
  imageModel?: string;
  webContentModel?: string;
  totalProcessingTime?: number;
  entryProcessingTime?: number;
}

const MetadataCard: React.FC<MetadataCardProps> = ({
  runId,
  timestamp,
  overallModel,
  imageModel,
  webContentModel,
  totalProcessingTime,
  entryProcessingTime,
}) => {
  return (
    <div className="bg-white dark:bg-slate-800 shadow-lg rounded-lg p-6 mb-6 border border-slate-200 dark:border-slate-700">
      <h2 className="text-xl font-semibold mb-3 text-slate-700 dark:text-slate-200">Run Metadata</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-x-4 gap-y-2 text-sm">
        <p><strong className="font-medium text-slate-600 dark:text-slate-400">Run ID:</strong> <span className="text-slate-800 dark:text-slate-300">{runId}</span></p>
        <p><strong className="font-medium text-slate-600 dark:text-slate-400">Timestamp:</strong> <span className="text-slate-800 dark:text-slate-300">{new Date(timestamp).toLocaleString()}</span></p>
        {overallModel && <p><strong className="font-medium text-slate-600 dark:text-slate-400">Overall Model:</strong> <span className="text-slate-800 dark:text-slate-300">{overallModel}</span></p>}
        {imageModel && <p><strong className="font-medium text-slate-600 dark:text-slate-400">Image Model:</strong> <span className="text-slate-800 dark:text-slate-300">{imageModel}</span></p>}
        {webContentModel && <p><strong className="font-medium text-slate-600 dark:text-slate-400">Web Content Model:</strong> <span className="text-slate-800 dark:text-slate-300">{webContentModel}</span></p>}
        {totalProcessingTime !== undefined && <p><strong className="font-medium text-slate-600 dark:text-slate-400">Total Processing Time:</strong> <span className="text-slate-800 dark:text-slate-300">{totalProcessingTime} ms</span></p>}
        {entryProcessingTime !== undefined && <p><strong className="font-medium text-slate-600 dark:text-slate-400">Primary Entry Processing Time:</strong> <span className="text-slate-800 dark:text-slate-300">{entryProcessingTime} ms</span></p>}
      </div>
    </div>
  );
};

export default MetadataCard; 