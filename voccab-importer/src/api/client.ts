import axios from 'axios';
import type { GroupsResponse, GenerateWordsRequest, GenerateWordsResponse, ImportWordsResponse, Group, Word } from '../types/api';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getGroups = async (): Promise<GroupsResponse> => {
  const response = await api.get('/groups');
  return response.data;
};

export const createGroup = async (name: string): Promise<Group> => {
  const response = await api.post('/groups', { name });
  return response.data;
};

export const generateWords = async (category: string): Promise<GenerateWordsResponse> => {
  const response = await api.post('/words/llm/generate-words', { category } as GenerateWordsRequest);
  return response.data;
};

export const importWords = async (groupId: number, words: Word[]): Promise<ImportWordsResponse> => {
  const response = await api.post('/words/import', { group_id: groupId, words });
  return response.data;
};