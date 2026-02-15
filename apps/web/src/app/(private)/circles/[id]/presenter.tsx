"use client";

import Image from "next/image";
import Link from "next/link";
import type { components } from "@/lib/api/schema";

type CircleResponse = components["schemas"]["dto.CircleResponse"];
type SignalResponse = components["schemas"]["dto.SignalResponse"];

type CircleDetailPresenterProps = {
  circle: CircleResponse | null;
  signals: SignalResponse[];
  initialError: string | null;
};

function getSignalIcon(type: string): string {
  switch (type) {
    case "same_day":
      return "🤝";
    case "same_hour":
      return "⏰";
    case "same_language":
      return "💻";
    default:
      return "✨";
  }
}

function getSignalMessage(signal: SignalResponse): string {
  const usernames = signal.users.map((u) => u.github_username).join("、");
  switch (signal.type) {
    case "same_day":
      return `${usernames}さんと同じ日にコミットしていました`;
    case "same_hour":
      return `${signal.detail}に${usernames}さんもコミットしていました`;
    case "same_language":
      return `${usernames}さんも${signal.detail}を使っていました`;
    default:
      return signal.detail;
  }
}

export function CircleDetailPresenter({
  circle,
  signals,
  initialError,
}: CircleDetailPresenterProps) {
  if (initialError) {
    return (
      <div className="p-6">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-md">
          {initialError}
        </div>
        <Link
          href="/circles"
          className="mt-4 inline-block text-sm text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200"
        >
          ← サークル一覧に戻る
        </Link>
      </div>
    );
  }

  if (!circle) return null;

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center gap-2">
        <Link
          href="/circles"
          className="text-sm text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200"
        >
          ← サークル一覧
        </Link>
      </div>

      <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
        <h1 className="text-2xl font-bold mb-2">{circle.name}</h1>
        <div className="flex items-center gap-2 mb-4">
          <span className="text-sm text-zinc-500 dark:text-zinc-400">
            招待コード:
          </span>
          <code className="text-sm bg-zinc-100 dark:bg-zinc-800 px-2 py-0.5 rounded font-mono">
            {circle.invite_code}
          </code>
        </div>

        <div className="flex items-center gap-2">
          <span className="text-sm text-zinc-500 dark:text-zinc-400">
            Members ({circle.members.length}):
          </span>
          <div className="flex -space-x-2">
            {circle.members.map((member) => (
              <div key={member.github_username} title={member.github_username}>
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

      <div className="space-y-4">
        <h2 className="text-lg font-semibold">並走シグナル</h2>
        {signals.length === 0 ? (
          <div className="bg-zinc-50 dark:bg-zinc-800 rounded-lg p-8 text-center">
            <p className="text-zinc-500 dark:text-zinc-400">
              まだ並走シグナルはありません。
              <br />
              メンバーがコミットすると検出されます。
            </p>
          </div>
        ) : (
          <div className="grid gap-3">
            {signals.map((signal, index) => (
              <div
                key={`${signal.type}-${signal.date}-${index}`}
                className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-4 flex items-start gap-3"
              >
                <span className="text-2xl">{getSignalIcon(signal.type)}</span>
                <div className="flex-1">
                  <p className="text-sm font-medium">
                    {getSignalMessage(signal)}
                  </p>
                  <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
                    {signal.date}
                  </p>
                </div>
                <div className="flex -space-x-1">
                  {signal.users.map((user) => (
                    <div
                      key={user.github_username}
                      title={user.github_username}
                    >
                      {user.avatar_url ? (
                        <Image
                          src={user.avatar_url}
                          alt={user.github_username}
                          width={24}
                          height={24}
                          className="rounded-full border border-white dark:border-zinc-900"
                        />
                      ) : (
                        <div className="w-6 h-6 rounded-full bg-zinc-300 dark:bg-zinc-600 border border-white dark:border-zinc-900 flex items-center justify-center text-xs">
                          {user.github_username.charAt(0)}
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
