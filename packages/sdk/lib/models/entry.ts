import { Locale } from "./locale";

export type EntryFieldKey = string;

export type EntryFieldValue = string | number;

export interface Entry {
  name: string;
  createdOn: string;
  updatedOn: string;
  fields: Record<Locale, Record<EntryFieldKey, EntryFieldValue>>;
}
