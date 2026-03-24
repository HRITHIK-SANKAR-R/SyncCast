import { defineConfig } from 'vite';

export default defineConfig({
  root: '.',
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/api': 'http://localhost:6969',
      '/ws': {
        target: 'ws://localhost:6969',
        ws: true,
      },
      '/media': 'http://localhost:6969',
    },
  },
});
