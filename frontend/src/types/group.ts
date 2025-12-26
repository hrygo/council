export interface Group {
    group_uuid: string;
    name: string;
    icon: string;
    system_prompt: string;
    default_agent_uuids: string[];
    created_at: string;
    updated_at: string;
}

export interface CreateGroupInput {
    name: string;
    icon?: string;
    system_prompt?: string;
    default_agent_uuids?: string[];
}
