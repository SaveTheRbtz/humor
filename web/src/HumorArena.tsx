import React, { useEffect, useState, MouseEvent, MouseEventHandler } from 'react';
import { ArenaApi, Configuration, V1GetChoicesResponse, V1Winner } from './apiClient';
import { getErrorMessage } from './errorUtils';
import './HumorArena.css';

const apiBasePath = process.env.REACT_APP_API_BASE_URL || '';

const config = new Configuration({ basePath: apiBasePath });
const api = new ArenaApi(config);

type Choice = {
  id: string;
  theme: string;
  leftJoke: string;
  rightJoke: string;
  known: V1Winner;
};

type JokeCardProps = {
  jokeText: string;
  onVote: MouseEventHandler<HTMLDivElement>;
};
const JokeCard: React.FC<JokeCardProps> = ({ jokeText, onVote,  }) => {
  const [isKnown, setIsKnown] = useState<boolean>(false);

  const handleMarkAsKnown = (event: MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    const newIsKnown = !isKnown;
    setIsKnown(newIsKnown);
  };

  return (
    <div
      className={`joke-card ${isKnown ? 'dimmed' : ''}`}
      onClick={onVote}
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

const JokeComparison: React.FC = () => {
  const [choice, setChoice] = useState<Choice | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchChoices = async () => {
    setLoading(true);
    setError(null);

    try {
      const response: V1GetChoicesResponse = await api.arenaGetChoices(
        {sessionId: sessionStorage.getItem('userId') || ''},
      );
      setChoice({
        id: response.id!,
        theme: response.theme!,
        leftJoke: response.leftJoke!,
        rightJoke: response.rightJoke!,
        known: V1Winner.None,
      });
    } catch (err: any) {
      const errorMessage = await getErrorMessage(err);
      setError(`Failed to fetch jokes: ${errorMessage}`);
    } finally {
      setLoading(false);
    }
  };

  const handleVote = async (winner: V1Winner) => {
    if (!choice) return;

    try {
      await api.arenaRateChoices({
        id: choice.id,
        body: {
          winner: winner,
        },
      });
      // Fetch new jokes after voting
      fetchChoices();
    } catch (err) {
      setError('Failed to submit your choice.');
    }
  };

  useEffect(() => {
    fetchChoices();
  }, []);

  if (loading) {
    return <div>Loading jokes...</div>;
  }

  if (error) {
    return (
      <div>
        <p>{error}</p>
        <button onClick={fetchChoices}>Retry</button>
      </div>
    );
  }

  if (!choice) {
    return <div>No jokes available.</div>;
  }

  return (
    <div className="joke-comparison">
      <h2>{choice.theme}</h2>
      <div className="jokes-container">
        <JokeCard
          jokeText={choice.leftJoke}
          onVote={() => handleVote(V1Winner.Left)}
        />
        <JokeCard
          jokeText={choice.rightJoke}
          onVote={() => handleVote(V1Winner.Right)}
        />
      </div>
      <div className="additional-options">
        <button className="both-button" onClick={() => handleVote(V1Winner.Both)}>Both are great</button>
        <button className="neither-button" onClick={() => handleVote(V1Winner.None)}>Neither are good</button>
      </div>
    </div>
  );
};

export default JokeComparison;
