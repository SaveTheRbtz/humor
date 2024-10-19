import React, { useEffect, useState, MouseEvent, MouseEventHandler } from 'react';
import { ArenaApi, Configuration, V1GetChoicesResponse, V1Winner } from './apiClient';
import {JokeCard} from './JokeCard';
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

const JokeComparison: React.FC = () => {
  const [choice, setChoice] = useState<Choice | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);


  const handleMarkAsKnown = (selected: V1Winner) => {
    if (!choice) return;

    let newKnown = choice.known;

    switch (selected) {
      case V1Winner.Left:
        if (choice.known === V1Winner.Left || choice.known === V1Winner.Both) {
          newKnown = choice.known === V1Winner.Both ? V1Winner.Right : V1Winner.None;
        } else {
          newKnown = choice.known === V1Winner.Right ? V1Winner.Both : V1Winner.Left;
        }
        break;
      case V1Winner.Right:
        if (choice.known === V1Winner.Right || choice.known === V1Winner.Both) {
          newKnown = choice.known === V1Winner.Both ? V1Winner.Left : V1Winner.None;
        } else {
          newKnown = choice.known === V1Winner.Left ? V1Winner.Both : V1Winner.Right;
        }
        break;
      default:
        break;
    }
    
    setChoice({ ...choice, known: newKnown });
  };

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
          known: choice.known,
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
          onVote={handleVote}
          selected={V1Winner.Left}
          isKnown={choice.known === V1Winner.Left || choice.known === V1Winner.Both}
          onMarkAsKnown={handleMarkAsKnown}
        />
        <JokeCard
          jokeText={choice.rightJoke}
          onVote={handleVote}
          selected={V1Winner.Right}
          isKnown={choice.known === V1Winner.Right || choice.known === V1Winner.Both}
          onMarkAsKnown={handleMarkAsKnown}
        />
      </div>
      <div className="additional-options">
      <button className="button" onClick={() => handleVote(V1Winner.Left)}>Left is better</button>
        <button className="both-button" onClick={() => handleVote(V1Winner.Both)}>Both are great</button>
        <button className="neither-button" onClick={() => handleVote(V1Winner.None)}>Neither are good</button>
        <button className="button" onClick={() => handleVote(V1Winner.Right)}>Right is better</button>
      </div>
    </div>
  );
};

export default JokeComparison;
