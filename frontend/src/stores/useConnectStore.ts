import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';
import type { WSCommand, WSMessage } from '../types/websocket';

interface ConnectState {
    socket: WebSocket | null;
    status: 'disconnected' | 'connecting' | 'connected' | 'reconnecting';
    lastError: string | null;
    reconnectAttempts: number;
    _lastMessage: WSMessage | null; // For router subscription

    // Actions
    connect: (url: string) => void;
    disconnect: () => void;
    send: <T>(command: WSCommand<T>) => void;

    // Internal
    _onMessage: (msg: WSMessage) => void;
    _scheduleReconnect: () => void;
    _startHeartbeat: () => void;
}

const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_DELAY = 3000;
const HEARTBEAT_INTERVAL = 30000;

export const useConnectStore = create<ConnectState>()(
    subscribeWithSelector((set, get) => {
        let heartbeatTimer: ReturnType<typeof setInterval> | null = null;
        let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
        let currentUrl: string | null = null;

        return {
            socket: null,
            status: 'disconnected',
            lastError: null,
            reconnectAttempts: 0,
            _lastMessage: null,

            connect: (url: string) => {
                // If already connected or connecting to the same URL, might skip, 
                // but let's check readyState to be safe.
                if (get().socket?.readyState === WebSocket.OPEN) return;

                currentUrl = url;
                set({ status: 'connecting' });
                const ws = new WebSocket(url);

                ws.onopen = () => {
                    set({ status: 'connected', reconnectAttempts: 0, lastError: null });
                    get()._startHeartbeat();
                };

                ws.onclose = (e) => {
                    set({ status: 'disconnected', socket: null });
                    if (heartbeatTimer) clearInterval(heartbeatTimer);

                    if (!e.wasClean) {
                        get()._scheduleReconnect();
                    }
                };

                ws.onerror = () => {
                    set({ lastError: 'WebSocket connection error' });
                };

                ws.onmessage = (event) => {
                    try {
                        const msg = JSON.parse(event.data) as WSMessage;
                        get()._onMessage(msg);
                    } catch (e) {
                        console.error('Failed to parse WS message:', e);
                    }
                };

                set({ socket: ws });
            },

            disconnect: () => {
                const { socket } = get();
                if (heartbeatTimer) clearInterval(heartbeatTimer);
                if (reconnectTimer) clearTimeout(reconnectTimer);

                socket?.close(1000, 'Client disconnect');
                set({ socket: null, status: 'disconnected', reconnectAttempts: 0 });
            },

            send: (command) => {
                const { socket, status } = get();
                if (socket && status === 'connected') {
                    socket.send(JSON.stringify(command));
                } else {
                    console.warn('Cannot send: not connected');
                }
            },

            _onMessage: (msg) => {
                set({ _lastMessage: msg });
            },

            _scheduleReconnect: () => {
                const { reconnectAttempts } = get();
                if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
                    set({ lastError: 'Max reconnection attempts reached' });
                    return;
                }

                set({ status: 'reconnecting', reconnectAttempts: reconnectAttempts + 1 });
                reconnectTimer = setTimeout(() => {
                    if (currentUrl) {
                        get().connect(currentUrl);
                    }
                }, RECONNECT_DELAY * (reconnectAttempts + 1));
            },

            _startHeartbeat: () => {
                if (heartbeatTimer) clearInterval(heartbeatTimer);
                heartbeatTimer = setInterval(() => {
                    const { socket } = get();
                    if (socket && socket.readyState === WebSocket.OPEN) {
                        socket.send(JSON.stringify({ cmd: 'ping' }));
                    }
                }, HEARTBEAT_INTERVAL);
            },
        };
    })
);
