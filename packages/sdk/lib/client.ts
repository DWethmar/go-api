import { Entry } from "@go-api/models";

export interface ClientConfig {
  host: string;
}

export class Client {
  private host: string;

  constructor(config: ClientConfig) {
    console.log("init client.");

    // Todo validate config.
    this.host = config.host;
  }

  getAll(): Promise<Entry[]> {
    return new Promise<Entry[]>((resolve, reject) => {
      fetch(this.host).then(function (response) {
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
}
