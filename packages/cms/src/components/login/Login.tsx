import React, { useState, FormEvent } from "react";
import "./Login.css";

function login(
  loginUrl: string,
  username: string,
  password: string
): Promise<boolean> {
  return fetch(loginUrl, {
    method: "POST", // or 'PUT'
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ username, password }),
  }).then((response) => response.status === 200);
}

export interface LoginProps {
  loginUrl: string;
}

export const Login: React.FC<LoginProps> = (props) => {
  const { loginUrl } = props;

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  function validateForm() {
    return email.length > 0 && password.length > 0;
  }

  function handleSubmit(event: FormEvent) {
    event.preventDefault();

    if (validateForm()) {
      login(loginUrl, email, password);
    }
    console.log(event);
  }

  return (
    <div className="Login">
      <form onSubmit={handleSubmit}>
        <fieldset>
          <label>Email</label>
          <input
            type="text"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </fieldset>
        <fieldset>
          <label>password</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </fieldset>
        <input type="submit" />
      </form>
    </div>
  );
};
