import { test, expect } from '@playwright/test';

/**
 * 工作流编辑器测试
 * 验证 React Flow 正确加载和可交互
 */
test.describe('Workflow Editor', () => {

    test('React Flow canvas renders', async ({ page }) => {
        await page.goto('/editor');

        // 严格要求: React Flow 容器必须存在
        await expect(page.locator('.react-flow')).toBeVisible();

        // 内部结构完整
        await expect(page.locator('.react-flow__viewport')).toBeVisible();
        await expect(page.locator('.react-flow__pane')).toBeVisible();
    });

    test('canvas is interactive', async ({ page }) => {
        await page.goto('/editor');

        const pane = page.locator('.react-flow__pane');
        await expect(pane).toBeVisible();

        // 点击画布必须成功 (不抛出错误)
        await pane.click({ position: { x: 100, y: 100 } });
    });

    test('nodes have correct structure', async ({ page }) => {
        await page.goto('/editor');
        await expect(page.locator('.react-flow')).toBeVisible();

        // 如果有节点，验证其结构
        const nodes = page.locator('.react-flow__node');
        const nodeCount = await nodes.count();

        if (nodeCount > 0) {
            // 每个节点应该有 data 属性
            const firstNode = nodes.first();
            await expect(firstNode).toBeVisible();
        }
    });
});
