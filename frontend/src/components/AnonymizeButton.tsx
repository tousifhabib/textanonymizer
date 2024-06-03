import React, { FC } from 'react';

interface AnonymizeButtonProps {
  onClick: () => void;
  disabled: boolean;
  children: React.ReactNode;
}

const AnonymizeButton: FC<AnonymizeButtonProps> = ({ onClick, disabled, children }) => (
  <button onClick={onClick} disabled={disabled}>
    {children}
  </button>
);

export default AnonymizeButton;