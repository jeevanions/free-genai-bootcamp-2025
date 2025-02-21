import { API_CONFIG } from '@/config/api';

export interface Word {
  id: number;
  italian: string;
  english: string;
  parts: {
    type: string;
    gender?: string;
    plural?: string;
    conjugation?: string;
    irregular?: boolean;
    usage?: string[];
  };
  correct_count: number;
  wrong_count: number;
}

export const getWordById = async (id: number): Promise<Word> => {
  const response = await fetch(`${API_CONFIG.baseURL}/api/words/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch word');
  }
  return response.json();
};
