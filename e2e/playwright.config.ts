import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright E2E 配置
 * 
 * 设计原则:
 * 1. 快速失败 - 短超时发现问题
 * 2. 严格断言 - 不容忍模糊状态
 * 3. 可观测 - 失败时有足够调试信息
 */
export default defineConfig({
    testDir: './tests',

    // 超时配置 - 平衡速度与稳定性
    timeout: 10 * 1000,          // 单个测试 10s (页面应在此时间内响应)
    expect: {
        timeout: 3000,            // 断言 3s (元素显示不应等太久)
    },

    // 执行配置
    fullyParallel: true,
    forbidOnly: !!process.env.CI,
    retries: process.env.CI ? 1 : 0,  // CI 重试一次
    workers: process.env.CI ? 2 : 4,

    // 报告配置
    reporter: [
        ['line'],
        ['html', { open: 'never' }],
    ],

    // 全局配置
    use: {
        baseURL: 'http://localhost:5173',

        // 调试信息
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure',
        video: 'retain-on-failure',

        // 性能优化
        launchOptions: {
            args: ['--no-sandbox', '--disable-gpu', '--disable-dev-shm-usage'],
        },

        // 严格模式
        actionTimeout: 5000,      // 操作超时 5s
        navigationTimeout: 8000, // 导航超时 8s
    },

    projects: [
        {
            name: 'chromium',
            use: { ...devices['Desktop Chrome'] },
        },
    ],
});
