import React from 'react';

interface PromptComponentDisplayProps {
  title: string;
  content: string | string[] | undefined;
}

const PromptComponentDisplay: React.FC<PromptComponentDisplayProps> = ({ title, content }) => {
  if (!content) return null;

  const renderContent = () => {
    if (Array.isArray(content)) {
      return (
        <ul className="list-disc list-inside pl-4 text-slate-700">
          {content.map((item, index) => (
            <li key={index} className="mb-1">{item}</li>
          ))}
        </ul>
      );
    }
    return <pre className="whitespace-pre-wrap text-sm bg-slate-100 text-slate-800 p-3 rounded-md shadow-sm">{content}</pre>;
  };

  return (
    <div className="bg-white shadow-lg rounded-lg p-6 mb-6">
      <h3 className="text-lg font-semibold mb-2 text-slate-700">{title}</h3>
      {renderContent()}
    </div>
  );
};

export default PromptComponentDisplay; 