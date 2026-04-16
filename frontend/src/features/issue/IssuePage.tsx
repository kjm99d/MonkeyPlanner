import { useParams } from 'react-router-dom';

export default function IssuePage() {
  const { issueId } = useParams<{ issueId: string }>();
  return (
    <section>
      <h1 className="text-3xl font-bold">이슈 {issueId}</h1>
      <p className="mt-2 text-ink-secondary">상세 페이지는 US-M3-4에서 구현됩니다.</p>
    </section>
  );
}
