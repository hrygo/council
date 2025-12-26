/// <reference types="vitest" />
import { defineConfig } from 'vitest/config'
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
      '/ws': {
        target: 'http://localhost:8080',
        ws: true,
        changeOrigin: true,
      },
    },
  },
  build: {
    // Bundle optimization
    rollupOptions: {
      output: {
        manualChunks: {
          // React vendor bundle
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          // ReactFlow vendor bundle
          'flow-vendor': ['@xyflow/react'],
          // Markdown and math rendering
          'markdown': ['react-markdown', 'rehype-katex', 'remark-math', 'remark-gfm'],
          // UI component libraries
          'ui-vendor': ['lucide-react'],
          // State management
          'state': ['zustand'],
        },
      },
    },
    // Target modern browsers for smaller bundle
    target: 'es2020',
    // Enable minification
    minify: 'esbuild',
  },
  test: {
    environment: 'happy-dom', // or jsdom if installed
    globals: true,
    setupFiles: './src/test-setup.ts', // This allows using describe, it, expect without importing? Wait, I imported them in test file. But globals: true is common.
    // I imported them explicitly in my test file, so globals: true is not strictly necessary but harmless.
  }
})
