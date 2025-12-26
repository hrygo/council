import React, { useEffect, useState } from 'react';

interface Knowledge {
  knowledge_uuid: string;
  title: string;
  summary: string;
  source: string;
  relevance: number;
  timestamp: Date;
  layer: 'sandboxed' | 'working' | 'long-term';
}

interface KnowledgePanelProps {
  sessionId: string;
}

export const KnowledgePanel: React.FC<KnowledgePanelProps> = ({ sessionId }) => {
  const [knowledge, setKnowledge] = useState<Knowledge[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [layerFilter, setLayerFilter] = useState<'all' | 'sandboxed' | 'working' | 'long-term'>('all');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!sessionId) return;

    const fetchKnowledge = async () => {
      setIsLoading(true);
      try {
        const params = new URLSearchParams();
        if (layerFilter !== 'all') params.append('layer', layerFilter);

        const response = await fetch(`/api/v1/sessions/${sessionId}/knowledge?${params}`);
        if (response.ok) {
          const data = await response.json();
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          const mappedItems = (data.items || []).map((item: any) => ({
            knowledge_uuid: item.knowledge_uuid,
            title: item.title,
            summary: item.summary,
            source: item.content, // Using content as source/details
            relevance: item.relevance_score / 5, // Normalize back if needed or use score directly
            timestamp: new Date(item.created_at),
            layer: item.memory_layer
          }));
          setKnowledge(mappedItems);
        }
      } catch (error: unknown) {
        console.error('Failed to fetch knowledge:', error instanceof Error ? error.message : String(error));
      } finally {
        setIsLoading(false);
      }
    };

    fetchKnowledge();
  }, [sessionId, layerFilter]);

  const filteredKnowledge = knowledge.filter(item =>
    searchQuery === '' ||
    item.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.summary.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="h-full flex flex-col bg-white dark:bg-gray-900 border-l border-gray-200 dark:border-gray-700">
      {/* Header */}
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <h3 className="text-lg font-semibold mb-3 text-gray-900 dark:text-gray-100">
          📚 相关知识
        </h3>

        {/* Search */}
        <input
          type="text"
          placeholder="搜索知识..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />

        {/* Layer Filter */}
        <select
          value={layerFilter}
          onChange={(e) => setLayerFilter(e.target.value as Knowledge['layer'] | 'all')}
          className="w-full mt-2 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="all">全部记忆</option>
          <option value="sandboxed">隔离区</option>
          <option value="working">工作记忆</option>
          <option value="long-term">长期记忆</option>
        </select>
      </div>

      {/* Knowledge List */}
      <div className="flex-1 overflow-y-auto">
        {isLoading ? (
          <div className="flex items-center justify-center h-32">
            <div className="text-gray-500 dark:text-gray-400">加载中...</div>
          </div>
        ) : filteredKnowledge.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-32 text-gray-500 dark:text-gray-400">
            <p className="text-sm">暂无相关知识</p>
            <p className="text-xs mt-1">知识将在会话进行中自动更新</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-100 dark:divide-gray-800">
            {filteredKnowledge.map((item) => (
              <KnowledgeItem key={item.knowledge_uuid} knowledge={item} />
            ))}
          </div>
        )}
      </div>

      {/* Footer */}
      {!isLoading && filteredKnowledge.length > 0 && (
        <div className="p-2 border-t border-gray-200 dark:border-gray-700 text-xs text-gray-500 dark:text-gray-400 text-center">
          显示 {filteredKnowledge.length} / {knowledge.length} 项
        </div>
      )}
    </div>
  );
};

interface KnowledgeItemProps {
  knowledge: Knowledge;
}

const KnowledgeItem: React.FC<KnowledgeItemProps> = ({ knowledge }) => {
  const getLayerBadgeColor = (layer: string) => {
    switch (layer) {
      case 'sandboxed': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      case 'working': return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case 'long-term': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      default: return 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200';
    }
  };

  const getLayerLabel = (layer: string) => {
    switch (layer) {
      case 'sandboxed': return '隔离';
      case 'working': return '工作';
      case 'long-term': return '长期';
      default: return '未知';
    }
  };

  const getRelevanceStars = (relevance: number) => {
    const stars = Math.round(relevance * 5);
    return '⭐'.repeat(Math.max(1, Math.min(5, stars)));
  };

  return (
    <div className="p-4 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors">
      {/* Title and Badge */}
      <div className="flex items-start justify-between mb-2">
        <h4 className="text-sm font-medium text-gray-900 dark:text-gray-100 flex-1 line-clamp-2">
          {knowledge.title}
        </h4>
        <span className={`ml-2 px-2 py-0.5 text-xs rounded ${getLayerBadgeColor(knowledge.layer)}`}>
          {getLayerLabel(knowledge.layer)}
        </span>
      </div>

      {/* Summary */}
      <p className="text-xs text-gray-600 dark:text-gray-400 mb-2 line-clamp-2">
        {knowledge.summary}
      </p>

      {/* Metadata */}
      <div className="flex items-center justify-between text-xs text-gray-500 dark:text-gray-500">
        <span className="flex items-center">
          <span className="mr-1">相关度:</span>
          <span>{getRelevanceStars(knowledge.relevance)}</span>
        </span>
        <span>{knowledge.source}</span>
      </div>
    </div>
  );
};
