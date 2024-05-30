import React, { useState, ChangeEvent, FC, useCallback } from 'react';
import './App.css';
import { fetchAnonymizedText } from './api';
import TextInput from './components/TextInput';
import AnonymizeButton from './components/AnonymizeButton';
import ErrorMessage from './components/ErrorMessage';

const App: FC = () => {
  const [text, setText] = useState('');
  const [anonymizedText, setAnonymizedText] = useState('');
  const [error, setError] = useState<string | null>(null);

  const handleTextChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    setText(e.target.value);
  };

  const handleAnonymize = useCallback(async () => {
    setError(null);
    try {
      const data = await fetchAnonymizedText(text);
      setAnonymizedText(data.anonymizedText);
    } catch (error) {
      console.error('Error:', error);
      setError('Failed to anonymize the text. Please try again.');
    }
  }, [text]);

  return (
    <div className="App">
      <header className="App-header">
        <h1>Text Anonymizer</h1>
        <TextInput text={text} onChange={handleTextChange} />
        <br />
        <AnonymizeButton onClick={handleAnonymize} disabled={!text} />
        {error && <ErrorMessage message={error} />}
        {anonymizedText && (
          <>
            <h2>Anonymized Text:</h2>
            <p>{anonymizedText}</p>
          </>
        )}
      </header>
    </div>
  );
};

export default App;
