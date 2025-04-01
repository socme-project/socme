import { writable } from "svelte/store";

export const user = writable<{
  ID: string;
  Username: string;
  GitHubID: string;
  GoogleID: string;
  Role: string;
  Avatar: string;
} | null>(null);
