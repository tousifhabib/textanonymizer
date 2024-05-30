import React, { FC } from 'react';

interface AnonymizeButtonProps {
  onClick: () => void;
  disabled: boolean;
}

const AnonymizeButton: FC<AnonymizeButtonProps> = ({ onClick, disabled }) => (
  <button onClick={onClick} disabled={disabled}>
    Anonymize
  </button>
);

export default AnonymizeButton;
