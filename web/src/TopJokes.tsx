import React, { useEffect, useState } from 'react';
import { FaTwitter } from 'react-icons/fa';
import './TopJokes.css';
import { ArenaApi, Configuration, V1GetTopJokesResponse, V1TopJokesEntry } from './apiClient';

const apiBasePath = process.env.REACT_APP_API_BASE_URL || '';
const config = new Configuration({ basePath: apiBasePath });
const api = new ArenaApi(config);

const TopJokes: React.FC = () => {
  const [topJokesEntries, setTopJokesEntries] = useState<V1TopJokesEntry[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const handleShare = (jokeText: string) => {
    const tweetText = `${jokeText}\n\nhttps://humor.ph34r.me/arena #humorarena`;
    const twitterUrl = `https://twitter.com/intent/tweet?text=${encodeURIComponent(tweetText)}`;
    window.open(twitterUrl, '_blank');
  };

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
      <div className="about-container">
        <p>
          Since our goal is to automate machine understanding of humor, here we demonstrate fully automatic sorting of 2000+ automatically generated jokes, publishing the top 50 according to machine scores.
          This is an early beta version and we are working on improving it, with code and algorithms to be published soon. Enjoy!
        </p>
      </div>
      <table className="top-jokes-table">
        <thead>
          <tr>
            <th>Rank</th>
            <th>Joke</th>
            <th>Share</th>
          </tr>
        </thead>
        <tbody>
          {topJokesEntries.map((entry) => (
            <tr key={entry.rank}>
              <td>{entry.rank}</td>
              <td>{entry.text}</td>
              <td>
                <button
                  className="share-icon-top-jokes"
                  onClick={() => handleShare(entry.text!)}
                  aria-label="Share on Twitter"
                >
                  <FaTwitter size={24} color="#1DA1F2" />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default TopJokes;
