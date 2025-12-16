# 1. 总体架构设计 (High-Level Architecture)

采用 **"模式分离 (Mode Separation)"** 与 **"重后端"** 策略。前端负责两种互斥视图的渲染，后端负责全生命周期的核心逻辑与防御机制。

```mermaid
graph TD
    subgraph Browser["Presentation Layer - Browser"]
        subgraph Mode_Build["Mode A: Builder (The IDE)"]
            UI_Canvas["React Flow Editor"]
            UI_PropPanel["Property Panel"]
        end
        
        subgraph Mode_Run["Mode B: Runner (The Meeting)"]
            UI_Timeline["Linear Timeline View"]
            UI_Stream["Chat Stream (DOM/SVG)"]
            UI_Dashboard["Cost Dashboard"]
        end
        
        Store["Zustand Store (Session/Layout)"]
    end

    subgraph Communication["Communication Layer - HTTPS/WSS"]
        API_HTTP["REST API"]
        API_WS["WebSocket Stream"]
    end

    subgraph Backend["Application Layer - Go Backend"]
        Gateway["Gin Server"]
        
        subgraph Middlewares["Defense Middlewares"]
            CircuitBreaker["Logic Circuit Breaker"]
            Auth["Auth & Encryption"]
        end
        
        subgraph Engines["Core Engines"]
            WorkflowEngine["Workflow Scheduler"]
            AIGateway["AI Model Router"]
            MemoryManager["Three-Tier Memory Protocol"]
        end
        
        subgraph Cache["Working Memory"]
            Redis["Hot Cache (24h TTL)"]
        end
    end

    subgraph Data["Data Layer - Docker"]
        DB[("PostgreSQL + pgvector")]
        FileSys["File System"]
    end

    UI_Canvas --> API_HTTP
    UI_Stream <--> API_WS
    
    API_HTTP --> Gateway
    API_WS --> Gateway
    Gateway --> Middlewares --> Engines
    
    WorkflowEngine --> AIGateway
    WorkflowEngine --> MemoryManager
    
    MemoryManager --> Redis
    MemoryManager --> DB
    AIGateway --> DB
```

### 关键组件职责更新

1.  **双模前端 (Dual-Mode Frontend)**:
    *   **Builder**: 重型 WebGL 画布，加载完整 React Flow 库。
    *   **Runner**: 轻量级 DOM 渲染，卸载 Canvas 逻辑以优化长对话性能。

2.  **防御性中间件 (Defense Middlewares)**:
    *   **Logic Circuit Breaker**: 在进入 Workflow Engine 前拦截异常流量（死循环/Token 激增）。

3.  **三层记忆架构**:
    *   引入 **Redis** 作为第2层工作记忆 (Working Memory) 的存储介质。
