import React from 'react';
import type { Word } from '../types/api';

interface WordReviewProps {
  words: Word[];
  onChange: (words: Word[]) => void;
  onImport: () => void;
  isValid: boolean;
}

export const WordReview: React.FC<WordReviewProps> = ({
  words,
  onChange,
  onImport,
  isValid,
}) => {
  const handleTextChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    try {
      const parsed = JSON.parse(event.target.value);
      onChange(parsed);
    } catch (e) {
      // Invalid JSON - parent component will handle this via isValid prop
    }
  };

  return (
    <div className="space-y-4">
      <textarea
        className="w-full h-96 font-mono text-sm p-4 border rounded-lg"
        value={JSON.stringify(words, null, 2)}
        onChange={handleTextChange}
      />
      <button
        className={`w-full py-2 px-4 rounded-md font-medium ${
          isValid
            ? 'bg-green-600 text-white hover:bg-green-700'
            : 'bg-gray-300 text-gray-500 cursor-not-allowed'
        }`}
        onClick={onImport}
        disabled={!isValid}
      >
        Import Words
      </button>
    </div>
  );
};