import React, { useState } from 'react';
import { useQuery } from "convex/react";
import { api } from "../convex/_generated/api";
import './index.css';

const Leaderboard: React.FC = () => {
  const topJokes = useQuery(api.jokesQueries.getTopJokes);
  const [popupText, setPopupText] = useState<string | null>(null);

  const handlePopupClick = (text: string | undefined) => {
    if (!text) return;
    setPopupText(text);
  };

  const handleClosePopup = () => {
    setPopupText(null);
  };

  return (
    <div className="leaderboard-container">
      <h1>Leaderboard</h1>
      {topJokes ? (
        <ul className="leaderboard-list">
          {topJokes.map((joke, _) => (
            <li key={joke._id} className="leaderboard-item">
              <p>{joke.text}</p>
              <p>
                Score: {joke.score.toString()} Views: {joke.views.toString()} Score/View: {(Number(joke.score) / Number(joke.views)).toFixed(2)}
                <span>, </span>
                Topic: <b>{joke.topicName}</b>
                <span>, </span>
                Policy: <span className="policy-name" onClick={() => handlePopupClick(joke.policyText)}>{joke.policyName}</span>
                <span>, </span>
                <span className="policy-name" onClick={() => handlePopupClick(joke.assocv1Text)}>Assocv1</span>
                <span>, </span>
                <span className="policy-name" onClick={() => handlePopupClick(joke.assocv2Text)}>Assocv2</span>
                <span>, </span>
                <span className="policy-name" onClick={() => handlePopupClick(joke.assocv3Text)}>Assocv3</span>
              </p>
            </li>
          ))}
        </ul>
      ) : (
        <p>Loading...</p>
      )}
      {popupText && (
        <div className="popup">
          <div className="popup-inner">
            <p>{popupText}</p>
            <button onClick={handleClosePopup}>Close</button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Leaderboard;