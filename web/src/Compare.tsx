import React, { useState, useEffect } from 'react';
import './index.css';
import { useQuery, useMutation } from "convex/react";
import { api } from "../convex/_generated/api";

const Compare: React.FC = () => {
  const getRandomJokesQuery = useQuery(api.jokesQueries.getRandomJokes, { count: 2 });
  const updateJokes = useMutation(api.jokesMutations.updateJokes);
  const [selectedJoke, setSelectedJoke] = useState<number | 'neither' | null>(null);

  const handleSelect = (selected: number | 'neither') => {
    setSelectedJoke(selected);
    if (getRandomJokesQuery) {
      const [joke1, joke2] = getRandomJokesQuery;
      if (selected === 'neither') {
        updateJokes({ jokeId1: joke1._id, jokeId2: joke2._id, winningJokeId: null });
      } else {
        const winningJokeId = selected === 0 ? joke1._id : joke2._id;
        updateJokes({ jokeId1: joke1._id, jokeId2: joke2._id, winningJokeId });
      }
    }
  };

  useEffect(() => {
    if (getRandomJokesQuery) {
      setSelectedJoke(null);
    }
  }, [getRandomJokesQuery]);

  return (
    <div className="compare-container">
      <h1>Which Joke is Funnier?</h1>
      {getRandomJokesQuery ? (
        <div className="jokes-container">
          <div className="joke-row">
            <div
              className={`joke-card ${selectedJoke === 0 ? 'selected-blue' : ''}`}
              onClick={() => handleSelect(0)}
            >
              <p>{getRandomJokesQuery[0].text}</p>
            </div>
            <div
              className={`joke-card ${selectedJoke === 1 ? 'selected-blue' : ''}`}
              onClick={() => handleSelect(1)}
            >
              <p>{getRandomJokesQuery[1].text}</p>
            </div>
          </div>
          <div className="joke-row">
            <div
              className={`joke-card neither ${selectedJoke === 'neither' ? 'selected-red' : ''}`}
              onClick={() => handleSelect('neither')}
            >
              <p>Both suck</p>
            </div>
          </div>
        </div>
      ) : (
        <p>Loading...</p>
      )}
    </div>
  );
};

export default Compare;