"use client";

import Image from "next/image";
import type { Session } from "next-auth";
import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";

type RivalsPresenterProps = {
  session: Session;
};

type Rival = components["schemas"]["Rival"];

export function RivalsPresenter({ session }: RivalsPresenterProps) {
  const [rivals, setRivals] = useState<Rival[]>([]);
  const [count, setCount] = useState(0);
  const [maxRivals, setMaxRivals] = useState(5);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [username, setUsername] = useState("");
  const [adding, setAdding] = useState(false);
  const [addError, setAddError] = useState<string | null>(null);

  const fetchRivals = useCallback(async () => {
    try {
      const { data, error: apiError } = await client.GET("/api/rivals", {
        headers: {
          "X-GitHub-User-ID": String(session.user.githubUserId),
        },
      });

      if (apiError) {
        setError(apiError.error);
        return;
      }

      setRivals(data.rivals);
      setCount(data.count);
      setMaxRivals(data.max_rivals);
    } catch {
      setError("ライバル一覧の取得に失敗しました");
    } finally {
      setLoading(false);
    }
  }, [session.user.githubUserId]);

  useEffect(() => {
    fetchRivals();
  }, [fetchRivals]);

  const handleAddRival = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim()) return;

    setAdding(true);
    setAddError(null);

    try {
      const { error: apiError } = await client.POST("/api/rivals", {
        headers: {
          "X-GitHub-User-ID": String(session.user.githubUserId),
        },
        body: {
          username: username.trim(),
        },
      });

      if (apiError) {
        setAddError(apiError.error);
        return;
      }

      setUsername("");
      fetchRivals();
    } catch {
      setAddError("ライバルの追加に失敗しました");
    } finally {
      setAdding(false);
    }
  };

  const handleRemoveRival = async (rivalId: number) => {
    if (!confirm("このライバルを削除しますか？")) return;

    try {
      await client.DELETE("/api/rivals/{id}", {
        params: {
          path: { id: rivalId },
        },
        headers: {
          "X-GitHub-User-ID": String(session.user.githubUserId),
        },
      });

      fetchRivals();
    } catch {
      setError("ライバルの削除に失敗しました");
    }
  };

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
        <form onSubmit={handleAddRival} className="space-y-4">
          <div className="flex gap-2">
            <Input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="GitHub username"
              disabled={adding || count >= maxRivals}
              className="flex-1"
            />
            <Button
              type="submit"
              disabled={adding || !username.trim() || count >= maxRivals}
            >
              {adding ? "Adding..." : "Add"}
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

      {loading && (
        <div className="flex justify-center items-center h-32">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-zinc-900 dark:border-zinc-50" />
        </div>
      )}

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      {!loading && !error && (
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
                        <div className="font-medium">
                          {rival.github_username}
                        </div>
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
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleRemoveRival(rival.id)}
                      className="text-red-600 dark:text-red-400 border-red-200 dark:border-red-800 hover:bg-red-50 dark:hover:bg-red-900/20"
                    >
                      Remove
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
