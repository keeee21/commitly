import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { NotificationsPresenter } from "./presenter";

type SlackNotificationSetting =
  components["schemas"]["dto.SlackNotificationSettingResponse"];

export async function NotificationsContainer() {
  const session = await auth();

  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  // Fetch Slack setting
  let slackSetting: SlackNotificationSetting | null = null;
  let fetchError: string | null = null;

  const { data, error, response } = await client.GET(
    "/api/notifications/slack",
    {
      headers: {
        "X-GitHub-User-ID": String(session.user.githubUserId),
      },
    },
  );

  if (response.status === 404) {
    // Setting not found - user hasn't connected yet
    slackSetting = null;
  } else if (error) {
    fetchError = error.error;
  } else {
    slackSetting = data;
  }

  return (
    <NotificationsPresenter
      slackSetting={slackSetting}
      initialSuccessMessage={null}
      initialErrorMessage={fetchError}
    />
  );
}
