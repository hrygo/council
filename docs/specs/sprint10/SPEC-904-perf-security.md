# SPEC-1004: 性能硬化与安全防护 (Perf & Security)

## 1. 需求背景
- **性能**: 当会话消息超过 50 条时，前端渲染性能开始下降，DOM 节点过多导致浏览器响应迟缓。
- **安全**: 系统目前缺乏基础的 CSRF 和 XSS 防护，尤其是 AI 输出内容直接渲染可能存在注入风险。

## 2. 目标定义
- 实现 ChatPanel 的虚拟滚动渲染。
- 集成基础的安全防御层。

## 3. 设计方案

### 3.1 虚拟滚动 (`frontend/src/components/chat/ChatPanel.tsx`)
1.  **引入库**: 使用 `react-window` 或 `tanstack-virtual`。
2.  **动态高度**: 考虑到 AI 消息长度不一，需支持动态高度计算（Item Mapping）。
3.  **缓存机制**: 缓存已渲染项的尺寸。

### 3.2 安全增强
1.  **XSS 防御**: 
    - 引入 `dompurify`。
    - 对所有 `Markdown` 渲染后的 HTML 进行 Sanitization。
2.  **CSRF 防御 (MVP 级别)**:
    - 后端 Middleware 增加 Referer 校验。
    - 前端配置 Axios 携带 `withCredentials: true`。

## 4.验收标准
- [ ] 构造一个包含 500 条消息的模拟会话，滚动流畅，F12 显示 DOM 节点数稳定在 20 个以内。
- [ ] 尝试输出恶意 `<script>` 脚本，系统应能正确清理且不执行。
- [ ] 所有的 `innerHTML` 操作都经过了清理函数处理。
