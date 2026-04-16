import { Routes, Route, Link } from 'react-router-dom';

function Home() {
  return (
    <main className="p-6">
      <h1 className="text-3xl font-bold">몽키 플래너</h1>
      <p className="mt-2 text-sm text-neutral-600 dark:text-neutral-300">
        코숭이 에이전트들의 작업 기억 저장소
      </p>
      <nav className="mt-6 flex gap-4">
        <Link to="/boards" className="underline">보드</Link>
        <Link to="/calendar" className="underline">캘린더</Link>
      </nav>
    </main>
  );
}

function BoardsPlaceholder() {
  return <main className="p-6"><h1 className="text-2xl">보드 (M3-3에서 구현)</h1></main>;
}

function CalendarPlaceholder() {
  return <main className="p-6"><h1 className="text-2xl">캘린더 (M3-5에서 구현)</h1></main>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/boards" element={<BoardsPlaceholder />} />
      <Route path="/calendar" element={<CalendarPlaceholder />} />
    </Routes>
  );
}
