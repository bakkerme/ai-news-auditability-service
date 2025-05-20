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
    <div className="bg-white shadow-lg rounded-lg p-6 mb-6">
      <h2 className="text-xl font-semibold mb-3 text-slate-700">Run Metadata</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-x-4 gap-y-2 text-sm">
        <p><strong className="font-medium text-slate-600">Run ID:</strong> <span className="text-slate-800">{runId}</span></p>
        <p><strong className="font-medium text-slate-600">Timestamp:</strong> <span className="text-slate-800">{new Date(timestamp).toLocaleString()}</span></p>
        {overallModel && <p><strong className="font-medium text-slate-600">Overall Model:</strong> <span className="text-slate-800">{overallModel}</span></p>}
        {imageModel && <p><strong className="font-medium text-slate-600">Image Model:</strong> <span className="text-slate-800">{imageModel}</span></p>}
        {webContentModel && <p><strong className="font-medium text-slate-600">Web Content Model:</strong> <span className="text-slate-800">{webContentModel}</span></p>}
        {totalProcessingTime !== undefined && <p><strong className="font-medium text-slate-600">Total Processing Time:</strong> <span className="text-slate-800">{totalProcessingTime} ms</span></p>}
        {entryProcessingTime !== undefined && <p><strong className="font-medium text-slate-600">Primary Entry Processing Time:</strong> <span className="text-slate-800">{entryProcessingTime} ms</span></p>}
      </div>
    </div>
  );
};

export default MetadataCard; 