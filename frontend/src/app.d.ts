// See https://svelte.dev/docs/kit/types#app.d.ts

import type PocketBase, { AuthModel } from "pocketbase";

// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
    interface Locals {
      pb: PocketBase;
      user: AuthModel | Null;
    }
    // interface PageData {}
    // interface PageState {}
    // interface Platform {}
  }
}

export { };
