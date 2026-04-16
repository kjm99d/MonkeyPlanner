import { useParams } from 'react-router-dom';

export default function BoardPage() {
  const { boardId } = useParams<{ boardId: string }>();
  return (
    <section>
      <h1 className="text-3xl font-bold">보드 {boardId}</h1>
      <p className="mt-2 text-ink-secondary">칸반 뷰는 다음 스토리(US-M3-3)에서 구현됩니다.</p>
    </section>
  );
}
