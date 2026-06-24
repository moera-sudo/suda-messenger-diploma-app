import { createHashRouter } from 'react-router-dom';
import HomePage from './pages/HomePage';
import SendPage from './pages/SendPage';
import BuyPage from './pages/BuyPage';
import ReceivePage from './pages/ReceivePage';
import HistoryPage from './pages/HistoryPage';
import ErrorPage from './pages/ErrorPage';

export const router = createHashRouter([
  { path: '/', element: <HomePage /> },
  { path: '/send', element: <SendPage /> },
  { path: '/buy', element: <BuyPage /> },
  { path: '/receive', element: <ReceivePage /> },
  { path: '/history', element: <HistoryPage /> },
  { path: '*', element: <ErrorPage code="NOT_FOUND" /> },
]);
