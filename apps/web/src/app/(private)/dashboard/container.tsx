import { auth } from "@/config/auth.config";
import { client } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { DashboardPresenter } from "./presenter";

type DashboardData = components["schemas"]["usecase.DashboardData"];

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

  return (
    <DashboardPresenter
      period={period}
      dashboardData={dashboardData}
      initialError={fetchError}
    />
  );
}
