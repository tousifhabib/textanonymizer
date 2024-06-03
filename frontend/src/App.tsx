import React, { useState, ChangeEvent, FC, useCallback } from 'react';
import './App.css';
import { fetchAnonymizedTextGPT, fetchAnonymizedTextSpacy } from './api';
import TextInput from './components/TextInput';
import AnonymizeButton from './components/AnonymizeButton';
import ErrorMessage from './components/ErrorMessage';

const App: FC = () => {
  const [text, setText] = useState('');
  const [anonymizedText, setAnonymizedText] = useState('');
  const [anonymizedTextSpacy, setAnonymizedTextSpacy] = useState('');
  const [error, setError] = useState<string | null>(null);

  const handleTextChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    setText(e.target.value);
  };

  const handleAnonymize = useCallback(async () => {
    setError(null);
    try {
      const data = await fetchAnonymizedTextGPT(text);
      setAnonymizedText(data.anonymizedText);
    } catch (error) {
      console.error('Error:', error);
      setError('Failed to anonymize the text. Please try again.');
    }
  }, [text]);

  const handleAnonymizeSpacy = useCallback(async () => {
    setError(null);
    try {
      const data = await fetchAnonymizedTextSpacy(text);
      console.log('Data received:', data);
      setAnonymizedTextSpacy(data.anonymizedTextSpacy || ''); // Handle the new key
    } catch (error) {
      console.error('Error:', error);
      setError('Failed to anonymize the text using spaCy. Please try again.');
    }
  }, [text]);

  return (
    <div className="App">
      <header className="App-header">
        <h1>Text Anonymizer</h1>
        <TextInput text={text} onChange={handleTextChange} />
        <br />
        <AnonymizeButton onClick={handleAnonymize} disabled={!text}>
          Anonymize with GPT
        </AnonymizeButton>
        <AnonymizeButton onClick={handleAnonymizeSpacy} disabled={!text}>
          Anonymize with spaCy
        </AnonymizeButton>
        {error && <ErrorMessage message={error} />}
        {anonymizedText && (
          <>
            <h2>Anonymized Text (GPT):</h2>
            <p>{anonymizedText}</p>
          </>
        )}
        {anonymizedTextSpacy && (
          <>
            <h2>Anonymized Text (spaCy):</h2>
            <p>{anonymizedTextSpacy}</p>
          </>
        )}
      </header>
    </div>
  );
};

export default App;
