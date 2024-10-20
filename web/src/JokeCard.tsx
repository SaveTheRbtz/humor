import React from 'react';
import { FaTwitter } from 'react-icons/fa';
import './JokeCard.css';

type JokeCardProps = {
  jokeText: string;
};

const JokeCard: React.FC<JokeCardProps> = ({ jokeText }) => {
  const handleShare = (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent any parent handlers from being notified of the event

    const tweetText = `${jokeText}\n\nhttps://humor.ph34r.me/arena #humorarena`;
    const twitterUrl = `https://twitter.com/intent/tweet?text=${encodeURIComponent(tweetText)}`;
    window.open(twitterUrl, '_blank');
  };

  return (
    <div className="joke-card">
      <p>{jokeText}</p>
      <button className="share-icon" onClick={handleShare} aria-label="Share on Twitter">
        <FaTwitter size={24} color="#1DA1F2" />
      </button>
    </div>
  );
};

export { JokeCard };