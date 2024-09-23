import React, { MouseEvent } from 'react';
import { V1Winner } from './apiClient';
type JokeCardProps = {
    jokeText: string;
    onVote: (winner: V1Winner) => void;
    selected: V1Winner;
    isKnown: boolean;
    onMarkAsKnown: (selected: V1Winner) => void;
};

const JokeCard: React.FC<JokeCardProps> = ({ jokeText, onVote, selected, isKnown, onMarkAsKnown }) => {
    const handleMarkAsKnown = (event: MouseEvent<HTMLButtonElement>) => {
        event.stopPropagation();
        onMarkAsKnown(selected);
    };

    return (
        <div
            className={`joke-card ${isKnown ? 'dimmed' : ''}`}
            onClick={() => onVote(selected)}
        >
            <button
                className="close-button"
                onClick={handleMarkAsKnown}
                title={isKnown ? 'Mark joke as unknown' : 'Mark joke as known'}
            >
                ×
            </button>
            <p>{jokeText}</p>
        </div>
    );
};

export { JokeCard };