import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { ActivityPresenter } from "./presenter";

type ActivityStreamResponse =
  components["schemas"]["dto.ActivityStreamResponse"];
type RhythmResponse = components["schemas"]["dto.RhythmResponse"];

export async function ActivityContainer() {
  const session = await auth();

  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  let streamData: ActivityStreamResponse | null = null;
  let rhythmData: RhythmResponse | null = null;
  let fetchError: string | null = null;

  const headers = {
    "X-GitHub-User-ID": String(session.user.githubUserId),
  };

  const [streamResult, rhythmResult] = await Promise.all([
    client.GET("/api/activity/stream", { headers }),
    client.GET("/api/activity/rhythm", { headers }),
  ]);

  if (streamResult.error) {
    fetchError = streamResult.error.error;
  } else {
    streamData = streamResult.data;
  }

  if (rhythmResult.error) {
    fetchError = rhythmResult.error.error;
  } else {
    rhythmData = rhythmResult.data;
  }

  return (
    <ActivityPresenter
      streamData={streamData}
      rhythmData={rhythmData}
      initialError={fetchError}
    />
  );
}
