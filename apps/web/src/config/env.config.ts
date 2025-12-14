export const envConfig = {
  auth: {
    github: {
      clientId: process.env.AUTH_GITHUB_ID || "",
      clientSecret: process.env.AUTH_GITHUB_SECRET || "",
    },
    jwt: {
      secret: process.env.AUTH_SECRET || "",
    },
  },
  nextjs: {
    runtime: process.env.NEXT_RUNTIME || "edge",
    apiBaseUrl:
      process.env.API_BASE_URL ||
      (process.env.NODE_ENV === "production" ? "" : "http://localhost:3120"),
    appEnv: process.env.NEXT_PUBLIC_APP_ENV || "development",
  },
} as const;
