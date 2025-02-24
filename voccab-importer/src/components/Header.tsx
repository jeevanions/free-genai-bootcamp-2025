import React from 'react';
import { Book } from 'lucide-react';

export const Header: React.FC = () => {
  return (
    <header className="bg-indigo-600 shadow-lg">
      <div className="max-w-7xl mx-auto px-4 py-6 sm:px-6 lg:px-8 flex items-center">
        <Book className="h-8 w-8 text-white mr-3" />
        <h1 className="text-3xl font-bold tracking-tight text-white">
          Vocabular Importer
        </h1>
      </div>
    </header>
  );
};