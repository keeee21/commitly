"use server";

import { revalidatePath } from "next/cache";
import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";

export type ActionState = {
  status: "SUCCESS" | "ERROR";
  message: string;
} | null;

export async function createCircle(
  _prevState: ActionState,
  formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const name = formData.get("name") as string;
  if (!name?.trim()) {
    return { status: "ERROR", message: "サークル名を入力してください" };
  }

  const { error } = await client.POST("/api/circles", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
    body: {
      name: name.trim(),
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/circles");
  return { status: "SUCCESS", message: "サークルを作成しました" };
}

export async function joinCircle(
  _prevState: ActionState,
  formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const inviteCode = formData.get("inviteCode") as string;
  if (!inviteCode?.trim()) {
    return { status: "ERROR", message: "招待コードを入力してください" };
  }

  const { error } = await client.POST("/api/circles/join", {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
    body: {
      invite_code: inviteCode.trim(),
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/circles");
  return { status: "SUCCESS", message: "サークルに参加しました" };
}

export async function leaveCircle(
  _prevState: ActionState,
  formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const circleId = Number(formData.get("circleId"));
  if (!circleId) {
    return { status: "ERROR", message: "無効なサークルIDです" };
  }

  const { error } = await client.DELETE("/api/circles/{id}/leave", {
    params: {
      path: { id: circleId },
    },
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/circles");
  return { status: "SUCCESS", message: "サークルを退会しました" };
}

export async function deleteCircle(
  _prevState: ActionState,
  formData: FormData,
): Promise<ActionState> {
  const session = await auth();
  if (!session) {
    return { status: "ERROR", message: "認証されていません" };
  }

  const circleId = Number(formData.get("circleId"));
  if (!circleId) {
    return { status: "ERROR", message: "無効なサークルIDです" };
  }

  const { error } = await client.DELETE("/api/circles/{id}", {
    params: {
      path: { id: circleId },
    },
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
  });

  if (error) {
    return { status: "ERROR", message: error.error };
  }

  revalidatePath("/circles");
  return { status: "SUCCESS", message: "サークルを削除しました" };
}
