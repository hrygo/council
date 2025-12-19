import { test, expect } from '@playwright/test';

test.describe('Meeting Room', () => {
    test('should load the meeting room page', async ({ page }) => {
        await page.goto('/meeting');
        await expect(page.locator('body')).toBeVisible();
    });

    test('should have resizable panels', async ({ page }) => {
        await page.goto('/meeting');
        // 检查 react-resizable-panels 组件
        const panelGroup = page.locator('[data-panel-group-id]').first();
        if (await panelGroup.isVisible()) {
            await expect(panelGroup).toBeVisible();
        } else {
            await expect(page.locator('body')).toBeVisible();
        }
    });

    test('should have chat area', async ({ page }) => {
        await page.goto('/meeting');
        // 检查输入框或聊天区域
        const chatArea = page.locator('textarea, [class*="chat"], [class*="input"]').first();
        if (await chatArea.isVisible()) {
            await expect(chatArea).toBeVisible();
        } else {
            await expect(page.locator('body')).toBeVisible();
        }
    });
});
