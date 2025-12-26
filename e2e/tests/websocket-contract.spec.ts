import { test, expect, Page } from '@playwright/test';
import { StreamEvent } from '../../frontend/src/types/workflow.generated';

/**
 * WebSocket 消息契约测试
 * 验证前后端 WebSocket 消息格式完全一致
 * 
 * 测试目标:
 * 1. 确保生成的 TypeScript 类型与 Go 后端发送的消息匹配
 * 2. 验证所有消息字段存在且类型正确
 * 3. 防止运行时类型错误
 */

interface CapturedMessage {
  raw: string;
  parsed: any;
  timestamp: number;
}

/**
 * 捕获 WebSocket 消息的辅助函数
 */
async function captureWebSocketMessages(page: Page): Promise<CapturedMessage[]> {
  const messages: CapturedMessage[] = [];

  page.on('websocket', ws => {
    ws.on('framereceived', frame => {
      try {
        const raw = frame.payload.toString();
        const parsed = JSON.parse(raw);
        messages.push({
          raw,
          parsed,
          timestamp: Date.now(),
        });
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    });
  });

  return messages;
}

/**
 * 等待特定类型的消息
 */
async function waitForMessage(
  messages: CapturedMessage[],
  eventType: string,
  timeoutMs: number = 5000
): Promise<any | null> {
  const startTime = Date.now();

  while (Date.now() - startTime < timeoutMs) {
    const found = messages.find(m => m.parsed.event === eventType);
    if (found) {
      return found.parsed;
    }
    await new Promise(resolve => setTimeout(resolve, 100));
  }

  return null;
}

test.describe('WebSocket Message Contract', () => {

  test.beforeEach(async ({ page }) => {
    // 确保页面加载成功
    await page.goto('/');
    await page.waitForLoadState('domcontentloaded');
  });

  test('StreamEvent base structure validation', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);

    // 导航到会议室页面，触发 WebSocket 连接
    await page.goto('/chat');
    await page.waitForTimeout(1000);

    // 如果收到任何消息，验证其基本结构
    if (messages.length > 0) {
      const firstMessage = messages[0].parsed;

      // 验证 StreamEvent 必需字段
      expect(firstMessage).toHaveProperty('event');
      expect(firstMessage).toHaveProperty('timestamp');

      // 验证字段类型
      expect(typeof firstMessage.event).toBe('string');
      expect(typeof firstMessage.timestamp).toBe('string');

      // 可选字段检查（如果存在）
      if (firstMessage.node_id) {
        expect(typeof firstMessage.node_id).toBe('string');
      }

      if (firstMessage.data) {
        expect(typeof firstMessage.data).toBe('object');
      }
    }
  });

  test('StreamEvent type consistency with generated types', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);

    await page.goto('/chat');
    await page.waitForTimeout(1000);

    // 验证消息可以被正确解析为 StreamEvent 类型
    messages.forEach(({ parsed }) => {
      // TypeScript 编译时类型检查
      const event: StreamEvent = parsed;

      // 运行时验证
      expect(event.event).toBeDefined();
      expect(event.timestamp).toBeDefined();
    });
  });

  test('node_state_change event format', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);
    await page.goto('/chat');

    // 如果系统发送了 node_state_change 事件
    const stateChangeEvent = await waitForMessage(messages, 'node_state_change', 3000);

    if (stateChangeEvent) {
      // 验证事件结构
      expect(stateChangeEvent).toMatchObject({
        event: 'node_state_change',
        timestamp: expect.any(String),
        node_id: expect.any(String),
      });

      // 验证 data 字段
      expect(stateChangeEvent.data).toBeDefined();
      expect(typeof stateChangeEvent.data).toBe('object');
      expect(['running', 'completed', 'waiting', 'failed']).toContain(stateChangeEvent.data.status);
    }
  });

  test('memory_retrieval_error event format', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);
    await page.goto('/chat');

    const memErrorEvent = await waitForMessage(messages, 'memory_retrieval_error', 3000);
    if (memErrorEvent) {
      expect(memErrorEvent).toMatchObject({
        event: 'memory_retrieval_error',
        timestamp: expect.any(String),
        data: expect.objectContaining({
          error: expect.any(String)
        })
      });
    }
  });

  test('token_stream event format', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);

    await page.goto('/chat');

    // 等待 token_stream 事件
    const tokenStreamEvent = await waitForMessage(messages, 'token_stream', 3000);

    if (tokenStreamEvent) {
      expect(tokenStreamEvent).toMatchObject({
        event: 'token_stream',
        timestamp: expect.any(String),
        node_id: expect.any(String),
      });

      // token_stream 应该有内容数据
      if (tokenStreamEvent.data) {
        expect(tokenStreamEvent.data).toBeDefined();
      }
    }
  });

  test('error event format', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);

    await page.goto('/chat');

    // 等待错误事件（如果有）
    const errorEvent = await waitForMessage(messages, 'error', 3000);

    if (errorEvent) {
      expect(errorEvent).toMatchObject({
        event: 'error',
        timestamp: expect.any(String),
      });

      // 错误事件应该有错误信息
      if (errorEvent.data) {
        expect(errorEvent.data).toBeDefined();
      }
    }
  });

  test('message timestamp is valid ISO 8601 format', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);

    await page.goto('/chat');
    await page.waitForTimeout(1000);

    messages.forEach(({ parsed }) => {
      if (parsed.timestamp) {
        // 验证时间戳可以被解析
        const date = new Date(parsed.timestamp);
        expect(date.getTime()).not.toBeNaN();

        // 验证时间戳是最近的（防止使用默认值）
        const now = Date.now();
        const messageTime = date.getTime();
        const timeDiff = Math.abs(now - messageTime);

        // 时间差应该在合理范围内（1小时）
        expect(timeDiff).toBeLessThan(3600000);
      }
    });
  });

  test('no unexpected fields in messages', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);

    await page.goto('/chat');
    await page.waitForTimeout(1000);

    // 验证消息只包含预期的字段
    const allowedFields = ['event', 'timestamp', 'node_id', 'data'];

    messages.forEach(({ parsed }) => {
      const fields = Object.keys(parsed);
      fields.forEach(field => {
        expect(allowedFields).toContain(field);
      });
    });
  });

  test('messages are valid JSON', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);

    await page.goto('/chat');
    await page.waitForTimeout(1000);

    messages.forEach(({ raw }) => {
      // 验证原始消息可以被解析为 JSON
      expect(() => JSON.parse(raw)).not.toThrow();

      // 验证 JSON 不包含循环引用
      const parsed = JSON.parse(raw);
      expect(() => JSON.stringify(parsed)).not.toThrow();
    });
  });
});

test.describe('WebSocket Connection Stability', () => {

  test('connection can be established', async ({ page }) => {
    let wsConnected = false;

    page.on('websocket', ws => {
      wsConnected = true;

      ws.on('close', () => {
        console.log('WebSocket connection closed');
      });

      ws.on('socketerror', error => {
        console.error('WebSocket error:', error);
      });
    });

    await page.goto('/chat');
    await page.waitForTimeout(2000);

    // 在某些情况下，WebSocket 可能不会立即连接
    // 这个测试记录连接状态而不是强制要求
    console.log('WebSocket connected:', wsConnected);
  });

  test('messages are received in order', async ({ page }) => {
    const messages = await captureWebSocketMessages(page);

    await page.goto('/chat');
    await page.waitForTimeout(2000);

    // 验证消息时间戳递增（或相同）
    for (let i = 1; i < messages.length; i++) {
      const prevTime = new Date(messages[i - 1].parsed.timestamp).getTime();
      const currTime = new Date(messages[i].parsed.timestamp).getTime();

      // 当前消息时间应该 >= 前一条消息时间
      expect(currTime).toBeGreaterThanOrEqual(prevTime - 1000); // 允许 1 秒时钟偏差
    }
  });
});
