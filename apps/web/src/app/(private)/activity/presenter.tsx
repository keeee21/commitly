"use client";

import Image from "next/image";
import type { components } from "@/lib/api/schema";

type ActivityStreamResponse =
  components["schemas"]["dto.ActivityStreamResponse"];
type RhythmResponse = components["schemas"]["dto.RhythmResponse"];
type UserRhythm = components["schemas"]["dto.UserRhythm"];

type ActivityPresenterProps = {
  streamData: ActivityStreamResponse | null;
  rhythmData: RhythmResponse | null;
  initialError: string | null;
};

export function ActivityPresenter({
  streamData,
  rhythmData,
  initialError,
}: ActivityPresenterProps) {
  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">Activity</h1>

      {initialError && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-md">
          {initialError}
        </div>
      )}

      {/* リズム可視化 */}
      {rhythmData && (
        <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
          <h2 className="text-lg font-semibold mb-1">Weekly Rhythm</h2>
          <p className="text-sm text-zinc-500 dark:text-zinc-400 mb-4">
            {rhythmData.period}
          </p>
          <div className="space-y-3">
            {rhythmData.users.map((user) => (
              <RhythmRow key={user.github_username} user={user} />
            ))}
          </div>
        </div>
      )}

      {/* アクティビティストリーム */}
      {streamData && (
        <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
          <h2 className="text-lg font-semibold mb-4">Activity Stream</h2>
          {streamData.activities.length > 0 ? (
            <div className="space-y-3">
              {streamData.activities.map((item, index) => (
                <div
                  key={`${item.github_username}-${item.repository}-${item.date}-${index}`}
                  className="flex items-center gap-3 py-2 border-b border-zinc-100 dark:border-zinc-800 last:border-0"
                >
                  {item.avatar_url && (
                    <Image
                      src={item.avatar_url}
                      alt={item.github_username}
                      width={32}
                      height={32}
                      className="rounded-full"
                    />
                  )}
                  <div className="flex-1 min-w-0">
                    <span className="font-medium text-sm">
                      {item.github_username}
                    </span>
                    <span className="text-zinc-400 text-sm mx-1">が</span>
                    <span className="text-sm text-zinc-600 dark:text-zinc-300 font-mono truncate">
                      {item.repository}
                    </span>
                    <span className="text-zinc-400 text-sm mx-1">に</span>
                    <span className="text-sm font-medium">
                      {item.commit_count} commits
                    </span>
                  </div>
                  <span className="text-xs text-zinc-400 whitespace-nowrap">
                    {item.date}
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-zinc-500 dark:text-zinc-400 text-center py-8">
              まだアクティビティがありません
            </p>
          )}
        </div>
      )}
    </div>
  );
}

const WEEKDAY_LABELS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

function RhythmRow({ user }: { user: UserRhythm }) {
  const days = [
    user.weekly_rhythm.mon,
    user.weekly_rhythm.tue,
    user.weekly_rhythm.wed,
    user.weekly_rhythm.thu,
    user.weekly_rhythm.fri,
    user.weekly_rhythm.sat,
    user.weekly_rhythm.sun,
  ];

  return (
    <div className="flex items-center gap-4">
      <div className="flex items-center gap-2 w-32 min-w-0">
        {user.avatar_url && (
          <Image
            src={user.avatar_url}
            alt={user.github_username}
            width={24}
            height={24}
            className="rounded-full"
          />
        )}
        <span className="text-sm font-medium truncate">
          {user.github_username}
        </span>
      </div>
      <div className="flex gap-1">
        {days.map((active, i) => (
          <div
            key={WEEKDAY_LABELS[i]}
            title={WEEKDAY_LABELS[i]}
            className={`w-8 h-8 rounded text-xs flex items-center justify-center ${
              active
                ? "bg-green-500 dark:bg-green-600 text-white"
                : "bg-zinc-100 dark:bg-zinc-800 text-zinc-400"
            }`}
          >
            {WEEKDAY_LABELS[i].charAt(0)}
          </div>
        ))}
      </div>
      <span className="text-xs text-zinc-500 dark:text-zinc-400 whitespace-nowrap">
        {user.pattern_label}
      </span>
    </div>
  );
}
