
import { useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { Routes, Route } from 'react-router-dom';
import Header from './Header';
import Footer from './Footer';
import About from './About';
import Leaderboard from './Leaderboard';
import Arena from './Arena';

function App() {
  useEffect(() => {
    let userId = sessionStorage.getItem('userId');
    if (!userId) {
      userId = uuidv4();
      sessionStorage.setItem('userId', userId);
    }
  }, []);

  return (
    <div className="App">
      <Header />
      <main>
        <Routes>
          <Route path="/" element={<About />} />
          <Route path="/about" element={<About />} />
          <Route path="/arena" element={<Arena />} />
          <Route path="/leaderboard" element={<Leaderboard />} />
        </Routes>
      </main>
      <Footer />
    </div>
  );
}

export default App;
