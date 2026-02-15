import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { CircleDetailPresenter } from "./presenter";

type CircleResponse = components["schemas"]["dto.CircleResponse"];
type SignalResponse = components["schemas"]["dto.SignalResponse"];

type CircleDetailContainerProps = {
  params: Promise<{ id: string }>;
};

export async function CircleDetailContainer({
  params,
}: CircleDetailContainerProps) {
  const session = await auth();

  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  const { id } = await params;

  let circle: CircleResponse | null = null;
  let signals: SignalResponse[] = [];
  let fetchError: string | null = null;

  const circlesRes = await client.GET("/api/circles", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
  });

  if (circlesRes.error) {
    fetchError = circlesRes.error.error;
  } else {
    circle = circlesRes.data.circles.find((c) => String(c.id) === id) ?? null;
    if (!circle) {
      fetchError = "サークルが見つかりません";
    }
  }

  if (circle) {
    const signalsRes = await client.GET("/api/circles/{id}/signals", {
      params: { path: { id: Number(id) } },
      headers: {
        "X-GitHub-User-ID": String(session.user.githubUserId),
      },
    });

    if (!signalsRes.error && signalsRes.data) {
      signals = signalsRes.data.signals;
    }
  }

  return (
    <CircleDetailPresenter
      circle={circle}
      signals={signals}
      initialError={fetchError}
    />
  );
}
