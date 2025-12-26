import useSWR from 'swr';

export interface ProviderOption {
    provider_id: string;
    name: string;
    icon: string;
    models: string[];
}

interface LLMOptionsResponse {
    providers: ProviderOption[];
}

const fetcher = (url: string) => fetch(url).then(res => res.json());

export const useLLMOptions = () => {
    const { data, error, isLoading } = useSWR<LLMOptionsResponse>('/api/v1/llm/providers', fetcher, {
        revalidateOnFocus: false,
        revalidateOnReconnect: false,
    });

    return {
        providers: data?.providers || [],
        isLoading,
        isError: error
    };
};
