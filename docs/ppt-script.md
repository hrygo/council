### 第一部分：视觉规范确认 (Design Protocol Initialization)

**视觉与风格规范 (Global Style & Visual DNA)**

**必须严格执行以下设计语言：**
* **核心风格 (Core Style)**: **Future Tech Glassmorphism (未来科技磨砂玻璃风)**
    * **视觉隐喻**: 悬浮的 2.5D 磨砂玻璃卡片 (Frosted Glass Cards) + 极简扁平矢量图标 (Clean Flat Vector Icons) + 柔和景深 (Soft Depth of Field)。
    * **背景环境**: 极简科技白 (`#ffffff` 至 `#f8f9fa`)，极低透明度 (3-5%) 的神经网络连线或网格纹理。
    * **字体与排版**: 标题位置固定左上角，采用 "Helvetica/Roboto" 风格，强调秩序感与留白。
* **品牌配色 (Color Palette)**:
    * **Deep Brand Blue (`#0073e5`)**: 主文字、边框、玻璃轮廓。
    * **Interactive Cyan (`#00bfff`)**: 玻璃渐变、发光点缀、AI 神经链路。
    * **Vivid Orange (`#ff531a`)**: 强调高光、Call-to-Action、关键数据指标。
* **通用绘图后缀 (Unified Suffix - Append to ALL Prompts)**:
    * > "2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background."

---

### 第二部分：分页脚本内容 (Slide Scripting)

#### Slide 01 :: LAUNCH // The Council (理事会)

`![[ppt-pages/P1.jpeg]]`

* **页面内容 (Slide Content)**:
    * **The Council (理事会)**
    * 可视化的多智能体协作决策系统 (Visualized Multi-Agent Collaborative System)
    * **v1.2.0 MVP_Simplified** 
    * 核心定位：**思维外脑容器 & 个人私有智库**

* **演讲稿 (Script)**:
    大家好。今天我们要发布的不是一个简单的聊天机器人，而是一个架构——一个属于你个人的“思维外脑容器”。我们称之为 **The Council (理事会)**。在过去，决策是孤独的；而在 The Council 中，我们将把单一的决策者转变为“理事会主席”。这是一个基于可视化的多智能体协作系统，专为极客与高管设计。让我们看看它是如何重新定义“思考”的。

* **Image Generation Prompt**:
    * `Subject`: A central, glowing 2.5D frosted glass sphere representing a "digital brain" or "core", surrounded by three orbiting smaller glass chips representing AI agents. The central sphere emits soft Cyan light (`#00bfff`).
    * `Composition`: Center aligned, clean and minimalist.
    * `Style Modifier`: 2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background.
    * `Text Rendering`: "The Council"

#### Slide 02 :: VISION // 概念模型与愿景

`![[ppt-pages/P2.jpeg]]`

* **页面内容 (Slide Content)**:
    * **不再孤独决策**：用户 = 理事会主席 (Chairman)
    * **五大核心实体**：
        1.  **Group** (群组)：业务场景容器 (如 SaaS项目组)
        2.  **Agent** (理事)：拥有独立人设的 AI 专家
        3.  **Proposal** (提案)：任务输入 (Text/PDF)
        4.  **Workflow** (流程)：会议逻辑 DAG 图
        5.  **Session** (会议)：具体运行实例

* **演讲稿 (Script)**:
    为了实现这一愿景，我们构建了一套严谨的概念模型。在 The Council 中，你不再直接与模型对话，而是管理一个体系。最顶层是 **Group (群组)**，代表你的业务场景；在这个场景下，你指挥不同的 **Agent (理事)**；而他们如何协作，由你编排的 **Workflow (流程)** 决定。每一次思考，都是一次 **Session (会议)**。这不仅是聊天，这是组织管理学的 AI 化映射。

* **Image Generation Prompt**:
    * `Subject`: A hierarchical diagram made of floating glass layers. Top layer is a large glass plate labeled "Group". Below it, three portrait-shaped glass cards labeled "Agents". At the bottom, a flowing line connecting them representing "Workflow".
    * `Composition`: Isometric view showing the layers stacking.
    * `Style Modifier`: 2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background.
    * `Text Rendering`: "Group / Agent / Workflow"

#### Slide 03 :: CORE_1 // 弹性 AI 班底构建

