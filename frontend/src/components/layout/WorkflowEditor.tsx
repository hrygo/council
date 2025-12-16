import type { FC } from 'react';
import { WorkflowCanvas } from '../modules/WorkflowCanvas';

export const WorkflowEditor: FC = () => {
    return (
        <div className="h-screen flex flex-col">
            <header className="h-14 border-b px-4 flex items-center justify-between bg-white dark:bg-gray-800">
                <h1 className="font-bold">Workflow Builder</h1>
                <div className="flex gap-2">
                    <button className="px-3 py-1 bg-blue-500 text-white rounded">Save</button>
                    <button className="px-3 py-1 bg-green-500 text-white rounded">Run Session</button>
                </div>
            </header>
            <div className="flex-1 overflow-hidden">
                <WorkflowCanvas />
            </div>
        </div>
    );
};
