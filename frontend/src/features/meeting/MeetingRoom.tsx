import { type FC, useRef, useEffect } from 'react';
import { PanelGroup, Panel, PanelResizeHandle, type ImperativePanelHandle } from 'react-resizable-panels';
import { Maximize2, Minimize2, ChevronLeft, ChevronRight } from 'lucide-react';
import { useLayoutStore } from '../../stores/useLayoutStore';
import { useWebSocketRouter } from '../../hooks/useWebSocketRouter';
import { useFullscreenShortcuts } from '../../hooks/useFullscreenShortcuts';

import WorkflowCanvas from '../../components/workflow/WorkflowCanvas';
import ChatPanel from '../../components/chat/ChatPanel';
import { DocumentReader } from '../../components/modules/DocumentReader';
import { HumanReviewModal } from '../execution/components/HumanReviewModal';
import { SessionStarter } from './SessionStarter';
import { RightPanel } from '../../components/panels/RightPanel';
import { useSessionStore } from '../../stores/useSessionStore';
import { useConnectStore } from '../../stores/useConnectStore';
import { useWorkflowRunStore } from '../../stores/useWorkflowRunStore';

// 面板全屏按钮
const PanelMaximizeButton: FC<{ panel: 'left' | 'center' | 'right' }> = ({ panel }) => {
    const { maximizedPanel, maximizePanel } = useLayoutStore();
    const isMaximized = maximizedPanel === panel;

    return (
        <button
            onClick={() => maximizePanel(isMaximized ? null : panel)}
            className="absolute top-2 right-2 z-10 p-1.5 bg-white/80 dark:bg-gray-800/80 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors shadow-sm backdrop-blur-sm text-gray-700 dark:text-gray-300"
            title={isMaximized ? "退出全屏" : "全屏聚焦"}
        >
            {isMaximized ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
        </button>
    );
};

// 侧边栏内的折叠触发器 (仅在展开时显示)
const SidebarCollapseTrigger: FC<{
    side: 'left' | 'right';
    onCollapse: () => void;
}> = ({ side, onCollapse }) => {
    const isLeft = side === 'left';
    return (
        <button
            onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onCollapse();
            }}
            className={`absolute top-1/2 -translate-y-1/2 z-20 p-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-full shadow-md hover:bg-gray-100 dark:hover:bg-gray-700 hover:scale-110 text-gray-500 dark:text-gray-400 transition-all cursor-pointer ${isLeft ? '-right-3' : '-left-3'
                }`}
            title="折叠面板"
            onMouseDown={(e) => e.stopPropagation()}
        >
            {isLeft ? <ChevronLeft size={14} /> : <ChevronRight size={14} />}
        </button>
    );
};

// 中间区域的展开触发器 (仅在折叠时显示)
const CenterExpandTrigger: FC<{
    side: 'left' | 'right';
    onExpand: () => void;
}> = ({ side, onExpand }) => {
    const isLeft = side === 'left';
    return (
        <button
            onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onExpand();
            }}
            className={`absolute top-1/2 -translate-y-1/2 z-20 py-4 px-0.5 bg-white dark:bg-gray-800 border-y border-r border-gray-200 dark:border-gray-700 shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 hover:w-6 transition-all group cursor-pointer flex items-center justify-center ${isLeft ? 'left-0 rounded-r-lg border-l-0' : 'right-0 rounded-l-lg border-l border-r-0'
                }`}
            title="展开面板"
        >
            <div className={`text-gray-400 hover:text-blue-500 transition-colors`}>
                {isLeft ? <ChevronRight size={14} /> : <ChevronLeft size={14} />}
            </div>
        </button>
    );
};