`![[ppt-pages/P3.jpeg]]`

* **页面内容 (Slide Content)**:
    * **Agent Factory (理事工厂)**：模型无关 (Model Agnostic)
    * **多供应商混用**：
        * CFO 角色 -> **Claude-3.5** (强逻辑)
        * 气氛组角色 -> **GPT-4o-mini** (低成本)
    * **能力开关**：
        * ✅ **联网搜索 (Tavily/Serper)**：默认开启，杜绝幻觉
        * ❌ 代码执行：Phase 2 待启用

* **演讲稿 (Script)**:
    一个优秀的理事会需要多样化的人才。我们的 **Agent Factory** 允许你打破模型的围墙。你可以让逻辑严密的 Claude-3.5 担任 CFO，同时让成本低廉的 GPT-4o-mini 负责头脑风暴。我们支持为每一个 Agent 单独配置模型供应商和参数 (`Temperature`, `Top_P`)。并且，MVP 版本默认集成联网搜索能力，确保每一次发言都经过事实核查，而非单纯的生成。

* **Image Generation Prompt**:
    * `Subject`: Three distinct glass profile cards standing side-by-side. Each card has a different colored glowing chip inside (representing different AI models) and a specific icon (Finance, Code, Creative). Connectors show integration with a cloud symbol.
    * `Composition`: Horizontal alignment, focus on the variety of the "chips" inside the glass cards.
    * `Style Modifier`: 2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background.
    * `Text Rendering`: "Model Agnostic"

#### Slide 04 :: CORE_2 // 可视化流程编排

`![[ppt-pages/P4.jpeg]]`

* **页面内容 (Slide Content)**:
    * **Workflow Canvas**：基于 React Flow 的无限画布
    * **节点类型 (Nodes)**：
        * 🟢 Start (提案挂载)
        * 🟦 Agent (发言)
        * 🔶 Logic: **Parallel** (并行) / **Vote** (表决) / **Loop** (辩论循环)
        * 🆕 **FactCheck** (事实核查) & **HumanReview** (人机回环)
    * **God Mode (上帝模式)**：向导生成 + 全参数微调

* **演讲稿 (Script)**:
    这是整个系统的核心引擎——**Workflow Canvas**。我们将“开会”变成了可编程的有向无环图 (DAG)。你可以拖拽设计复杂的逻辑：让三个 Agent 并行思考，进入互相反驳的 Loop 循环，最后通过 Vote 节点表决。特别值得一提的是，我们引入了 **FactCheck (事实核查)** 和 **HumanReview (人类裁决)** 节点，强制性地将人类智慧与 AI 检索能力由点串联成线，确保决策的安全与可控。

* **Image Generation Prompt**:
    * `Subject`: A complex flowchart interface displayed on a tilted glass tablet. Nodes are glowing circular glass buttons connected by animated cyan light streams. Shows a "Parallel" split and a "Vote" convergence.
    * `Composition`: Close-up isometric view of the canvas interface.
    * `Style Modifier`: 2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background.
    * `Text Rendering`: "Workflow Canvas"

#### Slide 05 :: UX // 弹性会议室体验

`![[ppt-pages/P5.jpeg]]`

* **页面内容 (Slide Content)**:
    * **三分布局 (Flexible Layout)**：
        * Left (20%): 流程监控 (实时高亮 Token 消耗)
        * Mid (50%): 结构化对话 (并行 UI 气泡)
        * Right (30%): 文档上下文 (双向索引 `[Ref: P3]`)
    * **UI 降噪**：折叠/展开/全屏专注
    * **成本预估**：会前显示 **~$0.35** 预估，拒绝“账单刺客”

* **演讲稿 (Script)**:
    强大的后台需要优雅的前台。我们设计了“弹性会议室”。左侧实时监控流程节点与 Token 消耗；中间是结构化的对话流，支持多 Agent 并行发言的并排 UI；右侧则是深度集成的文档阅读器。当 AI 引用文档时，点击链接即可自动跳转高亮。更重要的是，我们在“开始会议”前提供 **成本预估模块**，让每一次 Token 的消耗都明明白白，彻底告别 API 账单焦虑。

