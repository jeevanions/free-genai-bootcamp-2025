import React, { useState, useEffect } from 'react';
import { Header } from './components/Header';
import { CategorySelect } from './components/CategorySelect';
import { WordReview } from './components/WordReview';
import { getGroups, createGroup, generateWords, importWords } from './api/client';
import type { Group, Word } from './types/api';

function App() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);
  const [words, setWords] = useState<Word[]>([]);
  const [isGenerating, setIsGenerating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  useEffect(() => {
    loadGroups();
  }, []);

  const loadGroups = async () => {
    try {
      const response = await getGroups();
      setGroups(response.items);
    } catch (e) {
      setError('Failed to load categories');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateGroup = async (name: string) => {
    try {
      const newGroup = await createGroup(name);
      setGroups([...groups, newGroup]);
      setSelectedGroup(newGroup);
    } catch (e) {
      setError('Failed to create category');
    }
  };

  const handleGenerate = async () => {
    if (!selectedGroup) return;

    setIsGenerating(true);
    setError(null);

    try {
      const response = await generateWords(selectedGroup.name);
      setWords(response.words);
    } catch (e) {
      setError('Failed to generate words');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleImport = async () => {
    if (!selectedGroup) {
      setError('Please select a category first');
      return;
    }
    try {
      setError(null);
      setSuccessMessage(null);
      const response = await importWords(selectedGroup.id, words);
      setSuccessMessage(`Successfully imported ${response.imported_count} words into category "${selectedGroup.name}"`);
      setWords([]);
      setSelectedGroup(null);
      await loadGroups();
    } catch (e) {
      setError('Failed to import words');
    }
  };

  const isValidWords = words.length > 0 && words.every(
    word => word.italian && word.english && word.parts
  );

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {error && (
          <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}
        {successMessage && (
          <div className="mb-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
            {successMessage}
          </div>
        )}

        <div className="bg-white shadow rounded-lg p-6 space-y-6">
          <div className="space-y-4">
            <h2 className="text-lg font-medium">Select Category</h2>
            <CategorySelect
              groups={groups}
              isLoading={isLoading}
              value={selectedGroup}
              onChange={setSelectedGroup}
              onCreateGroup={handleCreateGroup}
            />
            <button
              className={`w-full py-2 px-4 rounded-md font-medium ${
                selectedGroup && !isGenerating
                  ? 'bg-indigo-600 text-white hover:bg-indigo-700'
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
              }`}
              onClick={handleGenerate}
              disabled={!selectedGroup || isGenerating}
            >
              {isGenerating ? 'Generating...' : 'Generate Words'}
            </button>
          </div>

          {words.length > 0 && (
            <div className="space-y-4">
              <h2 className="text-lg font-medium">Review Words</h2>
              <WordReview
                words={words}
                onChange={setWords}
                onImport={handleImport}
                isValid={isValidWords}
              />
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

export default App;