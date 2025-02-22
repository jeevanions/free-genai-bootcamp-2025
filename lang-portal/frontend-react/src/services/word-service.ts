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

export interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    current_page: number;
    total_pages: number;
    total_items: number;
    per_page: number;
  };
}

export const getWords = async (page: number = 1, search?: string): Promise<PaginatedResponse<Word>> => {
  const params = new URLSearchParams({
    page: page.toString(),
  });
  
  // Only add search param if it's not empty
  if (search && search.trim() !== '') {
    params.append('search', search.trim());
  }

  const response = await fetch(`${API_CONFIG.baseURL}/words?${params}`);
  if (!response.ok) {
    throw new Error('Failed to fetch words');
  }
  return response.json();
};

export const getWordById = async (id: number): Promise<Word> => {
  const response = await fetch(`${API_CONFIG.baseURL}/words/${id}`);
  if (!response.ok) {
    throw new Error('Failed to fetch word');
  }
  return response.json();
};
