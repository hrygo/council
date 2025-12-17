export interface ChatPanelProps {
    /**
     * 全屏模式
     */
    fullscreen?: boolean;

    /**
     * 退出全屏回调
     */
    onExitFullscreen?: () => void;

    /**
     * 只读模式 (禁用输入框)
     */
    readOnly?: boolean;

    /**
     * Session ID (用于发送用户消息)
     */
    sessionId?: string;
}
