import { auth } from "@/config/auth.config";
import { DashboardPresenter } from "./presenter";

export async function DashboardContainer() {
  const session = await auth();

  // middlewareで認証チェックされているため、ここではsessionは必ず存在する
  // 型安全性のため、念のためnullチェックを行う
  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  return <DashboardPresenter session={session} />;
}
