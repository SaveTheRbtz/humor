import React, { useEffect, useState } from 'react';
import './Leaderboard.css';
import { ArenaApi, Configuration, V1GetLeaderboardResponse, V1LeaderboardEntry } from './apiClient';

const apiBasePath = process.env.REACT_APP_API_BASE_URL || '';
const config = new Configuration({ basePath: apiBasePath });
const api = new ArenaApi(config);

const Leaderboard: React.FC = () => {
  const [leaderboardEntries, setLeaderboardEntries] = useState<V1LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchLeaderboard = async () => {
    setLoading(true);
    setError(null);

    try {
      const response: V1GetLeaderboardResponse = await api.arenaGetLeaderboard({});
      setLeaderboardEntries(response.entries || []);
    } catch (err: any) {
      setError('Failed to fetch leaderboard data.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLeaderboard();
  }, []);

  if (loading) {
    return <div>Loading leaderboard...</div>;
  }

  if (error) {
    return (
      <div>
        <p>{error}</p>
        <button onClick={fetchLeaderboard}>Retry</button>
      </div>
    );
  }

  return (
    <div className="leaderboard-container">
      <h1>Leaderboard</h1>
      <table className="leaderboard-table">
        <thead>
          <tr>
            <th>Rank</th>
            <th>Model</th>
            <th>Score</th>
          </tr>
        </thead>
        <tbody>
          {leaderboardEntries
            .sort((a, b) => b.bradleyterrScore! - a.bradleyterrScore!)
            .map((entry, index) => (
              <tr key={entry.model}>
                <td>{index + 1}</td>
                <td>{entry.model}</td>
                <td>{entry.bradleyterrScore!.toFixed(4)}</td>
              </tr>
            ))}
        </tbody>
      </table>
    </div>
  );
};

export default Leaderboard;