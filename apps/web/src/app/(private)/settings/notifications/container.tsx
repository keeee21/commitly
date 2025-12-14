import { auth } from "@/config/auth.config";
import { NotificationsPresenter } from "./presenter";

export async function NotificationsContainer() {
  const session = await auth();

  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  return <NotificationsPresenter session={session} />;
}
