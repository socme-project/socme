import { writable } from "svelte/store";

export type User = {
  ID: string;
  GithubID: string;
  Name: string;
  Role: string;
};

export const user = writable<User | null>(null);

if (typeof window !== "undefined") {
  const storedUser = localStorage.getItem("currentUser");
  if (storedUser) {
    user.set(JSON.parse(storedUser));
  }
  user.subscribe((value) => {
    if (value) {
      localStorage.setItem("currentUser", JSON.stringify(value));
    } else {
      localStorage.removeItem("currentUser");
    }
  });
}
