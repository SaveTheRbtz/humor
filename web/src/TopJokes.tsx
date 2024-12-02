import React, { useEffect, useState } from 'react';
import './TopJokes.css';
import { ArenaApi, Configuration, V1GetTopJokesResponse, V1TopJokesEntry } from './apiClient';

const apiBasePath = process.env.REACT_APP_API_BASE_URL || '';
const config = new Configuration({ basePath: apiBasePath });
const api = new ArenaApi(config);

const TopJokes: React.FC = () => {
  const [topJokesEntries, setTopJokesEntries] = useState<V1TopJokesEntry[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchTopJokes = async () => {
    setLoading(true);
    setError(null);

    try {
      const response: V1GetTopJokesResponse = await api.arenaGetTopJokes({});
      setTopJokesEntries(response.entries || []);
    } catch (err: any) {
      setError('Failed to fetch top jokes data.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTopJokes();
  }, []);

  if (loading) {
    return <div>Loading top jokes...</div>;
  }

  if (error) {
    return (
      <div>
        <p>{error}</p>
        <button onClick={fetchTopJokes}>Retry</button>
      </div>
    );
  }

  return (
    <div className="top-jokes-container">
      <h1>Top Jokes</h1>
      <table className="top-jokes-table">
        <thead>
          <tr>
            <th>Rank</th>
            <th>Joke</th>
          </tr>
        </thead>
        <tbody>
          {topJokesEntries
            .map((entry, _) => (
              <tr key={entry.rank}>
                <td>{entry.rank}</td>
                <td>{entry.text}</td>
              </tr>
            ))}
        </tbody>
      </table>
    </div>
  );
};

export default TopJokes;