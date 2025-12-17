import type { FC } from 'react';
import { PanelGroup, Panel, PanelResizeHandle } from 'react-resizable-panels';
import { Maximize2, Minimize2 } from 'lucide-react';
import { useLayoutStore } from '../../stores/useLayoutStore';
import { useWebSocketRouter } from '../../hooks/useWebSocketRouter';

import WorkflowCanvas from '../../components/workflow/WorkflowCanvas';
import ChatPanel from '../../components/chat/ChatPanel';
import { DocumentReader } from '../../components/modules/DocumentReader';

const PanelMaximizeButton: FC<{ panel: 'left' | 'center' | 'right' }> = ({ panel }) => {
    const { maximizedPanel, maximizePanel } = useLayoutStore();
    const isMaximized = maximizedPanel === panel;

    return (
        <button
            onClick={() => maximizePanel(isMaximized ? null : panel)}
            className="absolute top-2 right-2 z-10 p-1.5 bg-white/80 dark:bg-gray-800/80 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors shadow-sm backdrop-blur-sm"
            title={isMaximized ? "Exit Fullscreen" : "Fullscreen Focus"}
        >
            {isMaximized ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
        </button>
    );
};

export const MeetingRoom: FC = () => {
    useWebSocketRouter();
    const { maximizedPanel, panelSizes, leftCollapsed, rightCollapsed, setPanelSizes, toggleLeftPanel, toggleRightPanel, maximizePanel } = useLayoutStore();

    // Fullscreen Mode
    if (maximizedPanel) {
        const onExit = () => maximizePanel(null);
        const panelMap = {
            left: <WorkflowCanvas fullscreen onExitFullscreen={onExit} />,
            center: <ChatPanel fullscreen onExitFullscreen={onExit} />,
            right: <DocumentReader fullscreen onExitFullscreen={onExit} />,
        };
        return (
            <div className="h-screen w-screen fixed top-0 left-0 bg-white dark:bg-gray-900 z-50">
                {panelMap[maximizedPanel]}
            </div>
        );
    }

    const isRunning = true; // Placeholder for session status

    return (
        <div className="h-screen w-full bg-gray-100 dark:bg-gray-900 overflow-hidden">
            <PanelGroup direction="horizontal" onLayout={setPanelSizes}>
                {/* Left Panel: Workflow Canvas */}
                <Panel
                    defaultSize={panelSizes[0]}
                    minSize={10}
                    order={1}
                    collapsible
                    collapsedSize={0}
                    onCollapse={() => { if (!leftCollapsed) toggleLeftPanel() }}
                    onExpand={() => { if (leftCollapsed) toggleLeftPanel() }}
                >
                    <div className="relative h-full w-full group">
                        <PanelMaximizeButton panel="left" />
                        {/* Temporary: Hardcoded Workflow ID for demo. In real app, this comes from URL or Context */}
                        <WorkflowCanvas readOnly={isRunning} workflowId="1eb04085-f215-430b-9279-880c98f99e3a" />
                        {/* TODO: Add a floating 'Run' button if needed, but 'MeetingRoom' usually implies running session. */}
                    </div>
                </Panel>

                <PanelResizeHandle className="w-1 bg-gray-200 hover:bg-blue-400 dark:bg-gray-800 dark:hover:bg-blue-600 transition-colors" />

                {/* Center Panel: Chat Stream */}
                <Panel defaultSize={panelSizes[1]} minSize={30} order={2}>
                    <div className="relative h-full w-full">
                        <PanelMaximizeButton panel="center" />
                        <ChatPanel />
                    </div>
                </Panel>

                <PanelResizeHandle className="w-1 bg-gray-200 hover:bg-blue-400 dark:bg-gray-800 dark:hover:bg-blue-600 transition-colors" />

                {/* Right Panel: Document Reader */}
                <Panel
                    defaultSize={panelSizes[2]}
                    minSize={10}
                    order={3}
                    collapsible
                    collapsedSize={0}
                    onCollapse={() => { if (!rightCollapsed) toggleRightPanel() }}
                    onExpand={() => { if (rightCollapsed) toggleRightPanel() }}
                >
                    <div className="relative h-full w-full">
                        <PanelMaximizeButton panel="right" />
                        <DocumentReader />
                    </div>
                </Panel>
            </PanelGroup>
        </div>
    );
};
