# The Council - UI/UX 设计蓝图 (Design Blueprint)

> **版本**: v1.5 | **状态**: 正式版 (Production Ready)

---

## 1. 核心布局策略 (Global Layout Strategy)

采用 **"模式分离架构 (Mode Separation Architecture)"**，将复杂的编排与流式的执行在空间上解耦：

1.  **管理层 (Dashboard)**：以"卡片流"为主，强调群组概览与资产管理。
2.  **作业层 (Workspace)**：分为两种互斥的视图模式：
    *   **构建态 (Build Mode)**：全屏 IDE 体验，专注于 DAG 流程编排与参数微调。
    *   **运行态 (Run Mode)**：沉浸式会议体验，专注于阅读、决策与文档消费。

### 1.1 设计语言系统 (Design Language System)
*   **色调 (Palette)**：
    *   **材质**: 纯色背景 + 细微边框 (Subtle Borders) + 极其克制的半透明 (Micro-blur)。优先保证信息密度与性能。
    *   **阅读体验优化**: 正文文本颜色采用 `#E0E0E0` (非纯白)，降低高对比度导致的视疲劳。
    *   **背景**: Deep Space Blue (深空蓝/黑) - 营造专注、神秘的"理事会"氛围。
    *   **主色**: Electric Blue (电光蓝) - 强调 AI 智能与科技感。
    *   **辅助色**:
        *   Gold (暗金) - 象征核心决策/关键结论。
        *   Crimson (绯红) - 错误、驳回、高风险操作。
        *   Emerald (祖母绿) - 通过、验证成功、低成本。
*   **字体**: Inter (西文) + Noto Sans SC (中文) - 只有无衬线字体，极简现代。
*   **数据可视化**:
    *   **Dashboard Style**: 采用 Grafana/Bloomberg 风格的数据展示，使用等宽数字 (Monospace Numbers)，拒绝电商风的“价格标签”。

### 1.2 架构层防御机制 (Architectural Defense Mechanisms)
系统内置两层防御机制，以保障高风险决策的安全性和可靠性。

*   **逻辑熔断 (Logic Circuit Breaker)**:
    *   **机制**: 实时监控对话熵值。当检测到连续两轮对话无有效信息增量、陷入死循环或 Token 消耗速率异常激增时触发。
    *   **响应 (Hard Stop)**: 立即终止所有 LLM 进程，将系统状态锁定为 `SUSPENDED (LOCKED)`。
    *   **恢复流程**: 界面转为灰阶模式 (Grayscale)，仅保留红色警示。用户必须填写 "Risk Justification" (风险陈述) 或通过三项 "Safety Checks" (安全自查) 才能解锁系统，强制阻断无意识的操作惯性。

*   **防幻觉传播 (Anti-Hallucination Propagation)**:
    *   **机制**: 在 Agent 信息传递链路上部署 "Fact Verification Layer" (事实校验层)。
    *   **UI表现**: 涉及具体数据或外部引用的陈述，若未经交叉验证，气泡旁会自动标记 "Verify Pending" (黄色警示)，提示下一环节关注。

---

## 2. 界面结构详述 (Interface Detail)

### 2.1 首页：理事会大厅 (Dashboard)
*   **布局**: 顶栏 + 居中限宽内容区。
*   **群组网格**: 保持 GitHub Repo Card 风格的高密度卡片，展示 Token 消耗总额、活跃状态。

### 2.2 作业区 A：构建态 (Build Mode - The IDE)
*   **入口**: 提案向导结束或点击“编辑流程”时进入。
*   **布局**: **全屏画布 (Infinite Canvas)**。
*   **交互逻辑**:
    *   **画布**: 占据 100% 屏幕空间。
    *   **节点属性面板 (Floating Property Panel)**: 
        *   选中节点时在旁浮出或在屏幕右侧悬浮，避免视线长距离跳跃。
        *   配置模型参数 (Temperature, Top_P)、Prompt 模板、输入输出定义。
    *   **底部状态栏**: 显示流程复杂度预估、预计成本范围。

### 2.3 作业区 B：运行态 (Run Mode - The Meeting)
采用 **"轻量化左栏 + 宽幅阅读区"** 布局，避免运行时渲染重型 DAG，专注即时通讯体验。

#### A. 左栏：时间轴 (Timeline Sidebar)
*   **视图**: **线性进度视图 (Linear Step View)**。将执行路径扁平化为垂直时间轴。
*   **状态指示**: 🟢 完成 (Done)、🔵 进行中 (Spinner + 计时器)、⚪ 待定 (Pending)。
*   **循环组 (Nested Groups)**: 
    *   对于 `Loop` 节点，采用折叠容器设计 (例如 "Round 2/3")，保持时间轴整洁。
    *   支持点击展开查看历史轮次详情。
