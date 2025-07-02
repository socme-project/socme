import { user } from "$src/lib/stores/user";
import { SESSION_COOKIE_NAME } from "$src/lib/variables";
import axios from "axios";
import type { Hooks } from "sv-router";
import { navigate } from "sv-router/generated";
import { toast } from "svelte-sonner";
import { get } from "svelte/store";

const PUBLIC_ROUTES: string[] = ["", "/", "/auth/github/callback"];

function deleteCookie(cookieName: string) {
  document.cookie = `${cookieName}=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
}

function checkUserSession(): [boolean, string] {
  const sessionCookie = document.cookie
    .split("; ")
    .find((cookie) => cookie.startsWith(`${SESSION_COOKIE_NAME}=`));

  if (!sessionCookie) {
    deleteCookie(SESSION_COOKIE_NAME);
    deleteCookie("exp");
    return [false, ""];
  }

  const expCookie = document.cookie
    .split("; ")
    .find((cookie) => cookie.startsWith(`exp=`));

  if (!expCookie) {
    deleteCookie(SESSION_COOKIE_NAME);
    deleteCookie("exp");
    return [false, ""];
  }

  const now = Math.floor(Date.now() / 1000);
  if (expCookie && parseInt(expCookie.split("=")[1]) < now) {
    deleteCookie(SESSION_COOKIE_NAME);
    deleteCookie("exp");
    return [false, ""];
  }

  const u = get(user);
  if (user === null || u === undefined || !u?.Role) {
    deleteCookie(SESSION_COOKIE_NAME);
    deleteCookie("exp");
    return [false, ""];
  }

  const role = u.Role;

  return [true, role];
}

export default {
  beforeLoad({ pathname }) {
    const [isSessionValid, userRole] = checkUserSession();
    if (!isSessionValid || userRole === "guest") {
      if (PUBLIC_ROUTES.includes(pathname)) {
        // User not connected, public route
        return;
      } else {
        // User not connected, private route
        toast.error("You must be logged in to access this page.");
        throw navigate("/");
      }
    }

    // User connected, public route
    if (PUBLIC_ROUTES.includes(pathname)) {
      toast.error("You are already logged in.");
      throw navigate("/dashboard");
    }

    // User connected, private route
    return;
  },
} satisfies Hooks;
