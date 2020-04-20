import React from "react";
import "./Test.css";

function test(url: string): Promise<string> {
  return fetch(url, {
    method: "GET",
  }).then((response) => response.text());
}

export interface TestProps {
  url: string;
}

export const Test: React.FC<TestProps> = (props) => {
  const { url } = props;

  function handleClick() {
    console.log(test(url));
  }

  return (
    <div className="test">
      <button onClick={handleClick}>Test</button>
    </div>
  );
};