*   **技术栈**: 纯 DOM/SVG 渲染，无 WebGL 开销。

#### B. 中栏：众议院 (Chat Stream) - **视觉绝对焦点**
*   **布局**: 垂直流 (Vertical Stream)，占据 50% - 70% 宽幅。
*   **关键组件**:
    *   **成本仪表盘 (Cost Dashboard)**:
        *   采用 **“里程碑式汇报”** 策略：仅在阶段结束或暂停时更新，显示静态数值。
        *   避免秒级跳动的金额造成焦虑干扰。
    *   **并行发言组 (Parallel Group)**:
        *   遵循 **同时可见性原则 (Simultaneous Visibility)**。
        *   **布局**: 宽屏下并排展示卡片；窄屏下使用横向滚动容器或挤压布局。确保用户能同时对比多个 Agent 的状态，严禁使用 Tabs 隐藏内容。
    *   **Human Review (人类裁决) - Diff View**:
        *   **样式**: **语义级差异对比 (Semantic Diff)** (类似 Notion 建议模式)。
        *   基于单词/短语 (Word-level) 高亮修订，而非代码级的整行标红。
        *   操作: [ 驳回 ] [ 签署通过 ]。

#### C. 右栏：卷宗 (Context Drawer)
*   **状态**: 默认折叠或以窄宽 (25%) 存在。
*   **内容**: 文档阅读器、记忆检索结果。
*   **联动**: 点击对话引用链接时自动展开并定位。

---

## 3. 关键交互流程 (Interaction Flows)

### 3.1 发起提案 (Wizard)
*   三步向导：意图 -> 流程选择 -> 成本预估。
*   完成指引后，提供 "Run Immediately" (立即开会) 和 "Open in Builder" (进入构建态微调) 两个分支入口。

### 3.2 会议中干预 (Intervention)
*   用户输入消息时，系统自动在 Timeline 插入 "User Intervention" 节点。
*   **全局控制**: 底部常驻 [ ⏸ 暂停 ] [ ▶ 继续 ] [ ⏹ 强行终止 ] 按钮。

### 3.3 记忆净化协议 (Memory Purification Protocol)
为防止长期运行导致向量库污染，系统实施严格的分层记忆管理策略。

1.  **第一道防线：Cortex Quarantine (记忆隔离区)**
    *   会议产生的所有原始知识点、结论默认进入此临时缓冲区。
    *   此区域物理隔离于核心向量库 (Long-Term Vector DB)，确保核心记忆库的纯净。

2.  **第二道防线：Working Memory Buffer (工作记忆热层)**
    *   **功能**: 为当前会话及短期交互提供即时上下文支持（解决“归档前真空期”问题）。
    *   **Ingress Filter (入口过滤)**: 仅允许通过 **Self-Consistency Check (自洽性检查)** 的内容进入热层。
    *   **生命周期**: 
        *   **TTL**: **24小时** 自动过期。
        *   **Scope**: 严格限制在当前 Project ID 内，禁止跨项目访问。
    *   **UI标识**: 引用此类数据时，标记 "⚡️ Ephemeral Context" (临时上下文) 图标。

3.  **第三道防线：Knowledge Promotion (知识晋升)**
    *   **Smart Digest (智能简报)**: 系统按周生成 "Knowledge Promotion Digest"。
    *   **机制**: 自动聚类 Quarantine 中的碎片信息，生成 5-10 条核心洞察。
    *   **归档**: 用户基于简报进行 "One-click Promote"，将经过验证的高价值知识写入 Long-Term DB。

---

## 4. 响应式策略 (Responsive Strategy)

*   **Ultra Wide (>1600px)**: 三栏并存 (Timeline + Chat + Docs)。并行消息支持横向分栏。
*   **Desktop (1280px - 1600px)**: 左栏可见，右栏默认折叠。并行消息使用横向滚动卡片。
*   **Tablet/Small (<1280px)**: 左栏收缩为图标条，专注于对话流。并行消息采用纵向交错或横向滚动。

---

## 5. UI 组件规范 (Component Specs)

*   **Diff Editor**: 使用支持语义 Diff 的库 (如 `diff-match-patch`)。
*   **Typography**: 数据面板使用 `JetBrains Mono`，字号 `text-xs`，颜色 `text-slate-500`。
*   **Motion**: 使用精致的 SVG 动画 (Lottie/Rive) 代替原生 Spinner。

---

## 6. 附录：PRD 文档

见 ./PRD.md