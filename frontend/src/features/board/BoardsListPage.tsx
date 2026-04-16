import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useBoards, useCreateBoard } from '../../api/hooks';
import { Button } from '../../components/Button';
import { Card } from '../../components/Card';
import { Input } from '../../components/Input';

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
      <header className="flex items-end justify-between">
        <div>
          <h1 className="text-3xl font-bold">보드</h1>
          <p className="mt-1 text-ink-secondary">작업을 묶어둘 컨테이너입니다.</p>
        </div>
      </header>

      <form onSubmit={submit} className="flex gap-2">
        <Input
          placeholder="새 보드 이름"
          value={name}
          onChange={(e) => setName(e.target.value)}
          aria-label="새 보드 이름"
          className="flex-1"
        />
        <Button type="submit" disabled={createBoard.isPending}>
          생성
        </Button>
      </form>

      {boards.isLoading && <p className="text-ink-secondary">불러오는 중…</p>}
      {boards.data && boards.data.length === 0 && (
        <div className="flex flex-col items-center gap-4 py-16 text-center">
          <span className="text-5xl" aria-hidden>🐒</span>
          <div>
            <p className="text-lg font-semibold text-ink-primary">아직 보드가 없습니다</p>
            <p className="mt-1 text-sm text-ink-secondary">위 입력란에 이름을 적고 "생성"을 눌러 첫 보드를 만들어 보세요.</p>
          </div>
        </div>
      )}
      <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {boards.data?.map((b) => (
          <li key={b.id}>
            <Link
              to={`/boards/${b.id}`}
              className="block focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-500 rounded-lg"
            >
              <Card>
                <h2 className="text-lg font-semibold">{b.name}</h2>
                <p className="mt-1 text-xs text-ink-muted">
                  {new Date(b.createdAt).toLocaleDateString('ko-KR')} 생성 · {b.viewType}
                </p>
              </Card>
            </Link>
          </li>
        ))}
      </ul>
    </section>
  );
}
