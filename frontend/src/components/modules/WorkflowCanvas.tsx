import type { FC } from 'react';

export const WorkflowCanvas: FC<{ readOnly?: boolean; fullscreen?: boolean; onExitFullscreen?: () => void }> = ({ readOnly, fullscreen, onExitFullscreen }) => {
    return (
        <div className={`h-full w-full bg-slate-50 dark:bg-slate-900 border-r border-gray-200 dark:border-gray-800 flex flex-col ${fullscreen ? 'p-8' : 'p-4'}`}>
            <div className="flex justify-between items-center mb-4">
                <h2 className="font-bold text-lg">Workflow Canvas</h2>
                {fullscreen && <button onClick={onExitFullscreen} className="text-sm bg-gray-200 px-2 py-1 rounded">Exit Fullscreen</button>}
            </div>
            <div className="flex-1 border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center text-gray-500">
                {readOnly ? "Canvas (Read-Only Mode)" : "Canvas Editor"}
            </div>
        </div>
    );
};
