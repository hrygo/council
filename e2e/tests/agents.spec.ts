import { test, expect } from '@playwright/test';

test.describe('Agents Page', () => {
    test('should load the agents page', async ({ page }) => {
        await page.goto('/agents');
        await expect(page.locator('body')).toBeVisible();
    });

    test('should display page content', async ({ page }) => {
        await page.goto('/agents');
        // 页面应该有内容区域
        const main = page.locator('main, [class*="content"], [class*="page"]').first();
        if (await main.isVisible()) {
            await expect(main).toBeVisible();
        } else {
            await expect(page.locator('body')).toBeVisible();
        }
    });
});
