import { redirect } from "next/navigation";
import { auth } from "@/config/auth.config";
import { ROUTES } from "@/constants/routes";
import { LandingPresenter } from "./presenter";

export default async function Home() {
  const session = await auth();

  if (session) {
    redirect(ROUTES.DASHBOARD);
  }

  return <LandingPresenter />;
}
