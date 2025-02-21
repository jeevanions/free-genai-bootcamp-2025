import { API_CONFIG } from '../config/api';

class ApiService {
    private baseURL: string;

    constructor() {
        this.baseURL = API_CONFIG.baseURL;
    }

    async get(endpoint: string) {
        const response = await fetch(`${this.baseURL}${endpoint}`);
        return response.json();
    }

    async post(endpoint: string, data: any) {
        const response = await fetch(`${this.baseURL}${endpoint}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        return response.json();
    }
}

export const apiService = new ApiService();
