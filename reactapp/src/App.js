import logo from './logo.svg';
import './App.css';
import { useEffect } from 'react';

function App() {

  useEffect(() => {
    // Run once on page load
    const url = `${process.env.REACT_APP_API_URL}getPosts`; // get API host from .env file
    fetch(url)
      .then(res => res.json())
      .then(res => {
        console.log({
          res
        });
        // Printing out just for response visualisation
      })
      .catch(reason => {
        console.error({
          reason
        })
        // print reason if request failed
      })
  }, [])

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
