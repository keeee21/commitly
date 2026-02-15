import { CircleDetailContainer } from "./container";

type PageProps = {
  params: Promise<{ id: string }>;
};

export default async function CircleDetailPage({ params }: PageProps) {
  return <CircleDetailContainer params={params} />;
}
