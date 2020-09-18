import { Entry } from "./index";
import { RESTClient } from "./RESTClient";

export interface Client {
  getAll(): Promise<Entry[]>;
}

export function createClient(): Client {
  return new RESTClient({
    host: "http://localhost:8080",
  });
}
