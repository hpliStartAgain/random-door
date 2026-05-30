import { createBrowserRouter } from 'react-router-dom';
import { HomePage } from './pages/HomePage';
import { ModeSelectPage } from './pages/ModeSelectPage';
import { FreeExplorePage } from './pages/FreeExplorePage';
import { GameModePage } from './pages/GameModePage';
import { CityPage } from './pages/CityPage';
import { ChatPage } from './pages/ChatPage';
import { CheckinPage } from './pages/CheckinPage';
import { AchievementPage } from './pages/AchievementPage';

export const router = createBrowserRouter([
  { path: '/', element: <HomePage /> },
  { path: '/mode', element: <ModeSelectPage /> },
  { path: '/explore', element: <FreeExplorePage /> },
  { path: '/game', element: <GameModePage /> },
  { path: '/city/:id', element: <CityPage /> },
  { path: '/city/:id/chat/:cid', element: <ChatPage /> },
  { path: '/city/:id/checkin', element: <CheckinPage /> },
  { path: '/achievements', element: <AchievementPage /> },
]);
