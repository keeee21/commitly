"use client";

import Image from "next/image";
import { useActionState, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { components } from "@/lib/api/schema";
import { type ActionState, addRival, removeRival } from "./actions";

type Rival = components["schemas"]["dto.RivalResponse"];

type RivalsPresenterProps = {
  rivals: Rival[];
  count: number;
  maxRivals: number;
  initialError: string | null;
};

export function RivalsPresenter({
  rivals,
  count,
  maxRivals,
  initialError,
}: RivalsPresenterProps) {
  const [error, setError] = useState<string | null>(initialError);
  const [addError, setAddError] = useState<string | null>(null);

  const [, addRivalAction, isAddingPending] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await addRival(prevState, formData);
      if (result?.status === "ERROR") {
        setAddError(result.message);
      } else if (result?.status === "SUCCESS") {
        setAddError(null);
        // Clear the input by resetting the form
        const form = document.getElementById(
          "add-rival-form",
        ) as HTMLFormElement;
        form?.reset();
      }
      return result;
    },
    null,
  );

  const [, removeRivalAction, isRemovingPending] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await removeRival(prevState, formData);
      if (result?.status === "ERROR") {
        setError(result.message);
      } else if (result?.status === "SUCCESS") {
        setError(null);
      }
      return result;
    },
    null,
  );

  const handleRemove = (e: React.FormEvent<HTMLFormElement>) => {
    if (!confirm("このライバルを削除しますか？")) {
      e.preventDefault();
    }
  };

  const isPending = isAddingPending || isRemovingPending;

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Rivals</h1>
        <div className="text-sm text-zinc-500 dark:text-zinc-400">
          {count} / {maxRivals} rivals
        </div>
      </div>

      {/* ライバル追加フォーム */}
      <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
        <h2 className="text-lg font-semibold mb-4">Add a Rival</h2>
        <form id="add-rival-form" action={addRivalAction} className="space-y-4">
          <div className="flex gap-2">
            <Input
              type="text"
              name="username"
              placeholder="GitHub username"
              disabled={isPending || count >= maxRivals}
              className="flex-1"
            />
            <Button type="submit" disabled={isPending || count >= maxRivals}>
              {isAddingPending ? "Adding..." : "Add"}
            </Button>
          </div>
          {addError && (
            <p className="text-sm text-red-600 dark:text-red-400">{addError}</p>
          )}
          {count >= maxRivals && (
            <p className="text-sm text-amber-600 dark:text-amber-400">
              ライバル登録数が上限（{maxRivals}人）に達しています
            </p>
          )}
        </form>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      <div className="space-y-4">
        {rivals.length === 0 ? (
          <div className="bg-zinc-50 dark:bg-zinc-800 rounded-lg p-8 text-center">
            <p className="text-zinc-500 dark:text-zinc-400">
              ライバルがまだ登録されていません。
              <br />
              上のフォームからGitHubユーザー名を入力して追加してください。
            </p>
          </div>
        ) : (
          <div className="grid gap-4">
            {rivals.map((rival) => (
              <div
                key={rival.id}
                className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-4"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    {rival.avatar_url && (
                      <Image
                        src={rival.avatar_url}
                        alt={rival.github_username}
                        width={48}
                        height={48}
                        className="rounded-full"
                      />
                    )}
                    <div>
                      <div className="font-medium">{rival.github_username}</div>
                      <a
                        href={`https://github.com/${rival.github_username}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-sm text-blue-600 dark:text-blue-400 hover:underline"
                      >
                        View on GitHub
                      </a>
                    </div>
                  </div>
                  <form action={removeRivalAction} onSubmit={handleRemove}>
                    <input type="hidden" name="rivalId" value={rival.id} />
                    <Button
                      type="submit"
                      variant="outline"
                      size="sm"
                      disabled={isPending}
                      className="text-red-600 dark:text-red-400 border-red-200 dark:border-red-800 hover:bg-red-50 dark:hover:bg-red-900/20"
                    >
                      Remove
                    </Button>
                  </form>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
