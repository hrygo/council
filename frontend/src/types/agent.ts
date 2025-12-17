export interface ModelConfig {
    provider: 'openai' | 'anthropic' | 'google' | 'deepseek' | 'dashscope';
    model: string;
    temperature: number;
    top_p: number;
    max_tokens: number;
}

export interface Capabilities {
    web_search: boolean;
    search_provider: string; // e.g. "tavily"
    code_execution: boolean;
}

export interface Agent {
    id: string;
    name: string;
    avatar: string; // URL or emoji char
    description: string;
    persona_prompt: string;
    model_config: ModelConfig;
    capabilities: Capabilities;
    created_at: string;
    updated_at: string;
}

export interface CreateAgentInput {
    name: string;
    avatar?: string;
    description?: string;
    persona_prompt: string;
    model_config: ModelConfig;
    capabilities: Capabilities;
}
