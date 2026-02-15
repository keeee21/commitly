import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { DashboardPresenter } from "./presenter";

type DashboardData = components["schemas"]["usecase.DashboardData"];
type SignalResponse = components["schemas"]["dto.SignalResponse"];

type SearchParams = Promise<{
  period?: string;
}>;

type DashboardContainerProps = {
  searchParams: SearchParams;
};

export async function DashboardContainer({
  searchParams,
}: DashboardContainerProps) {
  const session = await auth();

  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  const params = await searchParams;
  const period = params.period === "monthly" ? "monthly" : "weekly";

  let dashboardData: DashboardData | null = null;
  let fetchError: string | null = null;

  const endpoint =
    period === "weekly" ? "/api/dashboard/weekly" : "/api/dashboard/monthly";

  const { data, error } = await client.GET(endpoint, {
    headers: {
      "X-GitHub-User-ID": String(session.user.githubUserId),
    },
  });

  if (error) {
    fetchError = error.error;
  } else {
    dashboardData = data;
  }

  // サークル情報とシグナルを取得
  const headers = {
    "X-GitHub-User-ID": String(session.user.githubUserId),
  };

  const circlesRes = await client.GET("/api/circles", { headers });
  const circles = circlesRes.data?.circles ?? [];
  const hasCircles = circles.length > 0;

  let signals: SignalResponse[] = [];
  if (hasCircles) {
    const signalsRes = await client.GET("/api/signals/recent", { headers });
    if (!signalsRes.error && signalsRes.data) {
      signals = signalsRes.data.signals;
    }
  }

  return (
    <DashboardPresenter
      period={period}
      dashboardData={dashboardData}
      initialError={fetchError}
      signals={signals}
      hasCircles={hasCircles}
    />
  );
}
