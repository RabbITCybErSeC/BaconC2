import { defineConfig, type UserConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig((): UserConfig => {
  return {
    plugins: [
      react(),
      tailwindcss(),
    ],
    server: {
      port: 5173,
      host: true,
      watch: {
        usePolling: true,
      },
    },
    esbuild: {
      target: 'esnext',
    }
  };
});
