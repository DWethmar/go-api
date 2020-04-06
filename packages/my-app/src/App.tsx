import React from 'react';
import { Hello } from "@go-api/my-monorepo-ui-lib"
import logo from './krat_logo.svg';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
      </header>
      <Hello></Hello>
    </div>
  );
}

export default App;
