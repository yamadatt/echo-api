const API_BASE_URL = 'https://ujxtgteuma.execute-api.ap-northeast-1.amazonaws.com/prod';

export interface ApiResponse {
  [key: string]: any;
}

export const apiService = {
  async get(): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`GET request failed: ${response.status}`);
    }

    return response.json();
  },

  async post(data: any = {}): Promise<ApiResponse> {
    const response = await fetch(`${API_BASE_URL}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`POST request failed: ${response.status}`);
    }

    return response.json();
  },
};