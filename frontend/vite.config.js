import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// SPA fallback: serve index.html for all HTML routes (/, /dashboard, etc.) in dev
function spaFallback() {
  return {
    name: 'spa-fallback',
    apply: 'serve',
    enforce: 'pre', // run before Vite's handlers so /dashboard gets index.html
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const url = req.url?.split('?')[0] ?? ''
        // Rewrite any non-file request (no dot in path segment) to index.html
        if (!url.includes('.') && !url.startsWith('/@') && !url.startsWith('/node_modules')) {
          req.url = '/index.html'
        }
        next()
      })
    },
  }
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss(), spaFallback()],
  appType: 'spa',
})
