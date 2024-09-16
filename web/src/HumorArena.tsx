import React, { useEffect, useState } from 'react';
import { ArenaApi, Configuration, V1GetChoicesResponse, V1Winner } from './apiClient';

const apiBasePath = process.env.REACT_APP_API_BASE_URL || '';

const config = new Configuration({ basePath: apiBasePath });
const api = new ArenaApi(config);

type Choice = {
  id: string;
  theme: string;
  leftJoke: string;
  rightJoke: string;
};

const JokeComparison: React.FC = () => {
  const [choice, setChoice] = useState<Choice | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchChoices = async () => {
    setLoading(true);
    setError(null);

    try {
      const response: V1GetChoicesResponse = await api.arenaGetChoices();
      setChoice({
        id: response.id!,
        theme: response.theme!,
        leftJoke: response.leftJoke!,
        rightJoke: response.rightJoke!,
      });
    } catch (err) {
      setError('Failed to fetch jokes.');
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
    // eslint-disable-next-line react-hooks/exhaustive-deps
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
    <div>
      <h2>Theme: {choice.theme}</h2>
      <div style={{ display: 'flex', justifyContent: 'space-around' }}>
        <div style={{ width: '40%' }}>
          <p>{choice.leftJoke}</p>
          <button onClick={() => handleVote(V1Winner.Left)}>Vote Left</button>
        </div>
        <div style={{ width: '40%' }}>
          <p>{choice.rightJoke}</p>
          <button onClick={() => handleVote(V1Winner.Right)}>Vote Right</button>
        </div>
      </div>
      <div style={{ marginTop: '20px' }}>
        <button onClick={() => handleVote(V1Winner.Both)}>Both are great!</button>
        <button onClick={() => handleVote(V1Winner.None)}>Neither</button>
      </div>
    </div>
  );
};

export default JokeComparison;
