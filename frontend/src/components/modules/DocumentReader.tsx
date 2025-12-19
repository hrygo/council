import type { FC } from 'react';

export const DocumentReader: FC<{ fullscreen?: boolean; onExitFullscreen?: () => void }> = ({ fullscreen, onExitFullscreen }) => {
    return (
        <div className={`h-full w-full bg-gray-50 dark:bg-gray-900 border-l border-gray-200 dark:border-gray-700 flex flex-col ${fullscreen ? 'p-8' : 'p-4'}`}>
            <div className="flex justify-between items-center mb-4">
                <h2 className="font-bold text-lg text-gray-800 dark:text-gray-200">Knowledge Base</h2>
                {fullscreen && <button onClick={onExitFullscreen} className="text-sm bg-gray-200 dark:bg-gray-700 px-2 py-1 rounded">Exit Fullscreen</button>}
            </div>
            <div className="flex-1 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg flex items-center justify-center text-gray-500 dark:text-gray-400">
                No documents loaded.
            </div>
        </div>
    );
};
