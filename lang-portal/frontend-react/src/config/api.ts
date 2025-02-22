export const API_CONFIG = {
    baseURL: 'http://localhost:8080/api',
};

export const getApiUrl = (endpoint: string): string => {
    return `${API_CONFIG.baseURL}${endpoint}`;
};

export const getDashboardLastStudySession = async () => {
    const response = await fetch(getApiUrl('/dashboard/last_study_session'));
    return response.json();
};

export const getDashboardQuickStats = async () => {
    const response = await fetch(getApiUrl('/dashboard/quick-stats'));
    return response.json();
};

export const getDashboardStudyProgress = async () => {
    const response = await fetch(getApiUrl('/dashboard/study_progress'));
    return response.json();
};
