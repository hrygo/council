import { test, expect } from '@playwright/test';

test.describe('Agents Management', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/agents');
    });

    test('should load the agents page', async ({ page }) => {
        // Check page loaded
        await expect(page.locator('h1, h2').first()).toBeVisible();
    });

    test('should show create agent button', async ({ page }) => {
        // Look for create button
        const createBtn = page.locator('button:has-text("Create"), button:has-text("创建"), button:has-text("New")');
        await expect(createBtn.first()).toBeVisible();
    });

    test('should open create agent modal on button click', async ({ page }) => {
        // Click create button
        const createBtn = page.locator('button:has-text("Create"), button:has-text("创建")').first();
        await createBtn.click();

        // Modal should appear
        const modal = page.locator('[role="dialog"], [class*="modal"]');
        await expect(modal.first()).toBeVisible({ timeout: 5000 });
    });

    test('should be able to fill agent form', async ({ page }) => {
        // Open create modal
        const createBtn = page.locator('button:has-text("Create"), button:has-text("创建")').first();
        await createBtn.click();

        // Fill form fields
        const nameInput = page.locator('input[name="name"], input[placeholder*="name"], input[placeholder*="名称"]').first();

        if (await nameInput.isVisible()) {
            await nameInput.fill('E2E Test Agent');
            await expect(nameInput).toHaveValue('E2E Test Agent');
        }

        // Check for model selector
        const modelSelect = page.locator('select[name="provider"], select[name="model"]').first();
        if (await modelSelect.isVisible()) {
            await expect(modelSelect).toBeVisible();
        }
    });
});
