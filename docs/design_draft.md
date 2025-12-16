# The Council - UI/UX 设计蓝图 (Design Blueprint)

> **版本**: v1.1 | **状态**: 修订版 | **基准**: PRD v1.2.0 & Debate Verdict 20251216_130019

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

### 1.2 架构层防御机制 (Architectural Defense Mechanisms) 🆕
*   **逻辑熔断 (Logic Circuit Breaker)**:
    *   **触发条件**: 当系统检测到连续两轮对话未产出有效差异、陷入无限循环或 Token 消耗速率异常激增时。
    *   **响应**: 系统自动挂起当前进程，切断 LLM 调用。
    *   **UI表现**: 界面四周出现红色脉冲呼吸灯 (Crimson Pulse)，并弹出模态窗口 "Logic Circuit Breaker Triggered"。用户必须手动点击 "Override" 才能继续，防止幻觉级联传播。
*   **防幻觉传播 (Anti-Hallucination Propagation)**:
    *   **机制**: 在 Agent 传递信息间隙增加 "Fact Verification Layer"。
    *   **UI表现**: 对于高风险陈述 (如具体数据、引用)，在气泡旁标记 "Verify Pending" (黄色警示)，提醒下一环节的 Agent 或人类注意核实。

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
        *   **交互**: 选中节点时，在节点旁浮出或在屏幕右侧悬浮显示（非定宽侧栏，避免视线长距离跳跃）。
        *   **内容**: 模型参数 (Temperature, Top_P)、Prompt 模板、输入输出定义。
    *   **底部状态栏**: 显示流程复杂度预估、预计成本范围。

### 2.3 作业区 B：运行态 (Run Mode - The Meeting)
*   **核心逻辑**: **"轻量化左栏 + 宽幅阅读区"**。避免在运行时渲染重型 DAG，释放 GPU 资源。

#### A. 左栏：时间轴 (Timeline Sidebar)
*   **宽度**: 固定 240px 或 收缩为 Icons (64px)。
*   **内容**: **线性化的进度视图 (Linear Step View)**。
    *   不显示复杂的 DAG 连线，仅将当前执行路径扁平化为垂直的时间轴节点。
    *   **状态指示**: 
        *   🟢 完成 (Done)
        *   🔵 进行中 (Spinner + 耗时计时器)
        *   ⚪ 待定 (Pending)
    *   **循环逻辑处理 (Loop Handling)** 🆕:
        *   对于 `Loop` 节点，不无限延伸线性列表。
        *   采用 **"嵌套组 (Nested Group)"** 设计：显示为 "Round 2/3 (Debating...)" 的可折叠容器。
        *   点击容器展开查看历史轮次的详细步骤。
*   **性能优势**: 纯 DOM/SVG 渲染，无 WebGL 开销。

#### B. 中栏：众议院 (Chat Stream) - **视觉绝对焦点**
*   **宽度**: 弹性占据 50% - 70% 屏幕空间。
*   **布局**: 垂直流 (Vertical Stream)。
*   **关键组件**:
    *   **专业级成本仪表盘 (Cost Dashboard)** 🆕:
        *   **位置**: 顶部导航栏右侧或吸附在输入框上方。
        *   **样式**: 单色、静止的低对比度文本。
        *   **策略**: **“里程碑式汇报 (Milestone Reporting)”**。
            *   **禁止**: 秒级跳动的实时金额。
            *   **逻辑**: 仅在流程关键节点（如阶段完成、暂停时）更新显示，或显示静态的“本次会议总成本”。
    *   **并行发言 (Parallel Group)** 🆕:
        *   **策略**: **同时可见性原则 (Simultaneous Visibility)**。
        *   **禁止**: 使用 Tab 切换隐藏内容 (无论屏幕宽度如何)。
        *   **布局**:
            *   **宽屏**: 并排卡片 (Side-by-Side Cards)。
            *   **窄屏**: **横向滚动容器 (Horizontal Scroll Container)** 或 **多栏挤压布局 (Squeezed Columns)**。必须保证用户能同时看到多个 Agent 的存在和状态，便于横向对比。
    *   **Human Review (人类裁决) - Diff View**:
        *   **优先级**: **P0 (Critical)**。
        *   **样式**: **语义级差异对比 (Semantic Diff)** (类似 Google Docs / Notion 建议模式)。
        *   **禁止**: 纯行的代码 Diff (Monaco Default)。
        *   **功能**: 基于单词/短语 (Word-level) 的高亮修订，而非整行标红。
        *   **动作**: [ 驳回 ] [ 签署通过 ]。

#### C. 右栏：卷宗 (Context Drawer)
*   **状态**: 默认折叠或以较窄宽度 (25%) 存在。
*   **内容**: 文档阅读器、记忆检索结果。
*   **联动**: 点击对话中的引用链接，右栏自动展开并定位。

---

## 3. 关键交互流程 (Interaction Flows)

### 3.1 发起提案 (Wizard)
*   保持原有三步设计：意图 -> 流程选择 -> 成本预估。
*   **新增**: 在确认流程后，用户可选择 "Run Immediately (立即开会)" 或 "Open in Builder (进入构建态微调)"。

### 3.2 会议中干预 (Intervention)
*   用户输入消息时，系统自动在 Timeline 插入 "User Intervention" 节点。
*   **暂停/继续**: 底部提供明显的全局控制按钮 [ ⏸ 暂停 ] [ ▶ 继续 ] [ ⏹ 强行终止 ] (用于熔断高消耗任务)。

### 3.3 记忆净化协议 (Memory Purification Protocol)
*   **风险阻断**: 针对 "Memory Pollution" (记忆污染) 风险，采取 **"默认丢弃 (Opt-in)"** 策略。
*   **交互流程**: 
    *   会议结束后，**强制**弹出 "本次会议知识沉淀 (Session Knowledge Impact)" 面板。
    *   系统列出提取的知识点，但 **默认全部不勾选** (尤其是从 Conflict 状态推导出的结论)。
    *   用户必须 **手动勾选** 确认高信度的条目，才能将其写入长期向量数据库。这一步是防止错误信息污染核心智库的最后一道防线。

---

## 4. 响应式策略 (Responsive Strategy)

*   **Ultra Wide (>1600px)**: 
    *   运行态：左(Timeline) + 中(Chat) + 右(Docs) 并存。
    *   并行消息：允许横向分栏显示 (Side-by-Side) 2-3 个 Agent。
*   **Desktop/Laptop (1280px - 1600px)**:
    *   左栏(Timeline) 保持可见。
    *   右栏(Docs) 默认折叠，需手动展开。
    *   并行消息：**横向滚动卡片**，保证内容不被 Tabs 隐藏。
*   **Small Laptop/Tablet (<1280px)**:
    *   左栏收起为图标条。
    *   专注于中间对话流。
    *   并行消息：**纵向交错 (Vertical Interleaved)** 或 **横向滚动**。

---

## 5. UI 组件规范 (Component Specs)

*   **Diff Editor**: 使用支持语义 Diff 的库 (如 `diff-match-patch` 封装组件)，而非简单的 Monaco Diff Editor。
*   **Cost Dashboard**: 数字使用 `JetBrains Mono` 或 `Roboto Mono`，字号 `text-xs`，颜色 `text-slate-500` (低调安静)。
*   **Loaders**: 使用精致的 SVG 动画代替系统默认 Spinner，体现 "Premium" 质感。

---

## 6. 附录：PRD 文档

见 ./PRD.md