import { useEffect } from 'react';
import { RouterProvider } from 'react-router-dom';
import { router } from './router';
import { useAuthStore } from './stores/authStore';
import { initAuth } from './lib/webapp';
import LoadingScreen from './components/LoadingScreen';
import ErrorPage from './pages/ErrorPage';

export default function App() {
  const status = useAuthStore((s) => s.status);
  const errorMsg = useAuthStore((s) => s.errorMsg);

  useEffect(() => {
    initAuth();
  }, []);

  if (status === 'idle' || status === 'authing') {
    return <LoadingScreen />;
  }

  if (status === 'error') {
    return <ErrorPage code={errorMsg || 'AUTH_ERROR'} />;
  }

  return <RouterProvider router={router} />;
}
