import { type FC, useCallback, useState, useRef } from 'react';
import { Upload, X, FileText, CheckCircle2 } from 'lucide-react';

interface FileUploadZoneProps {
    onContentChange: (content: string, fileName?: string) => void;
    accept?: string[]; // e.g., ['.md', '.txt']
    maxSizeMB?: number;
    initialContent?: string;
}

export const FileUploadZone: FC<FileUploadZoneProps> = ({
    onContentChange,
    accept = ['.md', '.txt'],
    maxSizeMB = 5,
    initialContent = ''
}) => {
    const [isDragging, setIsDragging] = useState(false);
    const [fileStats, setFileStats] = useState<{ name: string; size: number } | null>(null);
    const [content, setContent] = useState(initialContent);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const [error, setError] = useState<string | null>(null);

    const handleFile = useCallback((file: File) => {
        setError(null);

        // Validate extension
        const ext = '.' + file.name.split('.').pop()?.toLowerCase();
        if (!accept.includes(ext)) {
            setError(`Unsupported file type. Allowed: ${accept.join(', ')}`);
            return;
        }

        // Validate size
        if (file.size > maxSizeMB * 1024 * 1024) {
            setError(`File too large. Max size: ${maxSizeMB}MB`);
            return;
        }

        const reader = new FileReader();
        reader.onload = (e) => {
            const text = e.target?.result as string;
            setContent(text);
            setFileStats({ name: file.name, size: file.size });
            onContentChange(text, file.name);
        };
        reader.onerror = () => setError("Failed to read file");
        reader.readAsText(file);
    }, [accept, maxSizeMB, onContentChange]);

    const handleDrop = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(false);
        if (e.dataTransfer.files?.[0]) {
            handleFile(e.dataTransfer.files[0]);
        }
    }, [handleFile]);

    const handlePaste = useCallback((e: React.ClipboardEvent) => {
        const text = e.clipboardData.getData('text');
        if (text) {
            setContent(text);
            setFileStats(null); // Clear file stats as it's a paste
            onContentChange(text, undefined);
        }
    }, [onContentChange]);

    const clearFile = () => {
        setContent('');
        setFileStats(null);
        setError(null);
        onContentChange('', undefined);
        if (fileInputRef.current) fileInputRef.current.value = '';
    };

    return (
        <div className="space-y-2">
            <div
                onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
                onDragLeave={() => setIsDragging(false)}
                onDrop={handleDrop}
                className={`relative border-2 border-dashed rounded-xl transition-all duration-200 ${isDragging
                        ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                        : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600 bg-gray-50/50 dark:bg-gray-800/50'
                    }`}
            >
                {!content ? (
                    <div className="p-8 text-center cursor-pointer" onClick={() => fileInputRef.current?.click()}>
                        <Upload className="w-10 h-10 mx-auto text-gray-400 mb-3" />
                        <p className="text-sm font-medium text-gray-700 dark:text-gray-200">
                            Click to upload or drag & drop
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                            {accept.join(', ')} (Max {maxSizeMB}MB)
                        </p>
                    </div>
                ) : (
                    <div className="p-4">
                        <div className="flex items-center justify-between mb-3">
                            <div className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-green-400">
                                {fileStats ? (
                                    <>
                                        <FileText size={16} />
                                        <span className="truncate max-w-[200px]">{fileStats.name}</span>
                                        <span className="text-gray-400 text-xs">({(fileStats.size / 1024).toFixed(1)} KB)</span>
                                    </>
                                ) : (
                                    <>
                                        <CheckCircle2 size={16} className="text-green-500" />
                                        <span>Content Loaded (Pasted)</span>
                                    </>
                                )}
                            </div>
                            <button
                                onClick={(e) => { e.stopPropagation(); clearFile(); }}
                                className="p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-full text-gray-500"
                            >
                                <X size={16} />
                            </button>
                        </div>
                        <div className="relative">
                            <textarea
                                value={content}
                                onChange={(e) => {
                                    setContent(e.target.value);
                                    onContentChange(e.target.value, fileStats?.name);
                                }}
                                onPaste={handlePaste}
                                className="w-full h-32 text-xs font-mono bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-3 resize-none focus:outline-none focus:ring-1 focus:ring-blue-500"
                                placeholder="Edit content here..."
                            />
                        </div>
                    </div>
                )}

                <input
                    type="file"
                    ref={fileInputRef}
                    className="hidden"
                    accept={accept.join(',')}
                    onChange={(e) => e.target.files?.[0] && handleFile(e.target.files[0])}
                />
            </div>
            {error && (
                <p className="text-xs text-red-500 flex items-center gap-1">
                    <X size={12} /> {error}
                </p>
            )}
        </div>
    );
};
