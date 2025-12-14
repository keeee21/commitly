"use server";

import { signIn } from "@/config/auth.config";

export async function signInWithGitHub() {
  await signIn("github", { redirectTo: "/" });
}
