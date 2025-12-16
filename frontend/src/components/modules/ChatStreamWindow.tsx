import type { FC } from 'react';
import { CostEstimator } from './CostEstimator';
import { ParallelMessageRow } from './ParallelMessageRow';

export const ChatStreamWindow: FC<{ fullscreen?: boolean; onExitFullscreen?: () => void }> = ({ fullscreen, onExitFullscreen }) => {
    return (
        <div className={`h-full w-full bg-white dark:bg-gray-800 flex flex-col ${fullscreen ? 'p-8' : 'p-4'}`}>
            <div className="flex justify-between items-center mb-4">
                <h2 className="font-bold text-lg">Council Chamber</h2>
                {fullscreen && <button onClick={onExitFullscreen} className="text-sm bg-gray-200 px-2 py-1 rounded">Exit Fullscreen</button>}
            </div>

            {/* Session Header / Cost Estimator */}
            <CostEstimator />

            <div className="flex-1 overflow-y-auto space-y-4 p-2">
                <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                    <p className="text-sm font-semibold">System</p>
                    <p>Session started. Waiting for inputs...</p>
                </div>

                {/* Placeholder: Parallel Execution Visual */}
                <div className="my-4">
                    <p className="text-xs text-gray-400 mb-1 text-center">— Parallel Execution —</p>
                    <ParallelMessageRow messages={[
                        { agentName: "Architect", content: "**Plan**: We should use a Microservices pattern." },
                        { agentName: "Security", content: "**Risk**: Ensure mTLS between services." }
                    ]} />
                </div>

                {/* Placeholder for messages */}
            </div>
            <div className="mt-4 border-t pt-2">
                <input className="w-full p-2 border rounded bg-gray-50 dark:bg-gray-700" placeholder="Type your instruction..." />
            </div>
        </div>
    );
};
