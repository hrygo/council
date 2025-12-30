import { type FC, useState, useEffect, useMemo } from 'react';
import { FolderOpen, File, ChevronRight, ChevronDown, History, Code2 } from 'lucide-react';

interface FileEntity {
    file_uuid: string;
    session_uuid: string;
    path: string;
    version: number;
    content: string;
    author: string;
    reason: string;
    created_at: string;
}

interface VFSExplorerProps {
    sessionId: string;
}

interface TreeNode {
    name: string;
    path: string;
    isDir: boolean;
    children: TreeNode[];
    file?: FileEntity;
}

// Build tree structure from flat file list
function buildTree(files: FileEntity[]): TreeNode[] {
    const root: TreeNode[] = [];
    const pathMap = new Map<string, TreeNode>();

    // Get unique paths with latest version
    const latestFiles = new Map<string, FileEntity>();
    files.forEach(f => {
        const existing = latestFiles.get(f.path);
        if (!existing || f.version > existing.version) {
            latestFiles.set(f.path, f);
        }
    });

    latestFiles.forEach((file, path) => {
        const parts = path.split('/').filter(Boolean);
        let currentPath = '';
        let currentLevel = root;

        parts.forEach((part, idx) => {
            currentPath = currentPath ? `${currentPath}/${part}` : part;
            const isLast = idx === parts.length - 1;

            let node = pathMap.get(currentPath);
            if (!node) {
                node = {
                    name: part,
                    path: currentPath,
                    isDir: !isLast,
                    children: [],
                    file: isLast ? file : undefined,
                };
                pathMap.set(currentPath, node);
                currentLevel.push(node);
            }
            currentLevel = node.children;
        });
    });

    return root;
}

