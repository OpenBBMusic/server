import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { fileURLToPath, URL } from 'url';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 7012,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:9799',
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: [
      { find: '@', replacement: fileURLToPath(new URL('./src', import.meta.url)) },
      { find: '@wails', replacement: fileURLToPath(new URL('./wailsjs', import.meta.url)) },
    ],
  },
});
