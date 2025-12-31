package workflow

// PassthroughConfig 定义透传配置。
// 这是一个通用结构，不包含任何业务字段名。
type PassthroughConfig struct {
	// Keys 需要从 input 透传到 output 的字段名列表
	Keys []string
}

// ApplyPassthrough 将配置中指定的字段从 input 复制到 output。
// output 必须已初始化。
// 这是一个纯工具函数，不包含业务逻辑。
func ApplyPassthrough(input, output map[string]interface{}, config PassthroughConfig) {
	if output == nil {
		return
	}
	for _, key := range config.Keys {
		if val, ok := input[key]; ok {
			output[key] = val
		}
	}
}

// PromptSection defines a section in the user prompt constructed from input.
type PromptSection struct {
	Key   string
	Label string
}
