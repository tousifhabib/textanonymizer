export interface ApiResponse {
    anonymizedText: string;
  }
  
  export const fetchAnonymizedText = async (text: string): Promise<ApiResponse> => {
    const response = await fetch('http://localhost:8080/anonymize', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ text }),
    });
  
    if (!response.ok) {
      throw new Error('API request failed');
    }
  
    return response.json();
  };
  