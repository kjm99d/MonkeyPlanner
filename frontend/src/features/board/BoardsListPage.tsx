import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Plus } from 'lucide-react';
import { useBoards, useCreateBoard, useIssues } from '../../api/hooks';
import { Button } from '../../components/Button';
import { Card } from '../../components/Card';
import { Input } from '../../components/Input';
import { Skeleton } from '../../components/Skeleton';

function BoardCardWithCount({ id, name, viewType, createdAt }: { id: string; name: string; viewType: string; createdAt: string }) {
  const issues = useIssues({ boardId: id });
  const total = issues.data?.length ?? 0;
  const done = issues.data?.filter((i) => i.status === 'Done').length ?? 0;
  const pct = total > 0 ? Math.round((done / total) * 100) : 0;

  return (
    <Card className="relative overflow-hidden">
      <h2 className="text-lg font-semibold">{name}</h2>
      <p className="mt-1 text-xs text-ink-muted">
        {new Date(createdAt).toLocaleDateString('ko-KR')} 생성 · {viewType}
      </p>
      {total > 0 && (
        <div className="mt-3 flex flex-col gap-1.5">
          <div className="flex items-center justify-between text-xs">
            <span className="text-ink-secondary">{done}/{total} 완료</span>
            <span className="font-medium text-ink-primary">{pct}%</span>
          </div>
          <div className="h-1.5 w-full overflow-hidden rounded-full bg-surface-muted">
            <div
              className="h-full rounded-full bg-status-done transition-all duration-500"
              style={{ width: `${pct}%` }}
            />
          </div>
        </div>
      )}
      {total === 0 && (
        <p className="mt-3 text-xs text-ink-muted">이슈 없음</p>
      )}
    </Card>
  );
}

export default function BoardsListPage() {
  const [name, setName] = useState('');
  const boards = useBoards();
  const createBoard = useCreateBoard();

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    await createBoard.mutateAsync({ name: name.trim() });
    setName('');
  };

  return (
    <section className="flex flex-col gap-6">
      <header className="fade-up flex items-end justify-between">
        <div>
          <h1 className="text-3xl font-bold">보드</h1>
          <p className="mt-1 text-ink-secondary">작업을 묶어둘 컨테이너입니다.</p>
        </div>
      </header>

      <form onSubmit={submit} className="fade-up flex gap-2" style={{ animationDelay: '60ms' }}>
        <Input
          placeholder="새 보드 이름"
          value={name}
          onChange={(e) => setName(e.target.value)}
          aria-label="새 보드 이름"
          className="flex-1"
        />
        <Button type="submit" disabled={createBoard.isPending}>
          <Plus size={16} /> 생성
        </Button>
      </form>

      {boards.isLoading && (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <Skeleton className="h-32 rounded-lg" />
          <Skeleton className="h-32 rounded-lg" />
          <Skeleton className="h-32 rounded-lg" />
        </div>
      )}
      {boards.data && boards.data.length === 0 && (
        <div className="fade-up flex flex-col items-center gap-4 py-16 text-center" style={{ animationDelay: '120ms' }}>
          <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-brand-500/10">
            <Plus size={28} className="text-brand-500" />
          </div>
          <div>
            <p className="text-lg font-semibold text-ink-primary">아직 보드가 없습니다</p>
            <p className="mt-1 text-sm text-ink-secondary">위 입력란에 이름을 적고 "생성"을 눌러 첫 보드를 만들어 보세요.</p>
          </div>
        </div>
      )}
      <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {boards.data?.map((b, i) => (
          <li key={b.id} className="fade-up" style={{ animationDelay: `${120 + i * 50}ms` }}>
            <Link
              to={`/boards/${b.id}`}
              className="block rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
            >
              <BoardCardWithCount id={b.id} name={b.name} viewType={b.viewType} createdAt={b.createdAt} />
            </Link>
          </li>
        ))}
      </ul>
    </section>
  );
}
