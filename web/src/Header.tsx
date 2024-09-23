import React from 'react';
import './Header.css'; // We'll create this CSS file for styling

const Header: React.FC = () => {
  return (
    <header className="header">
      <h1>Humor Arena</h1>
      <nav>
        <ul className="menu">
          <li><a href="/">Joke Comparison</a></li>
          {/* <li><a href="/leaderboard">Leaderboard</a></li> */}
        </ul>
      </nav>
    </header>
  );
};

export default Header;