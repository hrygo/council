# 开发计划模板库

本目录包含创建和维护开发计划所需的标准模板和工具。

---

## 📚 模板清单

| 模板文件 | 用途 | 适用场景 |
| :------- | :--- | :------- |
| [development_plan_template.md](./development_plan_template.md) | 开发计划主文档模板 | 新项目启动、重大版本规划 |
| [spec_template.md](./spec_template.md) | 规格说明文档模板 | 编写任何新的 Spec 文档 |
| [sprint_readme_template.md](./sprint_readme_template.md) | Sprint README 模板 | 每个 Sprint 的索引和跟踪 |
| [development_plan_checklist.md](./development_plan_checklist.md) | 质量检查清单 | 文档创建、更新、评审 |

---

## 🚀 快速开始

### 场景 1: 创建新项目的开发计划

1. 复制 `development_plan_template.md` 到 `docs/development_plan.md`
2. 填写头部元信息（策略、质量方针）
3. 定义里程碑和时间节点
4. 规划 Sprint 并分解任务
5. 使用 `development_plan_checklist.md` 进行自检
6. 提交评审

### 场景 2: 创建新的 Sprint

1. 在 `docs/specs/` 下创建 `sprintX/` 目录
2. 复制 `sprint_readme_template.md` 到 `sprintX/README.md`
3. 填写 Sprint 目标和验收标准
4. 为每个任务创建 Spec 文档（使用 `spec_template.md`）
5. 在主开发计划中添加 Sprint 详情章节
6. 更新任务矩阵和规格索引

### 场景 3: 编写规格文档

1. 复制 `spec_template.md` 到相应的 Sprint 目录
2. 重命名为 `SPEC-XXX-描述.md`（遵循命名规范）
3. 填写所有必选章节
4. 在主开发计划的任务矩阵和规格索引中添加引用
5. 确保与相关任务和依赖保持一致

### 场景 4: 更新开发计划

1. 更新任务状态（✅/🔄/⬜）
2. 重新计算 Sprint 完成度
3. 更新验收标准完成情况
4. 记录变更日志
5. 使用检查清单验证一致性

---

## 📖 编号规范

### Sprint 编号

- **格式**: S1, S2, S3, ...
- **说明**: 从 1 开始递增，独立于时间

### Spec 编号

| 范围 | 编号区间 | 示例 |
| :--- | :------- | :--- |
| Sprint 1 | 001-099 | SPEC-001, SPEC-005 |
| Sprint 2 | 100-199 | SPEC-101, SPEC-105 |
| Sprint 3 | 200-299 | SPEC-201, SPEC-206 |
| Sprint 4 | 300-399 | SPEC-301, SPEC-305 |
| Backend | 400-499 | SPEC-401, SPEC-411 |
| Post-MVP | 500-599 | SPEC-501, SPEC-505 |
| Sprint 6 | 600-699 | SPEC-601, SPEC-609 |
| Sprint 7 | 700-799 | SPEC-701, SPEC-703 |
| Sprint 8+ | 800-899 | SPEC-801, SPEC-804 |

### 任务 ID

| 类型 | 格式 | 示例 |
| :--- | :--- | :--- |
| 前端任务 | {Sprint序号}.{序号} | 1.1, 2.3, 3.5 |
| 后端任务 | B.{序号} | B.1, B.5, B.11 |
| 专项任务 | {主题首字母}.{序号} | 5.1 (Post-MVP), Q.1 (QA) |

---

## 🎯 优先级定义

| 级别 | 说明 | 使用场景 |
| :--: | :--- | :------- |
| P0 | 阻塞性任务，必须在当前 Sprint 完成 | 核心功能、严重 Bug |
| P1 | 关键任务，影响核心功能 | 重要功能、性能优化 |
| P2 | 重要任务，可在必要时延期 | 次要功能、体验优化 |
| P3 | 优化任务，可推迟到后续迭代 | Nice-to-have、技术债 |

---

## ✅ 状态符号

| 符号 | 含义 | 适用对象 |
| :--: | :--- | :------- |
| ✅ | 已完成 | 任务、Spec、里程碑、验收标准 |
| 🔄 | 进行中 | Sprint、任务、问题处理 |
| ⬜ | 未开始 | 任务、Spec、里程碑 |
| 📅 | 计划中 | Sprint、未来任务 |
| ❌ | 已取消 | 任务、Spec |

