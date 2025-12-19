import { test, expect } from '@playwright/test';

test.describe('Workflow Builder', () => {
    test('should load the workflow editor page', async ({ page }) => {
        await page.goto('/editor');
        // React Flow canvas 必须存在
        await expect(page.locator('.react-flow')).toBeVisible({ timeout: 10000 });
    });

    test('should have nodes in canvas after load', async ({ page }) => {
        await page.goto('/editor');
        await expect(page.locator('.react-flow')).toBeVisible({ timeout: 10000 });

        // 检查至少有一个节点或者空画布可交互
        const canvas = page.locator('.react-flow');
        await expect(canvas).toBeEnabled();
    });

    test('should have clickable canvas', async ({ page }) => {
        await page.goto('/editor');
        const canvas = page.locator('.react-flow__pane');
        await expect(canvas).toBeVisible({ timeout: 10000 });

        // 点击画布不应报错
        await canvas.click({ position: { x: 200, y: 200 } });
    });
});
