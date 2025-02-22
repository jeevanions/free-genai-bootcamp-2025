import { API_CONFIG } from '@/config/api';

export interface Group {
  id: number;
  name: string;
  description?: string;
  word_count: number;
  words?: {
    id: number;
    italian: string;
    english: string;
    correct_count: number;
    wrong_count: number;
  }[];
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

export const getGroups = async (page: number = 1): Promise<PaginatedResponse<Group>> => {
  const params = new URLSearchParams({
    page: page.toString(),
  });

  const response = await fetch(`${API_CONFIG.baseURL}/groups?${params}`);
  if (!response.ok) {
    throw new Error('Failed to fetch groups');
  }
  return response.json();
};

export const getGroupById = async (id: number): Promise<Group> => {
  if (!id) throw new Error('Group ID is required');

  // First get the group details
  const groupResponse = await fetch(`${API_CONFIG.baseURL}/groups/${id}`);
  if (!groupResponse.ok) {
    throw new Error('Failed to fetch group');
  }
  const group = await groupResponse.json();

  // Then get the words for this group
  const wordsResponse = await fetch(`${API_CONFIG.baseURL}/groups/${id}/words`);
  if (!wordsResponse.ok) {
    throw new Error('Failed to fetch group words');
  }
  const words = await wordsResponse.json();

  return {
    ...group,
    words: words.items || [],
  };
};

export const getGroupWords = async (id: number): Promise<PaginatedResponse<Word>> => {
  const response = await fetch(`${API_CONFIG.baseURL}/groups/${id}/words`);
  if (!response.ok) {
    throw new Error('Failed to fetch group words');
  }
  return response.json();
};
