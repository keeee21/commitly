import { DashboardContainer } from "./container";

type PageProps = {
  searchParams: Promise<{
    period?: string;
  }>;
};

export default function DashboardPage({ searchParams }: PageProps) {
  return <DashboardContainer searchParams={searchParams} />;
}
