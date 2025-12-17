/// <reference types="vitest" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  test: {
    environment: 'happy-dom', // or jsdom if installed
    globals: true,
    setupFiles: './src/test-setup.ts', // This allows using describe, it, expect without importing? Wait, I imported them in test file. But globals: true is common.
    // I imported them explicitly in my test file, so globals: true is not strictly necessary but harmless.
  }
})
