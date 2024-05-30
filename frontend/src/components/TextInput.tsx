import React, { ChangeEvent, FC } from 'react';

interface TextInputProps {
  text: string;
  onChange: (e: ChangeEvent<HTMLTextAreaElement>) => void;
}

const TextInput: FC<TextInputProps> = ({ text, onChange }) => (
  <textarea
    value={text}
    onChange={onChange}
    rows={10}
    cols={50}
    placeholder="Enter text here..."
  />
);

export default TextInput;
