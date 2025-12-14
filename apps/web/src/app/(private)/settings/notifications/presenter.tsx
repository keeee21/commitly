"use client";

import { useActionState, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { components } from "@/lib/api/schema";
import {
  type ActionState,
  createSlackNotification,
  disconnectSlack,
  updateSlackEnabled,
} from "./actions";

type SlackNotificationSetting =
  components["schemas"]["SlackNotificationSetting"];

type NotificationsPresenterProps = {
  slackSetting: SlackNotificationSetting | null;
  initialSuccessMessage: string | null;
  initialErrorMessage: string | null;
};

export function NotificationsPresenter({
  slackSetting,
  initialSuccessMessage,
  initialErrorMessage,
}: NotificationsPresenterProps) {
  const [successMessage, setSuccessMessage] = useState(initialSuccessMessage);
  const [errorMessage, setErrorMessage] = useState(initialErrorMessage);

  const [, createAction, isCreating] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await createSlackNotification(prevState, formData);
      if (result?.status === "ERROR") {
        setErrorMessage(result.message);
        setSuccessMessage(null);
      } else if (result?.status === "SUCCESS") {
        setSuccessMessage(result.message);
        setErrorMessage(null);
      }
      return result;
    },
    null,
  );

  const [, updateEnabledAction, isUpdating] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await updateSlackEnabled(prevState, formData);
      if (result?.status === "ERROR") {
        setErrorMessage(result.message);
        setSuccessMessage(null);
      } else if (result?.status === "SUCCESS") {
        setSuccessMessage(result.message);
        setErrorMessage(null);
      }
      return result;
    },
    null,
  );

  const [, disconnectAction, isDisconnecting] = useActionState(
    async (prevState: ActionState, formData: FormData) => {
      const result = await disconnectSlack(prevState, formData);
      if (result?.status === "ERROR") {
        setErrorMessage(result.message);
        setSuccessMessage(null);
      } else if (result?.status === "SUCCESS") {
        setSuccessMessage(result.message);
        setErrorMessage(null);
      }
      return result;
    },
    null,
  );

  const handleDisconnect = (e: React.FormEvent<HTMLFormElement>) => {
    if (!confirm("Slackとの連携を解除しますか？")) {
      e.preventDefault();
    }
  };

  const isPending = isCreating || isUpdating || isDisconnecting;

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">Notification Settings</h1>

      <p className="text-zinc-500 dark:text-zinc-400">
        週次レポートの通知先を設定します。毎週月曜日の朝9時（JST）に配信されます。
      </p>

      {successMessage && (
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 text-green-700 dark:text-green-400 px-4 py-3 rounded-md">
          {successMessage}
        </div>
      )}

      {errorMessage && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-md">
          {errorMessage}
        </div>
      )}

      <div className="space-y-6">
        {/* Slack Section */}
        <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 rounded-full bg-purple-500 flex items-center justify-center text-white font-bold">
              S
            </div>
            <h2 className="text-lg font-semibold">Slack</h2>
          </div>

          {slackSetting ? (
            <div className="space-y-4">
              <div className="bg-zinc-50 dark:bg-zinc-800 rounded-lg p-4">
                <div className="text-sm">
                  <span className="text-zinc-500 dark:text-zinc-400">
                    Webhook URL
                  </span>
                  <p className="font-medium font-mono text-xs mt-1">
                    {slackSetting.webhook_url}
                  </p>
                </div>
              </div>

              <div className="flex items-center justify-between">
                <form action={updateEnabledAction}>
                  <input
                    type="hidden"
                    name="is_enabled"
                    value={String(!slackSetting.is_enabled)}
                  />
                  <div className="flex items-center gap-3">
                    <span className="text-sm text-zinc-600 dark:text-zinc-400">
                      通知を有効にする
                    </span>
                    <button
                      type="submit"
                      disabled={isPending}
                      className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                        slackSetting.is_enabled
                          ? "bg-green-500"
                          : "bg-zinc-300 dark:bg-zinc-700"
                      }`}
                    >
                      <span
                        className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                          slackSetting.is_enabled
                            ? "translate-x-6"
                            : "translate-x-1"
                        }`}
                      />
                    </button>
                  </div>
                </form>

                <form action={disconnectAction} onSubmit={handleDisconnect}>
                  <Button
                    type="submit"
                    variant="outline"
                    size="sm"
                    disabled={isPending}
                    className="text-red-600 dark:text-red-400 border-red-200 dark:border-red-800 hover:bg-red-50 dark:hover:bg-red-900/20"
                  >
                    連携を解除
                  </Button>
                </form>
              </div>
            </div>
          ) : (
            <form action={createAction} className="space-y-4">
              <div className="space-y-2">
                <label
                  htmlFor="webhook_url"
                  className="text-sm font-medium text-zinc-700 dark:text-zinc-300"
                >
                  Webhook URL
                </label>
                <Input
                  id="webhook_url"
                  name="webhook_url"
                  type="url"
                  placeholder="https://hooks.slack.com/services/..."
                  required
                  disabled={isPending}
                  className="font-mono text-sm"
                />
                <p className="text-xs text-zinc-500 dark:text-zinc-400">
                  SlackのIncoming Webhooks URLを入力してください。
                  <a
                    href="https://api.slack.com/messaging/webhooks"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-500 hover:underline ml-1"
                  >
                    設定方法はこちら
                  </a>
                </p>
              </div>
              <Button type="submit" disabled={isPending}>
                {isCreating ? "設定中..." : "設定を保存"}
              </Button>
            </form>
          )}
        </div>

        {/* LINE Section - Coming Soon */}
        <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6 opacity-60">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 rounded-full bg-green-500 flex items-center justify-center text-white font-bold">
              L
            </div>
            <h2 className="text-lg font-semibold">LINE</h2>
            <span className="text-xs bg-zinc-200 dark:bg-zinc-700 px-2 py-1 rounded">
              Coming Soon
            </span>
          </div>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            LINE通知は近日対応予定です。
          </p>
        </div>

        {/* Discord Section - Coming Soon */}
        <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6 opacity-60">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 rounded-full bg-indigo-500 flex items-center justify-center text-white font-bold">
              D
            </div>
            <h2 className="text-lg font-semibold">Discord</h2>
            <span className="text-xs bg-zinc-200 dark:bg-zinc-700 px-2 py-1 rounded">
              Coming Soon
            </span>
          </div>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            Discord通知は近日対応予定です。
          </p>
        </div>
      </div>
    </div>
  );
}
