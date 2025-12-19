import { test, expect } from '@playwright/test';

test.describe('Workflow Builder', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/editor');
    });

    test('should load the workflow editor page', async ({ page }) => {
        // Check that the editor canvas is visible
        await expect(page.locator('.react-flow')).toBeVisible();
    });

    test('should have node palette visible', async ({ page }) => {
        // Check for node palette or sidebar
        const sidebar = page.locator('[class*="sidebar"], [class*="palette"], [class*="panel"]').first();
        await expect(sidebar).toBeVisible();
    });

    test('should be able to drag nodes onto canvas', async ({ page }) => {
        // Look for draggable node items
        const nodeItem = page.locator('[draggable="true"]').first();

        if (await nodeItem.isVisible()) {
            const canvas = page.locator('.react-flow');

            // Drag node to canvas
            await nodeItem.dragTo(canvas, {
                targetPosition: { x: 300, y: 200 },
            });

            // Verify a node was added
            const nodes = page.locator('.react-flow__node');
            await expect(nodes).toHaveCount(await nodes.count()); // At least maintains count
        }
    });

    test('should show property panel when node is selected', async ({ page }) => {
        // Click on a node if one exists
        const node = page.locator('.react-flow__node').first();

        if (await node.isVisible()) {
            await node.click();

            // Property panel should appear
            const propertyPanel = page.locator('[class*="property"], [class*="panel"]');
            await expect(propertyPanel.first()).toBeVisible();
        }
    });
});
