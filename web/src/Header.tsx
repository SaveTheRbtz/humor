import React, { useState } from 'react';
import './Header.css';

const Header: React.FC = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  return (
    <header className="header">
      <h1>Humor Arena</h1>
      <nav>
        {/* Burger Menu Icon */}
        <div className="burger" onClick={toggleMenu}>
          <div className={`line ${isMenuOpen ? 'rotate1' : ''}`}></div>
          <div className={`line ${isMenuOpen ? 'fade' : ''}`}></div>
          <div className={`line ${isMenuOpen ? 'rotate2' : ''}`}></div>
        </div>
        {/* Navigation Menu */}
        <ul className={`menu ${isMenuOpen ? 'menu-open' : ''}`}>
          <li><a href="/arena">Arena (side-by-side)</a></li>
          <li><a href="/leaderboard">Leaderboard</a></li>
          <li><a href="/">About</a></li>
        </ul>
      </nav>
    </header>
  );
};

export default Header;