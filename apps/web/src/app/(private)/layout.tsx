import { AppSidebar } from "@/components/app-sidebar";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";

export default function PrivateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-14 items-center gap-4 border-b px-4">
          <SidebarTrigger />
        </header>
        <main className="flex-1 p-4">{children}</main>
      </SidebarInset>
    </SidebarProvider>
  );
}
