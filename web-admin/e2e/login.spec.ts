 import { test, expect } from '@playwright/test'

 test('login page renders', async ({ page }) => {
   await page.goto('/login')
   await expect(page.getByText('Asgard 管理控制台')).toBeVisible()
   await expect(page.getByRole('button', { name: '登录' })).toBeVisible()
 })
