import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright E2E Test Configuration
 * 
 * 优化重点:
 * 1. 移除 webServer 自动启动 - 要求先手动运行 `make start-frontend`
 * 2. 缩短超时时间提高反馈速度
 * 3. 使用 line reporter 实时显示进度
 */
export default defineConfig({
    testDir: './tests',

    // 超时配置 - 优化执行效率
    timeout: 15 * 1000,        // 单个测试 15s 超时
    expect: {
        timeout: 3000,          // 断言 3s 超时
    },

    // 执行配置
    fullyParallel: true,        // 并行执行提高效率
    forbidOnly: !!process.env.CI,
    retries: 0,                 // 本地不重试，快速失败
    workers: 4,                 // 4 个 worker 并行

    // 报告配置 - 实时进度
    reporter: [
        ['line'],               // 简洁的单行进度输出
        ['html', { open: 'never' }],
    ],

    // 浏览器配置
    use: {
        baseURL: 'http://localhost:5173',
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure',
        // 快速启动
        launchOptions: {
            args: ['--no-sandbox', '--disable-gpu'],
        },
    },

    // 仅 Chromium
    projects: [
        {
            name: 'chromium',
            use: { ...devices['Desktop Chrome'] },
        },
    ],

    // 不自动启动 webServer - 要求先运行 `make start-frontend`
    // webServer: { ... } 已移除
});
