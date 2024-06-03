export interface ApiResponse {
  anonymizedText: string;
  anonymizedTextSpacy?: string; // Marked as optional since it might not always be present
}

export const fetchAnonymizedTextGPT = async (text: string): Promise<ApiResponse> => {
  const response = await fetch('http://localhost:8080/anonymize-gpt', {
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

export const fetchAnonymizedTextSpacy = async (text: string): Promise<ApiResponse> => {
  const response = await fetch('http://localhost:8080/anonymize-spacy', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ text }),
  });

  if (!response.ok) {
    throw new Error('API request failed');
  }

  const data = await response.json();
  console.log('Response data:', data);
  return { anonymizedText: '', anonymizedTextSpacy: data.anonymizedText }; 
};
