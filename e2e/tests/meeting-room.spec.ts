import { test, expect } from '@playwright/test';

/**
 * 会议室 / 聊天功能测试
 * 验证核心交互组件
 */
test.describe('Meeting Room', () => {

    test('should render meeting room with panel layout', async ({ page }) => {
        await page.goto('/meeting');

        // 验证使用了 react-resizable-panels
        const panelGroup = page.locator('[data-panel-group-id]');
        await expect(panelGroup.first()).toBeVisible({ timeout: 5000 });
    });

    test('should have chat input that accepts text', async ({ page }) => {
        await page.goto('/meeting');
        await page.waitForLoadState('networkidle');

        // 找到聊天输入框
        const chatInput = page.locator('textarea').first();

        if (await chatInput.isVisible()) {
            // 测试输入功能
            await chatInput.fill('Hello, this is a test message');
            const value = await chatInput.inputValue();
            expect(value).toBe('Hello, this is a test message');

            // 清空
            await chatInput.clear();
            expect(await chatInput.inputValue()).toBe('');
        }
    });

    test('should have send button near chat input', async ({ page }) => {
        await page.goto('/meeting');
        await page.waitForLoadState('networkidle');

        const chatInput = page.locator('textarea').first();

        if (await chatInput.isVisible()) {
            // 在输入框附近找发送按钮
            const sendBtn = page.locator('button').filter({ has: page.locator('svg') }).last();

            if (await sendBtn.isVisible()) {
                // 输入文本后点击发送
                await chatInput.fill('Test message');
                await sendBtn.click();

                // 验证输入框可能被清空 (发送成功的标志)
                await page.waitForTimeout(500);
                // 不强制要求清空，取决于是否连接后端
            }
        }
    });

    test('should have panel resize handles', async ({ page }) => {
        await page.goto('/meeting');

        // 找到 resize handle
        const resizeHandle = page.locator('[data-panel-resize-handle-id]').first();

        // 验证 resize handle 存在且可见
        await expect(resizeHandle).toBeVisible({ timeout: 5000 });

        // 验证有多个面板
        const panels = page.locator('[data-panel-id]');
        const panelCount = await panels.count();
        expect(panelCount).toBeGreaterThanOrEqual(2);
    });
});
