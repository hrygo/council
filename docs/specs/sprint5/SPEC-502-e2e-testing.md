# SPEC-502: 端到端测试 (E2E Testing)

> **优先级**: P1  
> **类型**: QA  
> **预估工时**: 12h

## 1. 概述

引入端到端测试框架，覆盖核心用户流程，确保系统在真实浏览器环境中的稳定性。

## 2. 目标

- 集成 Playwright 测试框架
- 覆盖 ≥5 个核心用户场景
- 测试可在 CI 中自动运行
- 生成可视化测试报告

## 3. 技术方案

### 3.1 框架选型

| 框架           | 优点                             | 缺点               | 选择 |
| :------------- | :------------------------------- | :----------------- | :--- |
| **Playwright** | 多浏览器、自动等待、Trace Viewer | 学习曲线           | ✅    |
| Cypress        | 调试友好、生态丰富               | 跨域限制、单浏览器 | -    |
| Puppeteer      | 轻量、Chrome 原生                | 仅 Chromium        | -    |

### 3.2 项目结构

```
e2e/
├── playwright.config.ts
├── fixtures/
│   └── test-data.json
├── pages/
│   ├── WorkflowBuilderPage.ts
│   ├── ChatPanelPage.ts
│   └── AgentsPage.ts
├── tests/
│   ├── workflow-builder.spec.ts
│   ├── agent-management.spec.ts
│   ├── group-management.spec.ts
│   ├── chat-interaction.spec.ts
│   └── human-review.spec.ts
└── utils/
    └── auth.ts
```

### 3.3 配置文件

```typescript
// playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  timeout: 30 * 1000,
  expect: { timeout: 5000 },
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html', { open: 'never' }],
    ['junit', { outputFile: 'test-results/junit.xml' }],
  ],
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'firefox', use: { ...devices['Desktop Firefox'] } },
    { name: 'webkit', use: { ...devices['Desktop Safari'] } },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
});
```

## 4. 测试场景

### 4.1 工作流构建 (workflow-builder.spec.ts)

```typescript
import { test, expect } from '@playwright/test';

test.describe('Workflow Builder', () => {
  test('should create a simple workflow with two agents', async ({ page }) => {
    await page.goto('/builder');
    
    // Drag Start node
    await page.locator('[data-testid="node-palette-start"]').dragTo(
      page.locator('[data-testid="canvas"]')
    );
    
    // Drag Agent node
    await page.locator('[data-testid="node-palette-agent"]').dragTo(
      page.locator('[data-testid="canvas"]'), { targetPosition: { x: 300, y: 200 } }
    );
    
    // Connect nodes
    await page.locator('[data-testid="node-start"] .handle-source').dragTo(
      page.locator('[data-testid="node-agent"] .handle-target')
    );
    
    // Save workflow
    await page.click('[data-testid="save-button"]');
    await expect(page.locator('[data-testid="toast-success"]')).toBeVisible();
  });

  test('should estimate cost before execution', async ({ page }) => {
    await page.goto('/builder/123');
    await page.click('[data-testid="estimate-button"]');
    
    await expect(page.locator('[data-testid="cost-estimator"]')).toBeVisible();
    await expect(page.locator('[data-testid="estimated-tokens"]')).toHaveText(/\d+ tokens/);
  });
});
```

### 4.2 Agent 管理 (agent-management.spec.ts)

```typescript
test.describe('Agent Management', () => {
  test('should create a new agent', async ({ page }) => {
    await page.goto('/agents');
    await page.click('[data-testid="create-agent"]');
    
    await page.fill('[data-testid="agent-name"]', 'Test Agent');
    await page.fill('[data-testid="agent-persona"]', 'You are a helpful assistant.');
    await page.selectOption('[data-testid="model-provider"]', 'openai');
    
    await page.click('[data-testid="save-agent"]');
    
    await expect(page.locator('text=Test Agent')).toBeVisible();
  });

  test('should delete an agent', async ({ page }) => {
    await page.goto('/agents');
    await page.click('[data-testid="agent-card-test"] [data-testid="delete-button"]');
    await page.click('[data-testid="confirm-delete"]');
    
    await expect(page.locator('text=Test Agent')).not.toBeVisible();
  });
});
```

### 4.3 人工审核流程 (human-review.spec.ts)

```typescript
test.describe('Human Review', () => {
  test('should pause and wait for human decision', async ({ page }) => {
    // Start workflow with HumanReview node
    await page.goto('/meeting/workflow-with-review');
    await page.click('[data-testid="run-button"]');
    
    // Wait for review modal
    await expect(page.locator('[data-testid="human-review-modal"]')).toBeVisible({ timeout: 30000 });
    
    // Approve
    await page.click('[data-testid="approve-button"]');
    
    // Verify workflow continues
    await expect(page.locator('[data-testid="execution-status"]')).toHaveText('Completed', { timeout: 30000 });
  });
});
```

### 4.4 聊天交互 (chat-interaction.spec.ts)

```typescript
test.describe('Chat Interaction', () => {
  test('should send message and receive response', async ({ page }) => {
    await page.goto('/meeting/active-session');
    
    await page.fill('[data-testid="chat-input"]', 'Hello, agent!');
    await page.click('[data-testid="send-button"]');
    
    // Wait for response
    await expect(page.locator('[data-testid="message-bubble"]').last()).toContainText(/./);
  });
});
```

### 4.5 群组管理 (group-management.spec.ts)

```typescript
test.describe('Group Management', () => {
  test('should create and edit a group', async ({ page }) => {
    await page.goto('/groups');
    
    // Create
    await page.click('[data-testid="create-group"]');
    await page.fill('[data-testid="group-name"]', 'E2E Test Group');
    await page.click('[data-testid="save-group"]');
    
    // Verify
    await expect(page.locator('text=E2E Test Group')).toBeVisible();
    
    // Edit
    await page.click('[data-testid="group-card-e2e"] [data-testid="edit-button"]');
    await page.fill('[data-testid="group-name"]', 'Updated Group');
    await page.click('[data-testid="save-group"]');
    
    await expect(page.locator('text=Updated Group')).toBeVisible();
  });
});
```

## 5. CI 集成

```yaml
# .github/workflows/e2e.yml
name: E2E Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      
      - name: Install dependencies
        run: cd frontend && npm ci
      
      - name: Install Playwright
        run: npx playwright install --with-deps
      
      - name: Start backend
        run: make start-db && make start-backend &
      
      - name: Run E2E tests
        run: cd e2e && npx playwright test
      
      - uses: actions/upload-artifact@v4
        if: failure()
        with:
          name: playwright-report
          path: e2e/playwright-report/
```

## 6. 验收标准

- [ ] Playwright 配置完成，可本地运行
- [ ] 5 个核心场景测试通过率 100%
- [ ] CI 中 E2E 测试自动运行
- [ ] 失败时生成截图和 Trace
- [ ] 测试报告发布到 PR Comment

## 7. Makefile 集成

```makefile
# E2E Testing
e2e: ## 🧪 Run E2E tests
	@cd e2e && npx playwright test

e2e-ui: ## 🧪 Run E2E tests with UI
	@cd e2e && npx playwright test --ui

e2e-report: ## 📊 Open E2E test report
	@cd e2e && npx playwright show-report
```
