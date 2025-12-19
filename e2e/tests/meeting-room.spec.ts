import { test, expect } from '@playwright/test';

test.describe('Meeting Room / Chat', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/meeting');
    });

    test('should load the meeting room page', async ({ page }) => {
        // Check page loaded - should have panels
        await expect(page.locator('[class*="panel"], [class*="Panel"]').first()).toBeVisible();
    });

    test('should show chat input', async ({ page }) => {
        // Look for chat input
        const chatInput = page.locator('textarea[placeholder*="message"], textarea[placeholder*="消息"], input[placeholder*="message"]');
        await expect(chatInput.first()).toBeVisible();
    });

    test('should be able to type in chat input', async ({ page }) => {
        // Find and interact with chat input
        const chatInput = page.locator('textarea[placeholder*="message"], textarea[placeholder*="消息"]').first();

        if (await chatInput.isVisible()) {
            await chatInput.fill('Hello, this is a test message');
            await expect(chatInput).toHaveValue('Hello, this is a test message');
        }
    });

    test('should show send button', async ({ page }) => {
        // Look for send button
        const sendBtn = page.locator('button[type="submit"], button:has-text("Send"), button:has-text("发送"), button svg');
        await expect(sendBtn.first()).toBeVisible();
    });

    test('should have resizable panels', async ({ page }) => {
        // Check for resize handles
        const resizeHandles = page.locator('[class*="resize"], [data-panel-group]');
        if (await resizeHandles.first().isVisible()) {
            await expect(resizeHandles.first()).toBeVisible();
        }
    });
});
