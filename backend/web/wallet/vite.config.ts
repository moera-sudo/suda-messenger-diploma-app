import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  base: '/wallet/',
  resolve: {
    alias: {
      '@': '/src',
    },
  },
  server: {
    proxy: {
      '/api': { target: 'http://localhost', changeOrigin: true },
      '/ws': { target: 'ws://localhost', ws: true },
    },
  },
});
