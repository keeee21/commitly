import NextAuth from "next-auth";
import GitHub from "next-auth/providers/github";
import { envConfig } from "./env.config";

// GitHub Profile型定義
interface GitHubProfile {
  id: number;
  login: string;
  email: string | null;
  name: string | null;
  avatar_url: string;
}

export const { handlers, auth, signIn, signOut } = NextAuth({
  trustHost: true,
  providers: [
    GitHub({
      clientId: envConfig.auth.github.clientId,
      clientSecret: envConfig.auth.github.clientSecret,
    }),
  ],
  callbacks: {
    async jwt({ token, account, profile }) {
      // 初回ログイン時にGitHub情報をtokenに保存し、DBにも保存
      if (account && profile) {
        const githubProfile = profile as unknown as GitHubProfile;
        const githubUserId = githubProfile.id;
        const githubUsername = githubProfile.login;
        const email = githubProfile.email;

        token.githubUserId = githubUserId;
        token.githubUsername = githubUsername;

        // APIにユーザー情報を保存
        try {
          const apiUrl =
            process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
          const response = await fetch(`${apiUrl}/api/auth/callback`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              github_user_id: githubUserId,
              github_username: githubUsername,
              email: email || "",
              avatar_url: githubProfile.avatar_url,
            }),
          });

          if (!response.ok) {
            console.error("Failed to save user to API:", await response.text());
          } else {
            console.log("User saved to API successfully");
          }
        } catch (error) {
          console.error("Error saving user to API:", error);
        }
      }

      return token;
    },
    async session({ session, token }) {
      if (token?.sub) {
        session.user.id = token.sub;
      }

      // GitHub情報をセッションに追加
      if (
        token.githubUserId !== undefined &&
        token.githubUsername !== undefined
      ) {
        session.user.githubUserId = token.githubUserId as number;
        session.user.githubUsername = token.githubUsername as string;
      }

      return session;
    },
  },
});
