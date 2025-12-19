import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
    test('should load the home page', async ({ page }) => {
        await page.goto('/');
        // 等待页面加载完成
        await expect(page.locator('body')).toBeVisible();
    });

    test('should have navigation sidebar', async ({ page }) => {
        await page.goto('/');
        // 检查侧边栏存在
        const sidebar = page.locator('div').filter({ hasText: /Meeting|Builder|Groups|Agents|群组|智能体/ }).first();
        await expect(sidebar).toBeVisible();
    });

    test('should navigate to editor page', async ({ page }) => {
        await page.goto('/editor');
        // 检查 React Flow canvas 存在
        await expect(page.locator('.react-flow')).toBeVisible({ timeout: 10000 });
    });

    test('should navigate to groups page', async ({ page }) => {
        await page.goto('/groups');
        await expect(page.locator('body')).toBeVisible();
    });

    test('should navigate to agents page', async ({ page }) => {
        await page.goto('/agents');
        await expect(page.locator('body')).toBeVisible();
    });
});
