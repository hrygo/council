export interface Group {
    id: string;
    name: string;
    icon: string;
    system_prompt: string;
    default_agent_ids: string[];
    created_at: string;
    updated_at: string;
}

export interface CreateGroupInput {
    name: string;
    icon?: string;
    system_prompt?: string;
    default_agent_ids?: string[];
}
