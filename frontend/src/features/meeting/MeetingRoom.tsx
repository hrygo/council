import { type FC, useRef, useState } from 'react';
import { PanelGroup, Panel, PanelResizeHandle, type ImperativePanelHandle } from 'react-resizable-panels';
import { Maximize2, Minimize2, ChevronLeft, ChevronRight, RotateCcw, Check } from 'lucide-react';
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
            className="absolute top-2 right-2 z-10 p-1.5 bg-white/80 dark:bg-gray-800/80 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors shadow-sm backdrop-blur-sm text-gray-700 dark:text-gray-300"
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
    // Logic: 
    // Left Panel:  Collapsed(True) -> Show Right Arrow (to expand)
    //              Expanded(False) -> Show Left Arrow (to collapse)
    const Icon = isLeft
        ? (collapsed ? ChevronRight : ChevronLeft)
        : (collapsed ? ChevronLeft : ChevronRight);

    return (
        <button
            onClick={(e) => {
                e.preventDefault(); // Prevent accidental selection
                e.stopPropagation(); // Prevent panel focus steal
                onToggle();
            }}
            className={`absolute top-1/2 -translate-y-1/2 z-20 p-1.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-full shadow-md hover:bg-gray-100 dark:hover:bg-gray-700 hover:scale-110 text-gray-500 dark:text-gray-400 transition-all cursor-pointer ${isLeft ? '-right-3' : '-left-3'
                }`}
            title={collapsed ? '展开面板' : '折叠面板'}
            onMouseDown={(e) => e.stopPropagation()} // Prevent resize handle drag start if overlapping
        >
            <Icon size={14} />
        </button>
    );
};

// 重置布局按钮
const ResetLayoutButton: FC<{ onReset: () => void }> = ({ onReset }) => {
    const [feedback, setFeedback] = useState(false);

    const handleClick = () => {
        onReset();
        setFeedback(true);
        setTimeout(() => setFeedback(false), 1500);
    };

    return (
        <button
            onClick={handleClick}
            disabled={feedback}
            className={`absolute bottom-4 left-1/2 -translate-x-1/2 z-30 px-3 py-1.5 
            ${feedback ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 border-green-200 dark:border-green-800' : 'bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-400 border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700'} 
            border rounded-full shadow-lg transition-all flex items-center gap-1.5 text-xs font-medium`}
            title="重置为默认布局"
        >
            {feedback ? <Check size={12} /> : <RotateCcw size={12} />}
            {feedback ? '已重置' : '重置布局'}
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
        maximizePanel,
        resetLayout
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

    const handleReset = () => {
        resetLayout();
        // Imperatively reset sizes
        // react-resizable-panels API: resize(size: number)
        if (leftPanelRef.current) {
            leftPanelRef.current.resize(20);
            leftPanelRef.current.expand(); // Make sure it is expanded
        }
        if (centerPanelRef.current) centerPanelRef.current.resize(50);
        if (rightPanelRef.current) {
            rightPanelRef.current.resize(30);
            rightPanelRef.current.expand(); // Make sure it is expanded
        }
    };

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
                        <CollapseButton side="left" collapsed={leftCollapsed} onToggle={handleToggleLeft} />
                        <WorkflowCanvas readOnly={isRunning} workflowId="1eb04085-f215-430b-9279-880c98f99e3a" />
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
                        <ChatPanel />
                    </div>
                </Panel>

                <PanelResizeHandle className="w-1.5 bg-gray-200 hover:bg-blue-400 dark:bg-gray-800 dark:hover:bg-blue-600 transition-colors cursor-col-resize z-10" />

                {/* Right Panel: Document Reader */}
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
                        <CollapseButton side="right" collapsed={rightCollapsed} onToggle={handleToggleRight} />
                        <DocumentReader />
                    </div>
                </Panel>
            </PanelGroup>

            {/* Reset Layout Button with Callback */}
            <ResetLayoutButton onReset={handleReset} />

            <HumanReviewModal />
        </div>
    );
};
