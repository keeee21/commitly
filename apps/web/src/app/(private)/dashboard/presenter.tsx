"use client";

import Image from "next/image";
import Link from "next/link";
import type { components } from "@/lib/api/schema";

function formatMonth(dateStr: string): string {
  const date = new Date(dateStr);
  return `${date.getFullYear()}年${date.getMonth() + 1}月`;
}

type DashboardData = components["schemas"]["usecase.DashboardData"];

type SignalResponse = components["schemas"]["dto.SignalResponse"];

type DashboardPresenterProps = {
  period: "weekly" | "monthly";
  dashboardData: DashboardData | null;
  initialError: string | null;
  signals: SignalResponse[];
  hasCircles: boolean;
};

export function DashboardPresenter({
  period,
  dashboardData,
  initialError,
  signals,
  hasCircles,
}: DashboardPresenterProps) {
  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <div className="flex gap-2">
          <Link
            href="/dashboard?period=weekly"
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              period === "weekly"
                ? "bg-zinc-900 text-white dark:bg-zinc-50 dark:text-zinc-900"
                : "bg-zinc-100 text-zinc-900 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-50 dark:hover:bg-zinc-700"
            }`}
          >
            Weekly
          </Link>
          <Link
            href="/dashboard?period=monthly"
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              period === "monthly"
                ? "bg-zinc-900 text-white dark:bg-zinc-50 dark:text-zinc-900"
                : "bg-zinc-100 text-zinc-900 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-50 dark:hover:bg-zinc-700"
            }`}
          >
            Monthly
          </Link>
        </div>
      </div>

      {initialError && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-md">
          {initialError}
        </div>
      )}

      {dashboardData && (
        <div className="space-y-6">
          <div className="text-sm text-zinc-500 dark:text-zinc-400">
            {period === "monthly"
              ? formatMonth(dashboardData.start_date)
              : `${dashboardData.start_date} 〜 ${dashboardData.end_date}`}
          </div>

          {/* 自分のコミット統計 */}
          <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
            <h2 className="text-lg font-semibold mb-4">Your Commits</h2>
            <UserStatsCard stats={dashboardData.my_stats} isMe />
          </div>

          {/* ライバルとの比較 */}
          {dashboardData.rivals.length > 0 ? (
            <div className="space-y-4">
              <h2 className="text-lg font-semibold">Rivals</h2>
              <div className="grid gap-4">
                {dashboardData.rivals.map((rivalStats) => (
                  <div
                    key={rivalStats.github_user_id}
                    className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6"
                  >
                    <UserStatsCard stats={rivalStats} />
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <div className="bg-zinc-50 dark:bg-zinc-800 rounded-lg p-8 text-center">
              <p className="text-zinc-500 dark:text-zinc-400 mb-4">
                ライバルがまだ登録されていません
              </p>
              <Link
                href="/rivals"
                className="inline-block px-4 py-2 bg-zinc-900 text-white dark:bg-zinc-50 dark:text-zinc-900 rounded-md text-sm font-medium hover:opacity-80 transition-opacity"
              >
                ライバルを登録する
              </Link>
            </div>
          )}

          {/* 比較表 */}
          {dashboardData.rivals.length > 0 && (
            <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 overflow-hidden">
              <h2 className="text-lg font-semibold p-6 pb-4">
                Commit Comparison
              </h2>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800">
                      <th className="text-left px-6 py-3 text-sm font-medium text-zinc-500 dark:text-zinc-400">
                        User
                      </th>
                      <th className="text-right px-6 py-3 text-sm font-medium text-zinc-500 dark:text-zinc-400">
                        Total Commits
                      </th>
                      <th className="text-right px-6 py-3 text-sm font-medium text-zinc-500 dark:text-zinc-400">
                        Difference
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-zinc-200 dark:border-zinc-800 bg-blue-50 dark:bg-blue-900/20">
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-3">
                          {dashboardData.my_stats.avatar_url && (
                            <Image
                              src={dashboardData.my_stats.avatar_url}
                              alt={dashboardData.my_stats.github_username}
                              width={32}
                              height={32}
                              className="rounded-full"
                            />
                          )}
                          <span className="font-medium">
                            {dashboardData.my_stats.github_username} (You)
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-right font-mono text-lg font-semibold">
                        {dashboardData.my_stats.total_commits}
                      </td>
                      <td className="px-6 py-4 text-right">-</td>
                    </tr>
                    {dashboardData.rivals.map((rivalStats) => {
                      const diff =
                        dashboardData.my_stats.total_commits -
                        rivalStats.total_commits;
                      return (
                        <tr
                          key={rivalStats.github_user_id}
                          className="border-b border-zinc-200 dark:border-zinc-800"
                        >
                          <td className="px-6 py-4">
                            <div className="flex items-center gap-3">
                              {rivalStats.avatar_url && (
                                <Image
                                  src={rivalStats.avatar_url}
                                  alt={rivalStats.github_username}
                                  width={32}
                                  height={32}
                                  className="rounded-full"
                                />
                              )}
                              <span>{rivalStats.github_username}</span>
                            </div>
                          </td>
                          <td className="px-6 py-4 text-right font-mono text-lg">
                            {rivalStats.total_commits}
                          </td>
                          <td className="px-6 py-4 text-right">
                            <span
                              className={`font-mono ${
                                diff > 0
                                  ? "text-green-600 dark:text-green-400"
                                  : diff < 0
                                    ? "text-red-600 dark:text-red-400"
                                    : "text-zinc-500"
                              }`}
                            >
                              {diff > 0 ? "+" : ""}
                              {diff}
                            </span>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {/* サークルシグナル or 招待バナー */}
          {hasCircles ? (
            <div className="space-y-4">
              <h2 className="text-lg font-semibold">Circle Signals</h2>
              {signals.length > 0 ? (
                <div className="grid gap-3">
                  {signals.map((signal, index) => (
                    <div
                      key={`${signal.type}-${signal.date}-${index}`}
                      className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-4 flex items-start gap-3"
                    >
                      <span className="text-2xl">
                        {signal.type === "same_day"
                          ? "🤝"
                          : signal.type === "same_hour"
                            ? "⏰"
                            : "💻"}
                      </span>
                      <div className="flex-1">
                        <p className="text-sm font-medium">
                          {signal.users
                            .map((u) => u.github_username)
                            .join("、")}
                          さんと
                          {signal.type === "same_day"
                            ? "同じ日にコミット"
                            : signal.type === "same_hour"
                              ? `${signal.detail}にコミット`
                              : `${signal.detail}を使用`}
                        </p>
                        <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
                          {signal.circle_name} · {signal.date}
                        </p>
                      </div>
                    </div>
                  ))}
                  <Link
                    href="/circles"
                    className="text-sm text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200"
                  >
                    サークルを見る →
                  </Link>
                </div>
              ) : (
                <div className="bg-zinc-50 dark:bg-zinc-800 rounded-lg p-6 text-center">
                  <p className="text-zinc-500 dark:text-zinc-400">
                    まだ並走シグナルはありません。メンバーがコミットすると検出されます。
                  </p>
                </div>
              )}
            </div>
          ) : (
            <div className="bg-gradient-to-r from-indigo-50 to-purple-50 dark:from-indigo-900/20 dark:to-purple-900/20 rounded-lg border border-indigo-200 dark:border-indigo-800 p-6 text-center">
              <p className="text-lg font-semibold mb-2">
                サークルに参加して、仲間との並走シグナルを受け取ろう
              </p>
              <p className="text-sm text-zinc-500 dark:text-zinc-400 mb-4">
                同じ日・同じ時間帯・同じ言語でコミットすると、偶然の一致がシグナルとして届きます
              </p>
              <Link
                href="/circles"
                className="inline-block px-4 py-2 bg-zinc-900 text-white dark:bg-zinc-50 dark:text-zinc-900 rounded-md text-sm font-medium hover:opacity-80 transition-opacity"
              >
                サークルを作成する
              </Link>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

type UserStatsCardProps = {
  stats: components["schemas"]["usecase.UserCommitStats"];
  isMe?: boolean;
};

function UserStatsCard({ stats, isMe = false }: UserStatsCardProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4">
        {stats.avatar_url && (
          <Image
            src={stats.avatar_url}
            alt={stats.github_username}
            width={48}
            height={48}
            className="rounded-full"
          />
        )}
        <div>
          <div className="font-medium">
            {stats.github_username}
            {isMe && (
              <span className="ml-2 text-xs bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 px-2 py-0.5 rounded">
                You
              </span>
            )}
          </div>
          <div className="text-3xl font-bold mt-1">
            {stats.total_commits}
            <span className="text-sm font-normal text-zinc-500 dark:text-zinc-400 ml-2">
              commits
            </span>
          </div>
        </div>
      </div>

      {/* 日別コミット */}
      {stats.daily_stats.length > 0 && (
        <div>
          <h4 className="text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
            Daily Commits
          </h4>
          <div className="flex gap-1 flex-wrap">
            {stats.daily_stats.map((day) => (
              <div
                key={day.date}
                className="text-center"
                title={`${day.date}: ${day.commit_count} commits`}
              >
                <div
                  className={`w-8 h-8 rounded text-xs flex items-center justify-center ${
                    day.commit_count === 0
                      ? "bg-zinc-100 dark:bg-zinc-800 text-zinc-400"
                      : day.commit_count < 5
                        ? "bg-green-200 dark:bg-green-900 text-green-800 dark:text-green-200"
                        : day.commit_count < 10
                          ? "bg-green-400 dark:bg-green-700 text-white"
                          : "bg-green-600 dark:bg-green-500 text-white"
                  }`}
                >
                  {day.commit_count}
                </div>
                <div className="text-xs text-zinc-400 mt-1">
                  {new Date(day.date).getDate()}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* リポジトリ別コミット */}
      {stats.repo_stats.length > 0 && (
        <div>
          <h4 className="text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
            By Repository
          </h4>
          <div className="space-y-2">
            {stats.repo_stats.slice(0, 5).map((repo) => (
              <div key={repo.repository} className="flex justify-between">
                <span className="text-sm truncate max-w-[200px]">
                  {repo.repository}
                </span>
                <span className="text-sm font-mono">
                  {repo.commit_count} commits
                </span>
              </div>
            ))}
            {stats.repo_stats.length > 5 && (
              <div className="text-sm text-zinc-400">
                and {stats.repo_stats.length - 5} more...
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
