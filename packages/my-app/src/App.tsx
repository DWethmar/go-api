import React from 'react';
import { Hello, Bye } from "@go-api/elements"
import logo from './krat_logo.svg';

import './App.css';

import "@go-api/my-monorepo-ui-lib/lib/index.css"

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
      </header>
      <Hello></Hello>
      <Bye></Bye>
    </div>
  );
}

export default App;
