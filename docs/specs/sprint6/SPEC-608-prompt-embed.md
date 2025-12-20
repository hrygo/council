# SPEC-608: Prompt Embed 机制

> **优先级**: P0  
> **类型**: Infrastructure  
> **预估工时**: 4h

## 1. 概述

使用 Go 的 `//go:embed` 机制将 Prompt 文件嵌入到二进制中，替代 SQL 内嵌长文本的方案。

## 2. 目标

- Prompt 保持 `.md` 格式，便于编辑和版本控制
- 避免 SQL 转义问题
- Migration 时动态读取 Prompt 内容并插入数据库

## 3. 技术实现

### 3.1 目录结构

```
internal/resources/
  prompts/
    system_affirmative.md   # 完整 Prompt (从 example/prompts 复制)
    system_negative.md
    system_adjudicator.md
  embed.go                  # Go embed 定义
  seeder.go                 # 数据库初始化逻辑
```

### 3.2 Prompt 文件格式

```markdown
<!-- internal/resources/prompts/system_affirmative.md -->
---
name: Value Defender
provider: gemini
model: gemini-3-pro-preview
temperature: 0.9
max_tokens: 8192
capabilities:  # Issue Fix: Support dynamic capabilities
  web_search: false
  code_execution: false
---

### Role
...
```

### 3.3 Go Embed 定义

```go
// internal/resources/embed.go
package resources

import "embed"

//go:embed prompts/*.md
var PromptFiles embed.FS
```

### 3.4 Prompt 解析器

```go
// internal/resources/prompt_loader.go
package resources

import (
    "bytes"
    "fmt"
    "io/fs"
    "strings"
    
    "gopkg.in/yaml.v3"
)

type AgentConfig struct {
    Name         string          `yaml:"name"`
    Provider     string          `yaml:"provider"`
    Model        string          `yaml:"model"`
    Temperature  float64         `yaml:"temperature"`
    MaxTokens    int             `yaml:"max_tokens"`
    Capabilities map[string]bool `yaml:"capabilities"` // Updated
}

type AgentPrompt struct {
    Config  AgentConfig
    Content string
}

// LoadAllPrompts dynamic file loading (Issue Fix)
func LoadAllPrompts() (map[string]*AgentPrompt, error) {
    prompts := make(map[string]*AgentPrompt)
    
    entries, err := PromptFiles.ReadDir("prompts")
    if err != nil {
        return nil, err
    }
    
    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
            continue
        }
        
        prompt, err := LoadPrompt(entry.Name())
        if err != nil {
            return nil, err
        }
        
        // Remove .md extension for ID
        id := strings.TrimSuffix(entry.Name(), ".md")
        prompts[id] = prompt
    }
    
    return prompts, nil
}

func LoadPrompt(filename string) (*AgentPrompt, error) {
    // ... implementation same as before ...
}
```

### 3.5 数据库 Seeder

