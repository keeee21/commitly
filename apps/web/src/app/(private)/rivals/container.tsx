import { auth } from "@/config/auth.config";
import { RivalsPresenter } from "./presenter";

export async function RivalsContainer() {
  const session = await auth();

  if (!session) {
    throw new Error("Unauthorized: Session not found");
  }

  return <RivalsPresenter session={session} />;
}
