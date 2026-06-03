 import { defineConfig, devices } from '@playwright/test'

 export default defineConfig({
   testDir: './e2e',
   timeout: 30_000,
   expect: { timeout: 5000 },
   fullyParallel: true,
   forbidOnly: !!process.env.CI,
   retries: process.env.CI ? 2 : 0,
   reporter: [['list'], ['html', { open: 'never' }]],
   use: {
     baseURL: process.env.E2E_BASE_URL || 'http://localhost:5173',
     trace: 'on-first-retry',
     screenshot: 'only-on-failure',
   },
   webServer: process.env.CI
     ? undefined
     : {
         command: 'npm run dev',
         url: 'http://localhost:5173',
         reuseExistingServer: true,
         timeout: 60_000,
       },
   projects: [
     { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
   ],
 })
