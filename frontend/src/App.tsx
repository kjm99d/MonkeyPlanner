import { Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import HomePage from './features/home/HomePage';
import BoardsListPage from './features/board/BoardsListPage';
import BoardPage from './features/board/BoardPage';
import IssuePage from './features/issue/IssuePage';
import CalendarPage from './features/calendar/CalendarPage';
import ApprovalPage from './features/approval/ApprovalPage';

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/boards" element={<BoardsListPage />} />
        <Route path="/boards/:boardId" element={<BoardPage />} />
        <Route path="/issues/:issueId" element={<IssuePage />} />
        <Route path="/calendar" element={<CalendarPage />} />
        <Route path="/approve" element={<ApprovalPage />} />
        <Route path="*" element={<NotFound />} />
      </Route>
    </Routes>
  );
}

function NotFound() {
  return (
    <section>
      <h1 className="text-3xl font-bold">찾을 수 없습니다</h1>
      <p className="mt-2 text-ink-secondary">요청하신 경로가 존재하지 않습니다.</p>
    </section>
  );
}
