import React, { useState, useEffect } from "react";
import { Hello, Bye } from "@go-api/elements";
import { createClient, Entry } from "@go-api/sdk";
import logo from "./krat_logo.svg";

import "@go-api/elements/dist/index.css";

import "./App.css";

const client = createClient();

function App() {
  const [entries, setEntries] = useState<Entry[]>([]);

  useEffect(() => {
    client.getAll().then((entries) => setEntries(entries));
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
      </header>
      <Hello></Hello>
      <Bye></Bye>

      <ul>
        {entries.map((entry, index) => (
          <li key={index}>
            {entry.name} {entry.createdOn} {entry.fields.nl.a}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