---

## 📐 Spec 文档类型

| 类型 | 说明 | 示例 |
| :--- | :--- | :--- |
| Feature | 新功能开发 | Session Store、Workflow Builder |
| Bugfix | 缺陷修复 | WebSocket 连接修复 |
| Refactor | 代码重构 | LLM Registry 重构 |
| QA | 质量保证 | E2E 测试、性能优化 |
| Security | 安全强化 | 权限控制、数据加密 |
| Infrastructure | 基础设施 | Prompt Embed、CI/CD |

---

## 🛠️ 推荐工具

### 编辑器

- **VSCode** + Markdown Preview Enhanced
- **Typora**（所见即所得）
- **Obsidian**（知识管理）

### VSCode 插件

- Markdown All in One
- Markdown Preview Enhanced
- markdownlint
- Markdown Table
- Mermaid Markdown Syntax Highlighting

### 自动化工具

```bash
# 验证开发计划
python3 scripts/validate_dev_plan.py docs/development_plan.md

# 检查 Markdown 链接有效性
npx markdown-link-check docs/**/*.md

# 格式检查
npx markdownlint docs/**/*.md
```

---

## 📋 检查清单速查

### 创建新开发计划

- [ ] 填写策略和质量方针
- [ ] 定义所有里程碑
- [ ] 创建进度总览表
- [ ] 规划至少 3 个 Sprint
- [ ] 建立规格文档目录结构

### 创建新 Sprint

- [ ] 在 specs/ 下创建 sprintX/ 目录
- [ ] 创建 Sprint README
- [ ] 定义 Sprint 目标
- [ ] 分解任务和 Spec
- [ ] 绘制依赖关系图
- [ ] 定义验收标准

### 编写 Spec 文档

- [ ] 使用标准模板
- [ ] 填写所有必选章节
- [ ] 定义清晰的验收标准
- [ ] 标注依赖关系
- [ ] 在开发计划中建立索引

### 更新开发计划

- [ ] 更新任务状态
- [ ] 重新计算完成度
- [ ] 更新验收标准
- [ ] 记录变更日志
- [ ] 验证一致性

---

## 🎓 最佳实践

### 1. 任务分解

- 单个任务 2-8 小时
- 任务边界清晰
- 输入输出明确
- 低耦合，便于并行

### 2. 验收标准

- 可测试、可量化
- 独立验证
- 引用具体文件或命令
- 分类清晰

### 3. 依赖管理

- 明确前置依赖
- 标注阻塞关系
- 可视化依赖图
- 定期验证

### 4. 文档维护

- 及时更新状态
- 保持一致性
- 单一数据源
- 使用链接引用

### 5. 版本控制

- 记录变更日志
- 标注版本号
- 追踪决策依据
- 归档历史版本

---

## 🆘 常见问题

### Q: 如何选择合适的优先级？

**A**: 遵循以下原则：
- P0: 阻塞整个 Sprint 或里程碑
- P1: 影响核心功能，但可以 workaround
- P2: 重要但可延期
- P3: Nice-to-have

### Q: Sprint 工时如何估算？

**A**: 
- 单个任务：2-8 小时
- Sprint 总工时：建议不超过 40 小时
- 预留 20% 缓冲时间应对风险

### Q: 如何处理 Spec 编号冲突？

**A**: 
- 严格按照编号规则分配
- 使用编号区间隔离不同 Sprint
- 发现冲突时立即重新编号

### Q: 如何保持文档一致性？

**A**: 
- 使用检查清单定期验证
- 建立自动化检查（CI）
- 实施代码评审机制
- 单一数据源原则

---

## 📞 获取帮助

- **项目规范**: [GEMINI.md](../../GEMINI.md)
- **开发指南**: [TDD_GUIDE.md](../references/TDD_GUIDE.md)
- **PRD 文档**: [PRD.md](../references/PRD.md)

---

## 🔄 模板更新记录

| 日期 | 版本 | 变更内容 |
| :--- | :--- | :------- |
| 2025-12-26 | v1.0 | 初始模板库创建 |
