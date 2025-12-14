import createClient from "openapi-fetch";
import type { paths } from "./schema";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const client = createClient<paths>({ baseUrl: apiUrl });

// 認証付きAPIクライアントを作成するヘルパー
export function createAuthClient(githubUserId: number) {
  return createClient<paths>({
    baseUrl: apiUrl,
    headers: {
      "X-GitHub-User-ID": String(githubUserId),
    },
  });
}
