"use client";

import { Bell, GitCommit, Github, TrendingUp, Users } from "lucide-react";
import { signInWithGitHub } from "@/lib/auth";

export function LandingPresenter() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-zinc-50 to-zinc-100 dark:from-zinc-950 dark:to-zinc-900">
      {/* Hero Section */}
      <div className="mx-auto max-w-6xl px-4 py-16 sm:px-6 lg:px-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50 sm:text-5xl md:text-6xl">
            <span className="block">Commitly</span>
            <span className="mt-2 block text-2xl font-medium text-zinc-600 dark:text-zinc-400 sm:text-3xl">
              Track your commits, compete with friends
            </span>
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-lg text-zinc-600 dark:text-zinc-400">
            GitHubのコミット活動を可視化し、仲間と競い合いながらモチベーションを高めましょう。
            毎日のコミット数を追跡し、ライバルとの比較で成長を実感できます。
          </p>
          <div className="mt-10">
            <form action={signInWithGitHub}>
              <button
                type="submit"
                className="inline-flex items-center gap-3 rounded-lg bg-zinc-900 px-8 py-4 text-lg font-semibold text-white shadow-lg transition-all hover:bg-zinc-800 hover:shadow-xl dark:bg-zinc-50 dark:text-zinc-900 dark:hover:bg-zinc-200"
              >
                <Github className="h-6 w-6" />
                GitHubでログイン
              </button>
            </form>
          </div>
        </div>

        {/* Features Section */}
        <div className="mt-24">
          <h2 className="text-center text-2xl font-bold text-zinc-900 dark:text-zinc-50 sm:text-3xl">
            Features
          </h2>
          <div className="mt-12 grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
            <FeatureCard
              icon={<GitCommit className="h-8 w-8" />}
              title="コミット追跡"
              description="GitHubのコミット活動を自動で取得・集計。日々の進捗が一目でわかります。"
            />
            <FeatureCard
              icon={<Users className="h-8 w-8" />}
              title="ライバル機能"
              description="友達やチームメンバーをライバルとして追加。お互いのコミット数を比較できます。"
            />
            <FeatureCard
              icon={<TrendingUp className="h-8 w-8" />}
              title="ダッシュボード"
              description="コミット数の推移やライバルとの比較をグラフで可視化。成長を実感できます。"
            />
            <FeatureCard
              icon={<Bell className="h-8 w-8" />}
              title="通知機能"
              description="LINEでコミット状況を通知。モチベーション維持をサポートします。"
            />
          </div>
        </div>

        {/* CTA Section */}
        <div className="mt-24 text-center">
          <div className="rounded-2xl bg-zinc-900 px-8 py-12 dark:bg-zinc-800">
            <h2 className="text-2xl font-bold text-white sm:text-3xl">
              今すぐ始めよう
            </h2>
            <p className="mx-auto mt-4 max-w-xl text-zinc-300">
              GitHubアカウントでログインするだけで、すぐにコミット追跡を開始できます。
            </p>
            <div className="mt-8">
              <form action={signInWithGitHub}>
                <button
                  type="submit"
                  className="inline-flex items-center gap-3 rounded-lg bg-white px-8 py-4 text-lg font-semibold text-zinc-900 shadow-lg transition-all hover:bg-zinc-100"
                >
                  <Github className="h-6 w-6" />
                  GitHubでログイン
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>

      {/* Footer */}
      <footer className="mt-16 border-t border-zinc-200 py-8 dark:border-zinc-800">
        <div className="mx-auto max-w-6xl px-4 text-center text-sm text-zinc-500 dark:text-zinc-400">
          <p>&copy; 2025 Commitly. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}

interface FeatureCardProps {
  icon: React.ReactNode;
  title: string;
  description: string;
}

function FeatureCard({ icon, title, description }: FeatureCardProps) {
  return (
    <div className="rounded-xl bg-white p-6 shadow-md transition-shadow hover:shadow-lg dark:bg-zinc-800">
      <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-zinc-100 text-zinc-900 dark:bg-zinc-700 dark:text-zinc-50">
        {icon}
      </div>
      <h3 className="mt-4 text-lg font-semibold text-zinc-900 dark:text-zinc-50">
        {title}
      </h3>
      <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
        {description}
      </p>
    </div>
  );
}
