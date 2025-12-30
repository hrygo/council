import { type FC, useState, useMemo } from 'react';
import { FileCode, ArrowRight, Plus, Minus } from 'lucide-react';

interface DiffLine {
    type: 'unchanged' | 'added' | 'removed';
    content: string;
    lineNumber: { old?: number; new?: number };
}

interface CodeDiffViewerProps {
    originalContent: string;
    newContent: string;
    filePath: string;
    onContentChange?: (newContent: string) => void;
    editable?: boolean;
}

function computeDiff(original: string, modified: string): DiffLine[] {
    const originalLines = original.split('\n');
    const modifiedLines = modified.split('\n');
    const result: DiffLine[] = [];

    // Simple line-by-line diff (for production, use a proper diff algorithm)
    const maxLen = Math.max(originalLines.length, modifiedLines.length);
    let oldLineNum = 1;
    let newLineNum = 1;

    for (let i = 0; i < maxLen; i++) {
        const oldLine = originalLines[i];
        const newLine = modifiedLines[i];

        if (oldLine === undefined && newLine !== undefined) {
            // Added line
            result.push({
                type: 'added',
                content: newLine,
                lineNumber: { new: newLineNum++ },
            });
        } else if (newLine === undefined && oldLine !== undefined) {
            // Removed line
            result.push({
                type: 'removed',
                content: oldLine,
                lineNumber: { old: oldLineNum++ },
            });
        } else if (oldLine === newLine) {
            // Unchanged
            result.push({
                type: 'unchanged',
                content: oldLine,
                lineNumber: { old: oldLineNum++, new: newLineNum++ },
            });
        } else {
            // Changed: show as removed then added
            result.push({
                type: 'removed',
                content: oldLine,
                lineNumber: { old: oldLineNum++ },
            });
            result.push({
                type: 'added',
                content: newLine,
                lineNumber: { new: newLineNum++ },
            });
        }
    }

    return result;
}

export const CodeDiffViewer: FC<CodeDiffViewerProps> = ({
    originalContent,
    newContent,
    filePath,
    onContentChange,
    editable = false,
}) => {
    const [localContent, setLocalContent] = useState(newContent);
    const displayContent = editable ? localContent : newContent;

    const diffLines = useMemo(
        () => computeDiff(originalContent, displayContent),
        [originalContent, displayContent]
    );

    const stats = useMemo(() => {
        const added = diffLines.filter(l => l.type === 'added').length;
        const removed = diffLines.filter(l => l.type === 'removed').length;
        return { added, removed };
    }, [diffLines]);

    const handleContentChange = (value: string) => {
        setLocalContent(value);
        onContentChange?.(value);
    };

    return (
        <div className="h-full flex flex-col bg-gray-900 rounded-lg overflow-hidden">
            {/* Header */}
            <div className="flex items-center justify-between px-4 py-2 bg-gray-800 border-b border-gray-700">
                <div className="flex items-center gap-2 text-gray-300">
                    <FileCode size={16} className="text-blue-400" />
                    <span className="text-sm font-mono">{filePath}</span>
                </div>
                <div className="flex items-center gap-3 text-xs">
                    <span className="flex items-center gap-1 text-green-400">
                        <Plus size={12} />
                        {stats.added}
                    </span>
                    <span className="flex items-center gap-1 text-red-400">
                        <Minus size={12} />
                        {stats.removed}
                    </span>
                </div>
            </div>

            {/* Diff View */}
            <div className="flex-1 overflow-auto">
                {editable ? (
                    <div className="flex h-full">
                        {/* Original (read-only) */}
                        <div className="w-1/2 border-r border-gray-700 overflow-auto">
                            <div className="px-2 py-1 bg-gray-800 text-xs text-gray-400 border-b border-gray-700 sticky top-0">
                                Original
                            </div>
                            <pre className="text-sm font-mono p-2 text-gray-400 whitespace-pre-wrap">
                                {originalContent}
                            </pre>
                        </div>
                        {/* Modified (editable) */}
                        <div className="w-1/2 overflow-auto flex flex-col">
                            <div className="px-2 py-1 bg-gray-800 text-xs text-gray-400 border-b border-gray-700 sticky top-0">
                                Modified (Editable)
                            </div>
                            <textarea
                                value={localContent}
                                onChange={(e) => handleContentChange(e.target.value)}
                                className="flex-1 w-full bg-gray-900 text-gray-100 text-sm font-mono p-2 resize-none focus:outline-none focus:ring-1 focus:ring-blue-500"
                                spellCheck={false}
                            />
                        </div>
                    </div>
                ) : (
                    <div className="text-sm font-mono">
                        {diffLines.map((line, idx) => (
                            <div
                                key={idx}
                                className={`flex ${line.type === 'added'
                                        ? 'bg-green-900/30'
                                        : line.type === 'removed'
                                            ? 'bg-red-900/30'
                                            : ''
                                    }`}
                            >
                                {/* Line Numbers */}
                                <div className="w-10 text-right pr-2 text-gray-500 select-none border-r border-gray-700 bg-gray-800/50">
                                    {line.lineNumber.old || ''}
                                </div>
                                <div className="w-10 text-right pr-2 text-gray-500 select-none border-r border-gray-700 bg-gray-800/50">
                                    {line.lineNumber.new || ''}
                                </div>
                                {/* Indicator */}
                                <div className={`w-6 text-center select-none ${line.type === 'added'
                                        ? 'text-green-400'
                                        : line.type === 'removed'
                                            ? 'text-red-400'
                                            : 'text-gray-600'
                                    }`}>
                                    {line.type === 'added' ? '+' : line.type === 'removed' ? '-' : ' '}
                                </div>
                                {/* Content */}
                                <div className={`flex-1 px-2 whitespace-pre ${line.type === 'added'
                                        ? 'text-green-300'
                                        : line.type === 'removed'
                                            ? 'text-red-300'
                                            : 'text-gray-300'
                                    }`}>
                                    {line.content}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {/* Footer */}
            <div className="px-4 py-2 bg-gray-800 border-t border-gray-700 text-xs text-gray-400 flex items-center gap-2">
                <ArrowRight size={12} />
                <span>{diffLines.length} lines</span>
                {editable && (
                    <span className="ml-auto text-blue-400">
                        Editing enabled - changes will be applied on approval
                    </span>
                )}
            </div>
        </div>
    );
};
