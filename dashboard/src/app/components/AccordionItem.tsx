'use client';

import { useState } from 'react';

interface AccordionItemProps {
  title: string;
  children: React.ReactNode;
}

const AccordionItem: React.FC<AccordionItemProps> = ({ title, children }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="border border-slate-300 dark:border-slate-700 rounded-md mb-2">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full flex justify-between items-center p-4 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 focus:outline-none"
      >
        <span className="font-semibold text-slate-700 dark:text-slate-200">{title}</span>
        <span className="text-slate-700 dark:text-slate-200">{isOpen ? '▲' : '▼'}</span>
      </button>
      {isOpen && <div className="p-4 bg-white dark:bg-slate-800">{children}</div>}
    </div>
  );
};

export default AccordionItem; 