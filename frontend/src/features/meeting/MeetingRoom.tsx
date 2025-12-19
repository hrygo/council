import type { FC } from 'react';
import { PanelGroup, Panel, PanelResizeHandle } from 'react-resizable-panels';
import { Maximize2, Minimize2, ChevronLeft, ChevronRight, RotateCcw } from 'lucide-react';
import { useLayoutStore } from '../../stores/useLayoutStore';
import { useWebSocketRouter } from '../../hooks/useWebSocketRouter';
import { useFullscreenShortcuts } from '../../hooks/useFullscreenShortcuts';

import WorkflowCanvas from '../../components/workflow/WorkflowCanvas';
import ChatPanel from '../../components/chat/ChatPanel';
import { DocumentReader } from '../../components/modules/DocumentReader';
import { HumanReviewModal } from '../execution/components/HumanReviewModal';

// 面板全屏按钮
const PanelMaximizeButton: FC<{ panel: 'left' | 'center' | 'right' }> = ({ panel }) => {
    const { maximizedPanel, maximizePanel } = useLayoutStore();
    const isMaximized = maximizedPanel === panel;

    return (
        <button
            onClick={() => maximizePanel(isMaximized ? null : panel)}
            className="absolute top-2 right-2 z-10 p-1.5 bg-white/80 dark:bg-gray-800/80 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors shadow-sm backdrop-blur-sm"
            title={isMaximized ? "退出全屏" : "全屏聚焦"}
        >
            {isMaximized ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
        </button>
    );
};

// 折叠按钮
const CollapseButton: FC<{
    side: 'left' | 'right';
    collapsed: boolean;
    onToggle: () => void;
}> = ({ side, collapsed, onToggle }) => {
    const isLeft = side === 'left';
    const Icon = isLeft
        ? (collapsed ? ChevronRight : ChevronLeft)
        : (collapsed ? ChevronLeft : ChevronRight);

    return (
        <button
            onClick={onToggle}
            className={`absolute top-1/2 -translate-y-1/2 z-20 p-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-full shadow-md hover:bg-gray-100 dark:hover:bg-gray-700 transition-all ${isLeft ? '-right-3' : '-left-3'
                }`}
            title={collapsed ? '展开面板' : '折叠面板'}
        >
            <Icon size={14} className="text-gray-500" />
        </button>
    );
};

// 重置布局按钮
const ResetLayoutButton: FC = () => {
    const resetLayout = useLayoutStore((state) => state.resetLayout);

    return (
        <button
            onClick={resetLayout}
            className="absolute bottom-4 left-1/2 -translate-x-1/2 z-30 px-3 py-1.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-full shadow-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-all flex items-center gap-1.5 text-xs text-gray-600 dark:text-gray-400"
            title="重置为默认布局"
        >
            <RotateCcw size={12} />
            重置布局
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
                <HumanReviewModal />
            </div>
        );
    }

    const isRunning = true; // Placeholder for session status

    return (
        <div className="h-screen w-full bg-gray-100 dark:bg-gray-900 overflow-hidden relative">
            <PanelGroup direction="horizontal" onLayout={setPanelSizes}>
                {/* Left Panel: Workflow Canvas */}
                <Panel
                    defaultSize={panelSizes[0]}
                    minSize={leftCollapsed ? 0 : 15}
                    order={1}
                    collapsible
                    collapsedSize={0}
                    onCollapse={() => { if (!leftCollapsed) toggleLeftPanel() }}
                    onExpand={() => { if (leftCollapsed) toggleLeftPanel() }}
                >
                    <div className="relative h-full w-full group">
                        <PanelMaximizeButton panel="left" />
                        <CollapseButton side="left" collapsed={leftCollapsed} onToggle={toggleLeftPanel} />
                        <WorkflowCanvas readOnly={isRunning} workflowId="1eb04085-f215-430b-9279-880c98f99e3a" />
                    </div>
                </Panel>

                <PanelResizeHandle className="w-1.5 bg-gray-200 hover:bg-blue-400 dark:bg-gray-700 dark:hover:bg-blue-600 transition-colors cursor-col-resize" />

                {/* Center Panel: Chat Stream */}
                <Panel defaultSize={panelSizes[1]} minSize={25} order={2}>
                    <div className="relative h-full w-full">
                        <PanelMaximizeButton panel="center" />
                        <ChatPanel />
                    </div>
                </Panel>

                <PanelResizeHandle className="w-1.5 bg-gray-200 hover:bg-blue-400 dark:bg-gray-700 dark:hover:bg-blue-600 transition-colors cursor-col-resize" />

                {/* Right Panel: Document Reader */}
                <Panel
                    defaultSize={panelSizes[2]}
                    minSize={rightCollapsed ? 0 : 15}
                    order={3}
                    collapsible
                    collapsedSize={0}
                    onCollapse={() => { if (!rightCollapsed) toggleRightPanel() }}
                    onExpand={() => { if (rightCollapsed) toggleRightPanel() }}
                >
                    <div className="relative h-full w-full">
                        <PanelMaximizeButton panel="right" />
                        <CollapseButton side="right" collapsed={rightCollapsed} onToggle={toggleRightPanel} />
                        <DocumentReader />
                    </div>
                </Panel>
            </PanelGroup>

            {/* Reset Layout Button */}
            <ResetLayoutButton />

            <HumanReviewModal />
        </div>
    );
};
