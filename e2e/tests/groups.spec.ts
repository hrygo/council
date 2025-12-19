import { test, expect } from '@playwright/test';

/**
 * 群组管理页面测试
 * 验证页面结构和交互
 */
test.describe('Groups Page', () => {

    test('page renders with content', async ({ page }) => {
        await page.goto('/groups');

        // 必须有页面内容 (不能是空白页)
        const root = page.locator('#root');
        const children = await root.locator('> *').count();
        expect(children).toBeGreaterThan(0);
    });

    test('has expected UI elements', async ({ page }) => {
        await page.goto('/groups');
        await page.waitForLoadState('domcontentloaded');

        // 至少有按钮可以交互
        const buttons = page.locator('button');
        const buttonCount = await buttons.count();
        expect(buttonCount).toBeGreaterThanOrEqual(1);
    });
});
