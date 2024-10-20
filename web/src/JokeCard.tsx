import React, { MouseEvent } from 'react';
import { V1Winner } from './apiClient';
type JokeCardProps = {
    jokeText: string;
    onVote: (winner: V1Winner) => void;
    selected: V1Winner;
    isKnown: boolean;
};

const JokeCard: React.FC<JokeCardProps> = ({ jokeText, onVote, selected, isKnown}) => {
    return (
        <div
            className={`joke-card ${isKnown ? 'dimmed' : ''}`}
        >
            <p>{jokeText}</p>
        </div>
    );
};

export { JokeCard };