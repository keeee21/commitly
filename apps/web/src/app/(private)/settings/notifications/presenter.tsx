"use client";

import type { Session } from "next-auth";
import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";

type NotificationsPresenterProps = {
  session: Session;
};

type NotificationSetting = components["schemas"]["NotificationSetting"];
type ChannelType = "line" | "slack" | "discord";

const channelLabels: Record<ChannelType, string> = {
  line: "LINE",
  slack: "Slack",
  discord: "Discord",
};

export function NotificationsPresenter({
  session,
}: NotificationsPresenterProps) {
  const [settings, setSettings] = useState<NotificationSetting[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [newChannel, setNewChannel] = useState<ChannelType>("slack");
  const [newWebhookUrl, setNewWebhookUrl] = useState("");
  const [newLineUserId, setNewLineUserId] = useState("");
  const [adding, setAdding] = useState(false);
  const [addError, setAddError] = useState<string | null>(null);

  const fetchSettings = useCallback(async () => {
    try {
      const { data, error: apiError } = await client.GET("/api/notifications", {
        headers: {
          "X-GitHub-User-ID": String(session.user.githubUserId),
        },
      });

      if (apiError) {
        setError(apiError.error);
        return;
      }

      setSettings(data.settings);
    } catch {
      setError("通知設定の取得に失敗しました");
    } finally {
      setLoading(false);
    }
  }, [session.user.githubUserId]);

  useEffect(() => {
    fetchSettings();
  }, [fetchSettings]);

  const handleAddSetting = async (e: React.FormEvent) => {
    e.preventDefault();

    if (newChannel !== "line" && !newWebhookUrl.trim()) {
      setAddError("Webhook URLを入力してください");
      return;
    }
    if (newChannel === "line" && !newLineUserId.trim()) {
      setAddError("LINE User IDを入力してください");
      return;
    }

    setAdding(true);
    setAddError(null);

    try {
      const { error: apiError } = await client.POST("/api/notifications", {
        headers: {
          "X-GitHub-User-ID": String(session.user.githubUserId),
        },
        body: {
          channel_type: newChannel,
          webhook_url: newChannel !== "line" ? newWebhookUrl.trim() : undefined,
          line_user_id:
            newChannel === "line" ? newLineUserId.trim() : undefined,
        },
      });

      if (apiError) {
        setAddError(apiError.error);
        return;
      }

      setNewWebhookUrl("");
      setNewLineUserId("");
      fetchSettings();
    } catch {
      setAddError("通知設定の追加に失敗しました");
    } finally {
      setAdding(false);
    }
  };

  const handleToggleEnabled = async (setting: NotificationSetting) => {
    try {
      await client.PUT("/api/notifications/{id}", {
        params: {
          path: { id: setting.id },
        },
        headers: {
          "X-GitHub-User-ID": String(session.user.githubUserId),
        },
        body: {
          is_enabled: !setting.is_enabled,
          webhook_url: setting.webhook_url,
          line_user_id: setting.line_user_id,
        },
      });

      fetchSettings();
    } catch {
      setError("通知設定の更新に失敗しました");
    }
  };

  const handleDeleteSetting = async (settingId: number) => {
    if (!confirm("この通知設定を削除しますか？")) return;

    try {
      await client.DELETE("/api/notifications/{id}", {
        params: {
          path: { id: settingId },
        },
        headers: {
          "X-GitHub-User-ID": String(session.user.githubUserId),
        },
      });

      fetchSettings();
    } catch {
      setError("通知設定の削除に失敗しました");
    }
  };

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">Notification Settings</h1>

      <p className="text-zinc-500 dark:text-zinc-400">
        週次レポートの通知先を設定します。毎週月曜日の朝9時（JST）に配信されます。
      </p>

      {/* 通知設定追加フォーム */}
      <div className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
        <h2 className="text-lg font-semibold mb-4">Add Notification Channel</h2>
        <form onSubmit={handleAddSetting} className="space-y-4">
          <div>
            <span className="block text-sm font-medium mb-2">Channel Type</span>
            <div className="flex gap-2">
              {(["slack", "discord", "line"] as ChannelType[]).map(
                (channel) => (
                  <button
                    key={channel}
                    type="button"
                    onClick={() => setNewChannel(channel)}
                    className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                      newChannel === channel
                        ? "bg-zinc-900 text-white dark:bg-zinc-50 dark:text-zinc-900"
                        : "bg-zinc-100 text-zinc-900 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-50 dark:hover:bg-zinc-700"
                    }`}
                  >
                    {channelLabels[channel]}
                  </button>
                ),
              )}
            </div>
          </div>

          {newChannel !== "line" ? (
            <div>
              <label
                htmlFor="webhook-url"
                className="block text-sm font-medium mb-2"
              >
                Webhook URL
              </label>
              <Input
                id="webhook-url"
                type="url"
                value={newWebhookUrl}
                onChange={(e) => setNewWebhookUrl(e.target.value)}
                placeholder={`https://hooks.${newChannel}.com/...`}
                disabled={adding}
              />
              <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
                {newChannel === "slack"
                  ? "Slack Incoming Webhooks URL を入力してください"
                  : "Discord Webhook URL を入力してください"}
              </p>
            </div>
          ) : (
            <div>
              <label
                htmlFor="line-user-id"
                className="block text-sm font-medium mb-2"
              >
                LINE User ID
              </label>
              <Input
                id="line-user-id"
                type="text"
                value={newLineUserId}
                onChange={(e) => setNewLineUserId(e.target.value)}
                placeholder="U1234567890abcdef..."
                disabled={adding}
              />
              <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
                LINE公式アカウントから取得したUser IDを入力してください
              </p>
            </div>
          )}

          {addError && (
            <p className="text-sm text-red-600 dark:text-red-400">{addError}</p>
          )}

          <Button type="submit" disabled={adding}>
            {adding ? "Adding..." : "Add Channel"}
          </Button>
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
          <h2 className="text-lg font-semibold">Configured Channels</h2>
          {settings.length === 0 ? (
            <div className="bg-zinc-50 dark:bg-zinc-800 rounded-lg p-8 text-center">
              <p className="text-zinc-500 dark:text-zinc-400">
                通知チャンネルがまだ設定されていません。
                <br />
                上のフォームから追加してください。
              </p>
            </div>
          ) : (
            <div className="grid gap-4">
              {settings.map((setting) => (
                <div
                  key={setting.id}
                  className="bg-white dark:bg-zinc-900 rounded-lg border border-zinc-200 dark:border-zinc-800 p-4"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div
                        className={`w-10 h-10 rounded-full flex items-center justify-center text-white font-bold ${
                          setting.channel_type === "slack"
                            ? "bg-purple-500"
                            : setting.channel_type === "discord"
                              ? "bg-indigo-500"
                              : "bg-green-500"
                        }`}
                      >
                        {channelLabels[setting.channel_type][0]}
                      </div>
                      <div>
                        <div className="font-medium">
                          {channelLabels[setting.channel_type]}
                        </div>
                        <div className="text-sm text-zinc-500 dark:text-zinc-400 truncate max-w-[200px]">
                          {setting.channel_type === "line"
                            ? setting.line_user_id
                            : setting.webhook_url}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        onClick={() => handleToggleEnabled(setting)}
                        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                          setting.is_enabled
                            ? "bg-green-500"
                            : "bg-zinc-300 dark:bg-zinc-700"
                        }`}
                      >
                        <span
                          className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                            setting.is_enabled
                              ? "translate-x-6"
                              : "translate-x-1"
                          }`}
                        />
                      </button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleDeleteSetting(setting.id)}
                        className="text-red-600 dark:text-red-400 border-red-200 dark:border-red-800 hover:bg-red-50 dark:hover:bg-red-900/20"
                      >
                        Delete
                      </Button>
                    </div>
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
