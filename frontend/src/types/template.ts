import type { BackendGraph } from '../utils/graphUtils';

export type TemplateCategory = 'code_review' | 'business_plan' | 'quick_decision' | 'custom' | 'other';

export interface Template {
    template_uuid: string;
    name: string;
    description: string;
    category: TemplateCategory;
    is_system: boolean;
    graph: BackendGraph;
    created_at?: string;
    updated_at?: string;
}

export interface CreateTemplateInput {
    name: string;
    description: string;
    category: TemplateCategory;
    graph: BackendGraph;
}
