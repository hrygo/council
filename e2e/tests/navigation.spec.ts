import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
    test('should load the home page', async ({ page }) => {
        await page.goto('/');

        // Check that the page loads successfully
        await expect(page).toHaveTitle(/Council|理事会/);
    });

    test('should navigate to all main sections', async ({ page }) => {
        await page.goto('/');

        // Navigate to Meeting
        await page.click('button:has-text("Meeting")');
        await expect(page.url()).toContain('/meeting');

        // Navigate to Builder
        await page.click('button:has-text("Builder"), button:has-text("工作流")');
        await expect(page.url()).toContain('/editor');

        // Navigate to Groups
        await page.click('button:has-text("Groups"), button:has-text("群组")');
        await expect(page.url()).toContain('/groups');

        // Navigate to Agents
        await page.click('button:has-text("Agents"), button:has-text("智能体")');
        await expect(page.url()).toContain('/agents');
    });

    test('should switch language', async ({ page }) => {
        await page.goto('/');

        // Find the language switcher
        const langSwitcher = page.locator('select[aria-label="Select language"]');

        // Switch to English
        await langSwitcher.selectOption('en');
        await expect(page.locator('button:has-text("Groups")')).toBeVisible();

        // Switch to Chinese
        await langSwitcher.selectOption('zh-CN');
        await expect(page.locator('button:has-text("群组")')).toBeVisible();
    });
});
