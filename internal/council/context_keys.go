// Package council 包含 Council Debate/Optimize 工作流的应用层实现。
// 此包依赖 internal/core/workflow (骨架层)，但反之不可。
package council

// CouncilContextKeys 定义 Council 工作流使用的所有上下文字段。
// 这些字段名是业务概念，不应出现在骨架层代码中。
var CouncilContextKeys = []string{
	"document_content",       // 原始文档内容
	"proposal",               // 方案摘要
	"optimization_objective", // 优化目标
	"attachments",            // 附件列表
	"combined_context",       // 合并后的附件内容
	"session_id",             // 会话 ID
	"aggregated_outputs",     // 聚合的 Agent 输出
	"agent_output",           // 单个 Agent 输出
	"history_context",        // 历史上下文
}

// AgentPassthroughKeys 定义 Agent 节点需要透传的字段。
var AgentPassthroughKeys = []string{
	"document_content",
	"proposal",
	"optimization_objective",
	"attachments",
	"combined_context",
	"session_id",
	"aggregated_outputs",
}

// LoopPassthroughKeys 定义 Loop 节点需要透传的字段。
var LoopPassthroughKeys = []string{
	"document_content",
	"proposal",
	"optimization_objective",
	"combined_context",
	"session_id",
}

// StartOutputKeys 定义 Start 节点需要输出的字段。
var StartOutputKeys = []string{
	"document_content",
	"proposal",
	"optimization_objective",
	"attachments",
	"combined_context",
}

// EndInputKeys 定义 End 节点需要读取的字段。
var EndInputKeys = []string{
	"document_content",
	"combined_context",
	"proposal",
	"aggregated_outputs",
	"agent_output",
}
