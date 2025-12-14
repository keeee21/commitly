import type { DefaultSession } from "next-auth";

declare module "next-auth" {
  interface Session {
    user: {
      id: string;
      githubUserId?: number; // GitHub User ID (一意、変更不可、DBに保存)
      githubUsername?: string; // GitHub Username (変更可能、DBに保存)
    } & DefaultSession["user"];
  }

  interface User {
    id: string;
    name?: string | null;
    email?: string | null;
    image?: string | null;
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    sub: string;
    githubUserId?: number;
    githubUsername?: string;
  }
}
