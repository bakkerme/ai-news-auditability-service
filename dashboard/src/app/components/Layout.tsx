import React from 'react';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col bg-slate-50">
      <header className="bg-slate-800 text-white p-4 shadow-md">
        <h1 className="text-xl font-semibold">AI News Auditability - Latest Run</h1>
      </header>
      <main className="flex-grow container mx-auto p-4 md:p-6">
        {children}
      </main>
      <footer className="bg-slate-200 text-center p-4 text-sm text-slate-700 border-t border-slate-300">
        <p>Auditability Service v0.1.0</p>
      </footer>
    </div>
  );
};

export default Layout; 