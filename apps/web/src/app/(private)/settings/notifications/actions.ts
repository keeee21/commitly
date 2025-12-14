"use server";

import { revalidatePath } from "next/cache";
import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";

export type ActionState = {
  status: "SUCCESS" | "ERROR";
  message: string;
} | null;

export async function createSlackNotification(
  _prevState: ActionState,
  formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const webhookUrl = formData.get("webhook_url");
  if (!webhookUrl || typeof webhookUrl !== "string") {
    return { status: "ERROR", message: "Webhook URLを入力してください" };
  }

  const { error } = await client.POST("/api/notifications/slack", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
    body: {
      webhook_url: webhookUrl,
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/settings/notifications");
  return { status: "SUCCESS", message: "Slack通知を設定しました" };
}

export async function updateSlackEnabled(
  _prevState: ActionState,
  formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const isEnabled = formData.get("is_enabled") === "true";

  const { error } = await client.PUT("/api/notifications/slack", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
    body: {
      is_enabled: isEnabled,
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/settings/notifications");
  return { status: "SUCCESS", message: "設定を更新しました" };
}

export async function disconnectSlack(
  _prevState: ActionState,
  _formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const { error } = await client.DELETE("/api/notifications/slack", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/settings/notifications");
  return { status: "SUCCESS", message: "Slackとの連携を解除しました" };
}
