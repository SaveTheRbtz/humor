import Header from './Header';
import { useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import JokeComparison from './HumorArena';

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
        <JokeComparison />
      </main>
    </div>
  );
}

export default App;