export const MeetingRoom: FC = () => {
    useWebSocketRouter();
    useFullscreenShortcuts();
    const {
        maximizedPanel,
        panelSizes,
        leftCollapsed,
        rightCollapsed,
        setPanelSizes,
        toggleLeftPanel,
        toggleRightPanel,
        maximizePanel
    } = useLayoutStore();

    // Refs for imperative panel control
    const leftPanelRef = useRef<ImperativePanelHandle>(null);
    const centerPanelRef = useRef<ImperativePanelHandle>(null);
    const rightPanelRef = useRef<ImperativePanelHandle>(null);

    const handleToggleLeft = () => {
        const panel = leftPanelRef.current;
        if (panel) {
            if (leftCollapsed) {
                panel.expand();
            } else {
                panel.collapse();
            }
        }
    };

    const handleToggleRight = () => {
        const panel = rightPanelRef.current;
        if (panel) {
            if (rightCollapsed) {
                panel.expand();
            } else {
                panel.collapse();
            }
        }
    };


    // Check for active session
    const currentSession = useSessionStore(state => state.currentSession);
    const wsStatus = useConnectStore(state => state.status);
    const wsConnect = useConnectStore(state => state.connect);
    const graphDefinition = useWorkflowRunStore(state => state.graphDefinition);

    // Auto-connect WebSocket if session exists but WS is disconnected
    useEffect(() => {
        if (currentSession && wsStatus === 'disconnected') {
            const wsHost = window.location.host;
            const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${wsProtocol}//${wsHost}/ws`;  // Backend WebSocket route is /ws
            wsConnect(wsUrl);
        }
    }, [currentSession, wsStatus, wsConnect]);

    // Fullscreen Mode
    if (maximizedPanel) {
        const onExit = () => maximizePanel(null);
        const panelMap = {
            left: <WorkflowCanvas fullscreen onExitFullscreen={onExit} workflowId={currentSession?.workflow_uuid} graph={graphDefinition} readOnly={true} />,
            center: <ChatPanel fullscreen onExitFullscreen={onExit} />,
            right: currentSession ? <RightPanel sessionId={currentSession.session_uuid} fullscreen onExitFullscreen={onExit} /> : <DocumentReader fullscreen onExitFullscreen={onExit} />,
        };
        return (
            <div className="h-screen w-screen fixed top-0 left-0 bg-white dark:bg-gray-900 z-50">
                {panelMap[maximizedPanel]}
                <HumanReviewModal />
            </div>
        );
    }

    const isRunning = !!currentSession; // Session exists implies running/active

    return (
        <div className="h-screen w-full bg-gray-100 dark:bg-gray-900 overflow-hidden relative">

            {/* Session Starter Overlay */}
            {!currentSession && (
                <SessionStarter onStarted={() => { }} />
            )}

            <PanelGroup autoSaveId="council-layout-persistence" direction="horizontal" onLayout={setPanelSizes}>
                {/* Left Panel: Workflow Canvas */}
                <Panel
                    ref={leftPanelRef}
                    defaultSize={panelSizes[0]}
                    minSize={15}
                    order={1}
                    collapsible
                    collapsedSize={0}
                    onCollapse={() => { if (!leftCollapsed) toggleLeftPanel() }}
                    onExpand={() => { if (leftCollapsed) toggleLeftPanel() }}
                    className="flex flex-col relative transition-none"
                >
                    <div className="relative h-full w-full group">
                        <PanelMaximizeButton panel="left" />
                        {!leftCollapsed && <SidebarCollapseTrigger side="left" onCollapse={handleToggleLeft} />}
                        <WorkflowCanvas
                            readOnly={isRunning}
                            workflowId={currentSession?.workflow_uuid}
                            graph={graphDefinition}
                        />
                    </div>
                </Panel>

                <PanelResizeHandle className="w-1.5 bg-gray-200 hover:bg-blue-400 dark:bg-gray-800 dark:hover:bg-blue-600 transition-colors cursor-col-resize z-10" />

                {/* Center Panel: Chat Stream */}
                <Panel
                    ref={centerPanelRef}
                    defaultSize={panelSizes[1]}
                    minSize={25}
                    order={2}
                >
                    <div className="relative h-full w-full">
                        <PanelMaximizeButton panel="center" />

                        {/* Expand Triggers for collapsed side panels */}
                        {leftCollapsed && <CenterExpandTrigger side="left" onExpand={handleToggleLeft} />}
                        {rightCollapsed && <CenterExpandTrigger side="right" onExpand={handleToggleRight} />}

                        <ChatPanel />
                    </div>
                </Panel>

                <PanelResizeHandle className="w-1.5 bg-gray-200 hover:bg-blue-400 dark:bg-gray-800 dark:hover:bg-blue-600 transition-colors cursor-col-resize z-10" />

                {/* Right Panel: Knowledge Panel */}
                <Panel
                    ref={rightPanelRef}
                    defaultSize={panelSizes[2]}
                    minSize={15}
                    order={3}
                    collapsible
                    collapsedSize={0}
                    onCollapse={() => { if (!rightCollapsed) toggleRightPanel() }}
                    onExpand={() => { if (rightCollapsed) toggleRightPanel() }}
                >
                    <div className="relative h-full w-full">
                        <PanelMaximizeButton panel="right" />
                        {!rightCollapsed && <SidebarCollapseTrigger side="right" onCollapse={handleToggleRight} />}
                        {currentSession ? (
                            <RightPanel sessionId={currentSession.session_uuid} />
                        ) : (
                            <div className="h-full flex items-center justify-center text-gray-500 dark:text-gray-400">
                                <p>启动会话后查看相关知识</p>
                            </div>
                        )}
                    </div>
                </Panel>
            </PanelGroup>

            <HumanReviewModal />
        </div>
    );
};
