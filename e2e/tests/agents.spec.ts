import { test, expect } from '@playwright/test';

/**
 * 智能体管理页面测试
 * 验证页面结构和交互
 */
test.describe('Agents Page', () => {

    test('page renders with content', async ({ page }) => {
        await page.goto('/agents');

        // 必须有页面内容
        const root = page.locator('#root');
        const children = await root.locator('> *').count();
        expect(children).toBeGreaterThan(0);
    });

    test('has expected UI elements', async ({ page }) => {
        await page.goto('/agents');
        await page.waitForLoadState('domcontentloaded');

        // 至少有按钮可以交互
        const buttons = page.locator('button');
        const buttonCount = await buttons.count();
        expect(buttonCount).toBeGreaterThanOrEqual(1);
    });
});
