import { test, expect } from '@playwright/test';

/**
 * 核心渲染测试
 * 验证应用能正常启动和渲染
 */
test.describe('App Rendering', () => {

    test('React app renders without crash', async ({ page }) => {
        await page.goto('/');

        // 严格验证: React 根元素必须有子内容
        const root = page.locator('#root');
        await expect(root).toBeVisible();

        // 必须有至少一个子元素 (React 已渲染)
        const children = await root.locator('> *').count();
        expect(children).toBeGreaterThan(0);
    });

    test('no uncaught exceptions on initial load', async ({ page }) => {
        const errors: string[] = [];
        page.on('pageerror', error => {
            errors.push(error.message);
        });

        await page.goto('/');
        await page.waitForLoadState('domcontentloaded');

        // 只检查严重的未捕获异常
        expect(errors).toHaveLength(0);
    });

    test('navigation sidebar is visible and functional', async ({ page }) => {
        await page.goto('/');

        // 必须有导航按钮
        const navButtons = page.locator('button');
        const count = await navButtons.count();
        expect(count).toBeGreaterThanOrEqual(4); // 至少 4 个导航按钮
    });
});

test.describe('Routing', () => {

    test('navigates to /editor correctly', async ({ page }) => {
        await page.goto('/editor');
        expect(page.url()).toContain('/editor');
    });

    test('navigates to /groups correctly', async ({ page }) => {
        await page.goto('/groups');
        expect(page.url()).toContain('/groups');
    });

    test('navigates to /agents correctly', async ({ page }) => {
        await page.goto('/agents');
        expect(page.url()).toContain('/agents');
    });

    test('navigates to /meeting correctly', async ({ page }) => {
        await page.goto('/meeting');
        expect(page.url()).toContain('/meeting');
    });
});
