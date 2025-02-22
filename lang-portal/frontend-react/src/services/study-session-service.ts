import { API_CONFIG } from '@/config/api';
import { StudySession, PaginatedResponse } from './api-types';

export const getStudySessions = async (page: number = 1): Promise<PaginatedResponse<StudySession>> => {
  const response = await fetch(`${API_CONFIG.baseURL}/study_sessions?page=${page}`);
  if (!response.ok) {
    throw new Error('Failed to fetch study sessions');
  }
  return response.json();
};
