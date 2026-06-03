 import { defineConfig } from 'vitest/config'
 import vue from '@vitejs/plugin-vue'
 import { fileURLToPath, URL } from 'node:url'

 export default defineConfig({
   plugins: [vue()],
   resolve: {
     alias: {
       '@': fileURLToPath(new URL('./src', import.meta.url)),
     },
   },
   test: {
     globals: true,
     environment: 'jsdom',
     include: ['src/**/*.{test,spec}.ts'],
     coverage: {
       provider: 'v8',
       reporter: ['text', 'html'],
     },
   },
 })
