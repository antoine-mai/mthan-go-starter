import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

declare const process: { env: { [key: string]: string | undefined } }

const base = process.env.VITE_BASE || './'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: base,
  server: {
    port: 3001,
  },
  build: {
    outDir: 'build',
    emptyOutDir: true,
  }
})
