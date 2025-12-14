"use server";

import { revalidatePath } from "next/cache";
import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";

export type ActionState = {
  status: "SUCCESS" | "ERROR";
  message: string;
} | null;

export async function addRival(
  _prevState: ActionState,
  formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const username = formData.get("username") as string;
  if (!username?.trim()) {
    return { status: "ERROR", message: "ユーザー名を入力してください" };
  }

  const { error } = await client.POST("/api/rivals", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
    body: {
      username: username.trim(),
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/rivals");
  return { status: "SUCCESS", message: "ライバルを追加しました" };
}

export async function removeRival(
  _prevState: ActionState,
  formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const rivalId = Number(formData.get("rivalId"));
  if (!rivalId) {
    return { status: "ERROR", message: "無効なライバルIDです" };
  }

  const { error } = await client.DELETE("/api/rivals/{id}", {
    params: {
      path: { id: rivalId },
    },
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/rivals");
  return { status: "SUCCESS", message: "ライバルを削除しました" };
}
