import React from 'react';
import './index.css';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import Compare from './Compare';
import Leaderboard from './Leaderboard';

const App: React.FC = () => {
  return (
    <Router>
      <div className="app-container">
        <header className="header">
          <nav>
            <ul>
              <li>
                <Link to="/">Compare</Link>
              </li>
              <li>
                <Link to="/leaderboard">Leaderboard</Link>
              </li>
            </ul>
          </nav>
        </header>
        <main>
          <Routes>
            <Route path="/" element={<Compare />} />
            <Route path="/leaderboard" element={<Leaderboard />} />
          </Routes>
        </main>
        <footer className="footer">
          <p>Paper: <a href="https://arxiv.org/abs/2405.07280">Humor Mechanics: Advancing Humor Generation with Multistep Reasoning</a></p>
        </footer>
      </div>
    </Router>
  );
};

export default App;