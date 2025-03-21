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
  const [tooltipVisible, setTooltipVisible] = useState<boolean>(false);

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

  const toggleTooltip = () => {
    setTooltipVisible(!tooltipVisible);
  };

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
      <h1>Model Leaderboard</h1>
      <table className="leaderboard-table">
        <thead>
          <tr>
            <th>#</th>
            <th>Model</th>
            <th>Votes</th>
            {/* <th>
              Newman
              <span className="tooltip" onClick={toggleTooltip}>
                 &#9432;
                {tooltipVisible && (
                  <span className="tooltiptext">
                    <div>Newman Score:</div>
                    <a href="https://arxiv.org/abs/2207.00076" target="_blank" rel="noopener noreferrer">
                      Efficient computation of rankings from pairwise comparisons
                    </a>
                  </span>
                )}
              </span>
            </th> */}
            <th>
              Elo
            </th>
          </tr>
        </thead>
        <tbody>
          {leaderboardEntries
            .sort((a, b) => b.eloScore! - a.eloScore!)
            .map((entry, index) => (
              <tr key={entry.model}>
                <td>{index + 1}</td>
                <td>{entry.model}</td>
                <td>
                  {entry.votes}
                  <span className="ci">
                    &nbsp;↑{entry.votesGood!}/↓{entry.votesBad!}
                  </span>
                </td>
                {/* <td>
                  {entry.newmanScore!.toFixed(2)}
                  <span className="ci">
                    &nbsp;+{entry.newmanCIUpper!.toFixed(2)}/-{entry.newmanCILower!.toFixed(2)}
                  </span>
                </td> */}
                <td>
                  {entry.eloScore!.toFixed(0)}
                  <span className="ci">
                    &nbsp;+{entry.eloCIUpper!.toFixed(0)}/-{entry.eloCILower!.toFixed(0)}
                  </span>
                </td>
              </tr>
            ))}
        </tbody>
      </table>
    </div>
  );
};

export default Leaderboard;