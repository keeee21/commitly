"use server";

import { signIn, signOut } from "@/config/auth.config";

export async function signInWithGitHub() {
  await signIn("github", { redirectTo: "/" });
}

export async function signOutAction() {
  await signOut({ redirectTo: "/login" });
}
