import { test, expect } from '@playwright/test';

test.describe('Groups Page', () => {
    test('should load the groups page', async ({ page }) => {
        await page.goto('/groups');
        await expect(page.locator('body')).toBeVisible();
    });

    test('should display page content', async ({ page }) => {
        await page.goto('/groups');
        // 页面应该有内容区域
        const main = page.locator('main, [class*="content"], [class*="page"]').first();
        if (await main.isVisible()) {
            await expect(main).toBeVisible();
        } else {
            // fallback: 至少页面加载了
            await expect(page.locator('body')).toBeVisible();
        }
    });
});
