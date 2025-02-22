import { API_CONFIG } from '@/config/api';

export const resetHistory = async (): Promise<void> => {
  const response = await fetch(`${API_CONFIG.baseURL}/reset_history`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error('Failed to reset history');
  }
};

export const fullReset = async (): Promise<void> => {
  const response = await fetch(`${API_CONFIG.baseURL}/full_reset`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error('Failed to perform full reset');
  }
};
