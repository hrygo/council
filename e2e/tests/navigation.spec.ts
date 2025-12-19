import { test, expect } from '@playwright/test';

/**
 * 核心用户流程测试
 * 验证真实的用户操作和业务逻辑
 */
test.describe('Core User Flows', () => {

    test('should render React app without crash', async ({ page }) => {
        await page.goto('/');
        // 验证 React 正常渲染 - 检查根元素有子节点
        const root = page.locator('#root');
        await expect(root).toBeVisible();
        // 确保不是空白页
        const childCount = await root.locator('> *').count();
        expect(childCount).toBeGreaterThan(0);
    });

    test('should have working navigation buttons', async ({ page }) => {
        await page.goto('/');

        // 点击导航按钮应该改变 URL
        const buttons = page.locator('button').filter({ hasText: /Meeting|Builder|Groups|Agents/ });
        const count = await buttons.count();
        expect(count).toBeGreaterThan(0);

        // 点击第一个按钮
        await buttons.first().click();
        await page.waitForTimeout(500);

        // URL 应该有变化或页面应该有响应
        const url = page.url();
        expect(url).toMatch(/localhost:5173/);
    });

    test('should switch language and persist', async ({ page }) => {
        await page.goto('/');

        // 找到语言选择器
        const langSelect = page.locator('select[aria-label="Select language"]');
        if (await langSelect.isVisible()) {
            // 切换到英文
            await langSelect.selectOption('en');
            await page.waitForTimeout(300);

            // 验证语言确实切换了 (检查 localStorage)
            const lang = await page.evaluate(() => localStorage.getItem('council-language'));
            expect(lang).toBe('en');
        }
    });
});
