"use client";

import Image from "next/image";
import Link from "next/link";
import { useActionState, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { components } from "@/lib/api/schema";
import {
  type ActionState,
  createCircle,
  deleteCircle,
  joinCircle,
  leaveCircle,
} from "./actions";

type CircleResponse = components["schemas"]["dto.CircleResponse"];

type CirclesPresenterProps = {
  circles: CircleResponse[];
  count: number;
  maxCircles: number;
  initialError: string | null;
};

export function CirclesPresenter({
  circles,
  count,
  maxCircles,
  initialError,
}: CirclesPresenterProps) {
  const [error, setError] = useState<string | null>(initialError);
  const [createError, setCreateError] = useState<string | null>(null);
  const [joinError, setJoinError] = useState<string | null>(null);

  const [, createAction, isCreatingPending] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await createCircle(prevState, formData);
      if (result?.status === "ERROR") {
        setCreateError(result.message);
      } else if (result?.status === "SUCCESS") {
        setCreateError(null);
        const form = document.getElementById(
          "create-circle-form",
        ) as HTMLFormElement;
        form?.reset();
      }
      return result;
    },
    null,
  );

  const [, joinAction, isJoiningPending] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await joinCircle(prevState, formData);
      if (result?.status === "ERROR") {
        setJoinError(result.message);
      } else if (result?.status === "SUCCESS") {
        setJoinError(null);
        const form = document.getElementById(
          "join-circle-form",
        ) as HTMLFormElement;
        form?.reset();
      }
      return result;
    },
    null,
  );

  const [, leaveAction, isLeavingPending] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await leaveCircle(prevState, formData);
      if (result?.status === "ERROR") {
        setError(result.message);
      } else if (result?.status === "SUCCESS") {
        setError(null);
      }
      return result;
    },
    null,
  );

  const [, deleteAction, isDeletingPending] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await deleteCircle(prevState, formData);
      if (result?.status === "ERROR") {
        setError(result.message);
      } else if (result?.status === "SUCCESS") {
        setError(null);
      }
      return result;
    },
    null,
  );

  const isPending =
    isCreatingPending ||
    isJoiningPending ||
    isLeavingPending ||
    isDeletingPending;

  const handleLeave = (e: React.FormEvent<HTMLFormElement>) => {
    if (!confirm("このサークルを退会しますか？")) {
      e.preventDefault();
    }
  };

  const handleDelete = (e: React.FormEvent<HTMLFormElement>) => {
    if (!confirm("このサークルを削除しますか？メンバー全員が退会されます。")) {
      e.preventDefault();
    }
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Circles</h1>
        <div className="text-sm text-zinc-500 dark:text-zinc-400">
          {count} / {maxCircles} circles
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {/* サークル作成フォーム */}
        <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
          <h2 className="text-lg font-semibold mb-4">Create Circle</h2>
          <form
            id="create-circle-form"
            action={createAction}
            className="space-y-4"
          >
            <div className="flex gap-2">
              <Input
                type="text"
                name="name"
                placeholder="サークル名"
                disabled={isPending || count >= maxCircles}
                className="flex-1"
              />
              <Button type="submit" disabled={isPending || count >= maxCircles}>
                {isCreatingPending ? "Creating..." : "Create"}
              </Button>
            </div>
            {createError && (
              <p className="text-sm text-red-600 dark:text-red-400">
                {createError}
              </p>
            )}
            {count >= maxCircles && (
              <p className="text-sm text-amber-600 dark:text-amber-400">
                サークル作成数が上限（{maxCircles}個）に達しています
              </p>
            )}
          </form>
        </div>

        {/* 招待コード参加フォーム */}
        <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
          <h2 className="text-lg font-semibold mb-4">Join Circle</h2>
          <form id="join-circle-form" action={joinAction} className="space-y-4">
            <div className="flex gap-2">
              <Input
                type="text"
                name="inviteCode"
                placeholder="招待コード"
                disabled={isPending}
                className="flex-1"
              />
              <Button type="submit" disabled={isPending}>
                {isJoiningPending ? "Joining..." : "Join"}
              </Button>
            </div>
            {joinError && (
              <p className="text-sm text-red-600 dark:text-red-400">
                {joinError}
              </p>
            )}
          </form>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      <div className="space-y-4">
        {circles.length === 0 ? (
          <div className="bg-zinc-50 dark:bg-zinc-800 rounded-lg p-8 text-center">
            <p className="text-zinc-500 dark:text-zinc-400">
              サークルがまだありません。
              <br />
              新しいサークルを作成するか、招待コードで参加してください。
            </p>
          </div>
        ) : (
          <div className="grid gap-4">
            {circles.map((circle) => (
              <div
                key={circle.id}
                className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6"
              >
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <Link
                      href={`/circles/${circle.id}`}
                      className="text-lg font-semibold hover:underline"
                    >
                      {circle.name}
                    </Link>
                    <div className="flex items-center gap-2 mt-1">
                      <span className="text-sm text-zinc-500 dark:text-zinc-400">
                        招待コード:
                      </span>
                      <code className="text-sm bg-zinc-100 dark:bg-zinc-800 px-2 py-0.5 rounded font-mono">
                        {circle.invite_code}
                      </code>
                      <CopyButton text={circle.invite_code} />
                    </div>
                  </div>
                  <div className="flex gap-2">
                    {circle.is_owner ? (
                      <form action={deleteAction} onSubmit={handleDelete}>
                        <input
                          type="hidden"
                          name="circleId"
                          value={circle.id}
                        />
                        <Button
                          type="submit"
                          variant="outline"
                          size="sm"
                          disabled={isPending}
                          className="text-red-600 dark:text-red-400 border-red-200 dark:border-red-800 hover:bg-red-50 dark:hover:bg-red-900/20"
                        >
                          Delete
                        </Button>
                      </form>
                    ) : (
                      <form action={leaveAction} onSubmit={handleLeave}>
                        <input
                          type="hidden"
                          name="circleId"
                          value={circle.id}
                        />
                        <Button
                          type="submit"
                          variant="outline"
                          size="sm"
                          disabled={isPending}
                        >
                          Leave
                        </Button>
                      </form>
                    )}
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <span className="text-sm text-zinc-500 dark:text-zinc-400">
                    Members ({circle.members.length}):
                  </span>
                  <div className="flex -space-x-2">
                    {circle.members.map((member) => (
                      <div
                        key={member.github_username}
                        title={member.github_username}
                      >
                        {member.avatar_url ? (
                          <Image
                            src={member.avatar_url}
                            alt={member.github_username}
                            width={32}
                            height={32}
                            className="rounded-full border-2 border-white dark:border-zinc-900"
                          />
                        ) : (
                          <div className="w-8 h-8 rounded-full bg-zinc-300 dark:bg-zinc-600 border-2 border-white dark:border-zinc-900 flex items-center justify-center text-xs">
                            {member.github_username.charAt(0)}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <button
      type="button"
      onClick={handleCopy}
      className="text-xs text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200"
    >
      {copied ? "Copied!" : "Copy"}
    </button>
  );
}
