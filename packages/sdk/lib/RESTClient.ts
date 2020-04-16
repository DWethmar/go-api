import { Entry } from "./models/entry";
import { Client } from "./client";

export interface RESTClientConfig {
  host: string;
}

function jsonRequest(request: Request) {
  return new Promise<Entry[]>((resolve, reject) => {
    fetch(request).then(function (response) {
      if (response.status !== 200) {
        reject(
          new Error(
            "Looks like there was a problem. Status Code: " + response.status
          )
        );
        return;
      }

      // Examine the text in the response
      response.json().then(function (data) {
        resolve(data);
      });
    });
  });
}

export class RESTClient implements Client {
  private host: string;

  constructor(config: RESTClientConfig) {
    console.log("init client.");

    // Todo validate config.
    this.host = config.host;
  }

  getAll(): Promise<Entry[]> {
    const request = new Request(this.host, {
      method: "get",
    });
    // request.headers.append("ananas", "kaas");
    return jsonRequest(request);
  }
}