// Tree Node Component
const TreeItem: FC<{
    node: TreeNode;
    selectedPath: string | null;
    onSelect: (node: TreeNode) => void;
    level: number;
}> = ({ node, selectedPath, onSelect, level }) => {
    const [expanded, setExpanded] = useState(level < 2);
    const isSelected = selectedPath === node.path;

    return (
        <div>
            <div
                className={`flex items-center gap-1.5 py-1 px-2 cursor-pointer rounded text-sm transition-colors ${isSelected
                        ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300'
                        : 'hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-700 dark:text-gray-300'
                    }`}
                style={{ paddingLeft: `${level * 12 + 8}px` }}
                onClick={() => {
                    if (node.isDir) {
                        setExpanded(!expanded);
                    } else {
                        onSelect(node);
                    }
                }}
            >
                {node.isDir ? (
                    <>
                        {expanded ? (
                            <ChevronDown size={14} className="text-gray-400" />
                        ) : (
                            <ChevronRight size={14} className="text-gray-400" />
                        )}
                        <FolderOpen size={14} className="text-yellow-500" />
                    </>
                ) : (
                    <>
                        <span className="w-3.5" />
                        <File size={14} className="text-gray-400" />
                    </>
                )}
                <span className="truncate">{node.name}</span>
                {node.file && (
                    <span className="ml-auto text-xs text-gray-400">v{node.file.version}</span>
                )}
            </div>
            {node.isDir && expanded && (
                <div>
                    {node.children.map((child) => (
                        <TreeItem
                            key={child.path}
                            node={child}
                            selectedPath={selectedPath}
                            onSelect={onSelect}
                            level={level + 1}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

// File Content Viewer
const FileViewer: FC<{
    file: FileEntity;
    history: FileEntity[];
    onVersionSelect: (version: number) => void;
}> = ({ file, history, onVersionSelect }) => {
    return (
        <div className="h-full flex flex-col">
            {/* Header */}
            <div className="flex items-center justify-between px-4 py-2 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
                <div className="flex items-center gap-2">
                    <Code2 size={16} className="text-blue-500" />
                    <span className="font-medium text-sm text-gray-900 dark:text-gray-100">
                        {file.path}
                    </span>
                </div>
                <div className="flex items-center gap-2">
                    <History size={14} className="text-gray-400" />
                    <select
                        value={file.version}
                        onChange={(e) => onVersionSelect(Number(e.target.value))}
                        className="text-xs bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded px-2 py-1"
                    >
                        {history.map((h) => (
                            <option key={h.version} value={h.version}>
                                v{h.version} - {h.author} ({new Date(h.created_at).toLocaleTimeString()})
                            </option>
                        ))}
                    </select>
                </div>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-auto bg-gray-900 p-4">
                <pre className="text-sm text-gray-100 font-mono whitespace-pre-wrap">
                    {file.content}
                </pre>
            </div>

            {/* Footer */}
            <div className="px-4 py-2 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 text-xs text-gray-500">
                <span>Modified by <strong>{file.author}</strong></span>
                {file.reason && <span className="ml-2">• {file.reason}</span>}
            </div>
        </div>
    );
};

export const VFSExplorer: FC<VFSExplorerProps> = ({ sessionId }) => {
    const [files, setFiles] = useState<FileEntity[]>([]);
    const [selectedPath, setSelectedPath] = useState<string | null>(null);
    const [selectedVersion, setSelectedVersion] = useState<number | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Fetch files
    useEffect(() => {
        if (!sessionId) return;

        const fetchFiles = async () => {
            setLoading(true);
            try {
                const res = await fetch(`/api/v1/sessions/${sessionId}/files`);
                if (!res.ok) throw new Error('Failed to fetch files');
                const data = await res.json();
                setFiles(data.files || []);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Unknown error');
            } finally {
                setLoading(false);
            }
        };

        fetchFiles();
    }, [sessionId]);

    // Build tree
    const tree = useMemo(() => buildTree(files), [files]);

    // Get selected file and its history
    const selectedFile = useMemo(() => {
        if (!selectedPath) return null;
        const fileVersions = files.filter(f => f.path === selectedPath).sort((a, b) => b.version - a.version);
        if (fileVersions.length === 0) return null;
        const version = selectedVersion ?? fileVersions[0].version;
        return fileVersions.find(f => f.version === version) || fileVersions[0];
    }, [files, selectedPath, selectedVersion]);

    const fileHistory = useMemo(() => {
        if (!selectedPath) return [];
        return files.filter(f => f.path === selectedPath).sort((a, b) => b.version - a.version);
    }, [files, selectedPath]);

    const handleNodeSelect = (node: TreeNode) => {
        setSelectedPath(node.path);
        setSelectedVersion(null); // Reset to latest
    };

    if (loading) {
        return (
            <div className="h-full flex items-center justify-center text-gray-500">
                <div className="animate-spin w-6 h-6 border-2 border-blue-500 border-t-transparent rounded-full" />
            </div>
        );
    }

    if (error) {
        return (
            <div className="h-full flex items-center justify-center text-red-500">
                <p>{error}</p>
            </div>
        );
    }

    if (files.length === 0) {
        return (
            <div className="h-full flex flex-col items-center justify-center text-gray-500 dark:text-gray-400 gap-2">
                <FolderOpen size={32} className="text-gray-300 dark:text-gray-600" />
                <p className="text-sm">No files in VFS yet</p>
                <p className="text-xs">Files will appear here after the Surgeon modifies them.</p>
            </div>
        );
    }

    return (
        <div className="h-full flex">
            {/* File Tree */}
            <div className="w-48 border-r border-gray-200 dark:border-gray-700 overflow-y-auto bg-white dark:bg-gray-900">
                <div className="p-2 border-b border-gray-200 dark:border-gray-700">
                    <h3 className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase">
                        Files ({files.length})
                    </h3>
                </div>
                <div className="py-1">
                    {tree.map((node) => (
                        <TreeItem
                            key={node.path}
                            node={node}
                            selectedPath={selectedPath}
                            onSelect={handleNodeSelect}
                            level={0}
                        />
                    ))}
                </div>
            </div>

            {/* File Content */}
            <div className="flex-1 bg-white dark:bg-gray-900">
                {selectedFile ? (
                    <FileViewer
                        file={selectedFile}
                        history={fileHistory}
                        onVersionSelect={setSelectedVersion}
                    />
                ) : (
                    <div className="h-full flex items-center justify-center text-gray-400">
                        <p className="text-sm">Select a file to view its contents</p>
                    </div>
                )}
            </div>
        </div>
    );
};
