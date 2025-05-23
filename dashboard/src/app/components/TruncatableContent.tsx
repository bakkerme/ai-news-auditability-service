'use client';

import { useState } from 'react';

interface TruncatableContentProps {
  title: string;
  content: string;
  initialLineLimit?: number;
  initialCharLimit?: number;
}

const TruncatableContent: React.FC<TruncatableContentProps> = ({ title, content, initialLineLimit = 3, initialCharLimit = 300 }) => {
  const [isExpanded, setIsExpanded] = useState(false);

  // Ensure content is a string before splitting
  const safeContent = typeof content === 'string' ? content : '';
  const lines = safeContent.replace(/\r\n/g, '\n').split('\n');

  let preview: string;
  let isActuallyTruncated: boolean;

  if (lines.length > initialLineLimit) {
    preview = lines.slice(0, initialLineLimit).join('\n');
    if (preview.length > initialCharLimit) {
      preview = preview.substring(0, initialCharLimit);
    }
    isActuallyTruncated = true;
  } else {
    if (safeContent.length > initialCharLimit) {
      preview = safeContent.substring(0, initialCharLimit);
      isActuallyTruncated = true;
    } else {
      preview = safeContent;
      isActuallyTruncated = false;
    }
  }

  const displayedContent = isExpanded ? safeContent : preview;

  return (
    <div className="mb-3 p-3 border border-slate-200 dark:border-slate-700 rounded-md bg-white dark:bg-slate-800">
      <h4 className="text-md font-semibold text-slate-700 dark:text-slate-200 mb-1">{title}</h4>
      <pre className="text-sm text-slate-600 dark:text-slate-300 bg-slate-50 dark:bg-slate-700 p-2 rounded whitespace-pre-wrap font-sans">
        {displayedContent}
        {!isExpanded && isActuallyTruncated && '...'}
      </pre>
      {isActuallyTruncated && (
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="text-xs text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline mt-1 focus:outline-none"
        >
          {isExpanded ? 'Show Less' : 'Show More'}
        </button>
      )}
    </div>
  );
};

export default TruncatableContent; 