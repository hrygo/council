# Backend Specifications: 节点处理器与 API

> **并行 Sprint**: Week 1-4  
> **目标**: 实现所有节点处理器和新增 API

---

## 📋 Backend Specs 索引

| Spec ID  | 文档                                                         | 类型        | 优先级 | Sprint |
| -------- | ------------------------------------------------------------ | ----------- | ------ | ------ |
| SPEC-401 | [SequenceProcessor](./SPEC-401-sequence-processor.md)        | Processor   | P1     | 1-2    |
| SPEC-402 | [VoteProcessor](./SPEC-402-vote-processor.md)                | Processor   | P1     | 1-2    |
| SPEC-403 | [LoopProcessor](./SPEC-403-loop-processor.md)                | Processor   | P2     | 1-2    |
| SPEC-404 | [FactCheckProcessor](./SPEC-404-factcheck-processor.md)      | Processor   | P1     | 3-4    |
| SPEC-405 | [HumanReviewProcessor](./SPEC-405-human-review-processor.md) | Processor   | P0     | 3-4    |
| SPEC-406 | [Templates API](./SPEC-406-templates-api.md)                 | API         | P1     | 3-4    |
| SPEC-407 | [Cost Estimation API](./SPEC-407-cost-estimation-api.md)     | API         | P1     | 3-4    |
| SPEC-408 | [三层记忆协议](./SPEC-408-memory-protocol.md)                | Core        | P0     | 3-4    |
| SPEC-409 | [逻辑熔断](./SPEC-409-circuit-breaker.md)                    | Core        | P0     | 3-4    |
| SPEC-410 | [防幻觉传播](./SPEC-410-anti-hallucination.md)               | Core        | P1     | 3-4    |
| SPEC-411 | [联网搜索集成](./SPEC-411-search-integration.md)             | Integration | P1     | 1-2    |

---

## 🎯 验收标准

- [ ] 所有节点处理器实现并通过测试
- [ ] Templates CRUD API 可用
- [ ] Cost Estimation API 返回合理预估
