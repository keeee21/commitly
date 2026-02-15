import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { RivalsPresenter } from "./presenter";

type Rival = components["schemas"]["dto.RivalResponse"];

export async function RivalsContainer() {
  const session = await auth();

  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  let rivals: Rival[] = [];
  let count = 0;
  let maxRivals = 5;
  let fetchError: string | null = null;

  const { data, error } = await client.GET("/api/rivals", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
  });

  if (error) {
    fetchError = error.error;
  } else {
    rivals = data.rivals;
    count = data.count;
    maxRivals = data.max_rivals;
  }

  return (
    <RivalsPresenter
      rivals={rivals}
      count={count}
      maxRivals={maxRivals}
      initialError={fetchError}
    />
  );
}
