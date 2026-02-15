import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { CirclesPresenter } from "./presenter";

type CircleResponse = components["schemas"]["dto.CircleResponse"];

export async function CirclesContainer() {
  const session = await auth();

  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  let circles: CircleResponse[] = [];
  let count = 0;
  let maxCircles = 3;
  let fetchError: string | null = null;

  const { data, error } = await client.GET("/api/circles", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
  });

  if (error) {
    fetchError = error.error;
  } else {
    circles = data.circles;
    count = data.count;
    maxCircles = data.max_circles;
  }

  return (
    <CirclesPresenter
      circles={circles}
      count={count}
      maxCircles={maxCircles}
      initialError={fetchError}
    />
  );
}