```go
// internal/resources/seeder.go

func (s *Seeder) SeedAgents(ctx context.Context) error {
    // Dynamic loading from loader
    prompts, err := LoadAllPrompts()
    if err != nil {
        return err
    }
    
    for agentID, prompt := range prompts {
        modelConfig, _ := json.Marshal(map[string]interface{}{
            "provider":    prompt.Config.Provider,
            "model":       prompt.Config.Model,
            "temperature": prompt.Config.Temperature,
            "max_tokens":  prompt.Config.MaxTokens,
        })
        
        // Use capabilities from prompt config (Issue Fix)
        // Default to safe defaults if nil
        caps := prompt.Config.Capabilities
        if caps == nil {
            caps = map[string]bool{
                "web_search":     false, 
                "code_execution": false,
            }
        }
        capabilities, _ := json.Marshal(caps)
        
        _, err = s.db.ExecContext(ctx, `
            INSERT INTO agents (id, name, persona_prompt, model_config, capabilities, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
            ON CONFLICT (id) DO NOTHING
        `, agentID, prompt.Config.Name, prompt.Content, modelConfig, capabilities)
        
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### 3.6 集成到启动流程

```go
// cmd/council/main.go
func main() {
    // ... 初始化 DB ...
    
    seeder := resources.NewSeeder(db)
    if err := seeder.SeedAgents(context.Background()); err != nil {
        log.Fatalf("Failed to seed agents: %v", err)
    }
    
    // ... 启动服务 ...
}
```

## 4. Prompt 文件内容

### 4.1 system_affirmative.md

从 `example/prompts/affirmative.md` 完整复制，添加 YAML Front Matter。

### 4.2 system_negative.md

从 `example/prompts/negative.md` 完整复制，添加 YAML Front Matter。

### 4.3 system_adjudicator.md

从 `example/prompts/adjudicator.md` 完整复制，添加 YAML Front Matter。

#### 增强说明 (Issue 5 Remediation)

> **位置明确**: 以下评分矩阵指引应追加到 `internal/resources/prompts/system_adjudicator.md` 文件的 **Output Format 部分之后**，作为 Adjudicator 输出的强制格式。
>
> **来源**: 此增强内容是 `skill.md` Step 3 (Verify Consistency) 在 Workflow 系统中的实现，用于支持 HumanReview 节点的评分显示和退出条件判断。

**追加内容** (在 `### 6. 元裁决（Meta-Verdict）` 部分之后):

```markdown
---

## 🎯 结构化评分输出 (Structured Score Output)

> **系统解析区**: 以下 JSON 块将被 Workflow Engine 解析，用于驱动循环退出条件。

\`\`\`json
{
  "score": {
    "strategic_alignment": XX,
    "practical_value": XX,
    "logical_consistency": XX,
    "weighted_total": XX
  },
  "verdict": "直接通过 | 细节优化 | 逻辑完善 | 深度重构 | 彻底驳回",
  "exit_recommendation": true | false
}
\`\`\`

### 评分矩阵

| 维度     | 权重 | 得分 (0-100) | 说明 |
| -------- | ---- | ------------ | ---- |
| 战略对齐 | 40%  | ?            | ...  |
| 实用价值 | 30%  | ?            | ...  |
| 逻辑一致 | 30%  | ?            | ...  |

**综合得分**: ?

### 行动建议

- [ ] 继续优化 (Score < 90)
- [ ] 直接通过 (Score >= 90)
```

#### Workflow Engine 集成

```go
// internal/core/workflow/nodes/agent.go
// Adjudicator 输出后，解析 JSON 块
type StructuredScore struct {
    Score struct {
        StrategicAlignment int `json:"strategic_alignment"`
        PracticalValue     int `json:"practical_value"`
        LogicalConsistency int `json:"logical_consistency"`
        WeightedTotal      int `json:"weighted_total"`
    } `json:"score"`
    Verdict           string `json:"verdict"`
    ExitRecommendation bool  `json:"exit_recommendation"`
}

func parseAdjudicatorOutput(content string) (*StructuredScore, error) {
    // Extract JSON block from markdown using regex
    re := regexp.MustCompile(`(?s)\x60\x60\x60json\s*(\{.*?\})\s*\x60\x60\x60`)
    matches := re.FindStringSubmatch(content)
    if len(matches) < 2 {
        return nil, fmt.Errorf("no structured score found")
    }
    
    var score StructuredScore
    if err := json.Unmarshal([]byte(matches[1]), &score); err != nil {
        return nil, err
    }
    return &score, nil
}
```

## 5. 验收标准

- [ ] `internal/resources/prompts/` 目录存在且包含 3 个 `.md` 文件
- [ ] Prompt 文件格式正确 (YAML Front Matter + Markdown Body)
- [ ] `LoadPrompt()` 函数可正确解析 Prompt 文件
- [ ] `SeedAgents()` 可将 Prompt 插入数据库
- [ ] 服务启动后，数据库中 Agent 的 `persona_prompt` 字段完整

## 6. 测试

```go
func TestLoadPrompt(t *testing.T) {
    prompt, err := LoadPrompt("system_affirmative.md")
    
    assert.NoError(t, err)
    assert.Equal(t, "Value Defender", prompt.Config.Name)
    assert.Equal(t, "gemini", prompt.Config.Provider)
    assert.Contains(t, prompt.Content, "价值辩护人")
}
```

## 7. 与 SPEC-601 关系

**SPEC-601 更新**: 原有的 SQL Migration 仅创建表结构，Prompt 数据通过本 Seeder 注入。

Migration SQL 简化为：
```sql
-- 确保 agents 表存在 (已在 schema migration 中)
-- Seed data 由 Go Seeder 处理
```