* **Image Generation Prompt**:
    * `Subject`: A wide monitor screen divided into three vertical glass panes. Left pane shows a graph, middle pane shows chat bubbles (two side-by-side), right pane shows a document with highlighted text. Orange highlight on a "Cost: $0.35" tag.
    * `Composition`: Front-facing UI mockup on a glass surface.
    * `Style Modifier`: 2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background.
    * `Text Rendering`: "Flexible Workspace"

#### Slide 06 :: ENGINE // 双层记忆与隐私

`![[ppt-pages/P6.jpeg]]`

* **页面内容 (Slide Content)**:
    * **Dual-Layer Memory (双层记忆)**：
        1.  **Group Memory**: 场景隔离，严禁跨群泄露
        2.  **Agent Memory**: 角色经验沉淀
    * **Memory Sanitation (记忆清洗)**：
        * 会议复盘 -> 手动纠错 -> 向量库软删除
    * **Privacy First**:
        * **Docker PostgreSQL** 本地部署
        * API Key 浏览器加密存储

* **演讲稿 (Script)**:
    The Council 不仅协助决策，更在沉淀智慧。我们采用了 **双层记忆架构**：群组记忆确保业务上下文隔离，Agent 记忆积累角色经验。独有的 **记忆清洗机制** 允许你在复盘时标记错误结论，防止 AI 的“坏习惯”进入长期记忆。这一切都建立在绝对的隐私之上——所有数据存储在你私有的 Docker PostgreSQL 中，除了你，没有人能窥探你的理事会。

* **Image Generation Prompt**:
    * `Subject`: A secure glass vault or cube structure (`#0073e5` borders). Inside, data streams are organized into clean layers. A shield icon represents privacy. Outside the cube, a "Docker" whale icon in glass style.
    * `Composition`: Isometric view, emphasizing security and containment.
    * `Style Modifier`: 2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background.
    * `Text Rendering`: "Private Memory"

#### Slide 07 :: TECH // 技术栈与性能

`![[ppt-pages/P7.jpeg]]`

* **页面内容 (Slide Content)**:
    * **Backend**: **Go (Goroutines)** - 高并发支持 3+ Agent 同时生成
    * **Frontend**: React SPA + WebSocket (< 500ms 延迟)
    * **Database**: PostgreSQL + pgvector (Dockerized)
    * **Extensibility**: OpenAI Function Calling 标准兼容

* **演讲稿 (Script)**:
    为了支撑复杂的并发思考，我们选择了 **Go** 作为后端核心，利用 Goroutines 轻松处理多个 Agent 的并行生成。前端通过 WebSocket 实现了低于 500ms 的首字延迟，带来行云流水的体验。整个技术栈基于 Docker 容器化，既保证了部署的便捷性，也为未来的插件扩展 (Function Calling) 打下了坚实的工程基础。

* **Image Generation Prompt**:
    * `Subject`: Floating glass blocks with logos engraved: "Go" (Gopher), "React" (Atom), "PostgreSQL" (Elephant). They are connected by high-speed glowing cyan data pipes.
    * `Composition`: Dynamic, showing speed and connectivity.
    * `Style Modifier`: 2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background.
    * `Text Rendering`: "Go + React + Docker"

#### Slide 08 :: CLOSING // 你的私有智库

`![[ppt-pages/P8.jpeg]]`

* **页面内容 (Slide Content)**:
    * **The Council v1.2.0**
    * **Empower Your Decision.**
    * **Build Your Board.**
    * **Own Your Wisdom.**
    * Next Step: **Deploy on Docker**

* **演讲稿 (Script)**:
    从今天起，不再是一个人面对复杂的代码审查，不再是一个人苦思商业计划的漏洞。组建你的群组，定义你的理事，编排你的流程。The Council —— 让 AI 成为你的幕僚，让智慧成为你的资产。现在，就可以通过 Docker 启动属于你的私有智库。谢谢大家。

* **Image Generation Prompt**:
    * `Subject`: An open glass hand offering a glowing, complex geometric polygon (symbolizing the "Wisdom" or "Product"). The polygon is bright Orange (`#ff531a`) contrasting with the Blue/White environment.
    * `Composition`: Center, symbolic, inspiring.
    * `Style Modifier`: 2.5D glassmorphism layout, floating frosted glass cards, clean flat vector icons, soft depth of field, subtle neural mesh background, high tech business style, white background.
    * `Text Rendering`: "The Council"