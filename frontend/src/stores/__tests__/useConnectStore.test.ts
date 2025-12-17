import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { act } from '@testing-library/react';
import { useConnectStore } from '../useConnectStore';
import type { WSMessage } from '../../types/websocket';

// Mock WebSocket
class MockWebSocket {
    static CONNECTING = 0;
    static OPEN = 1;
    static CLOSING = 2;
    static CLOSED = 3;

    url: string;
    onopen: () => void = () => { };
    onclose: () => void = () => { };
    onmessage: (event: { data: string }) => void = () => { };
    onerror: () => void = () => { };
    send: (data: string) => void = vi.fn();
    close: () => void = vi.fn();
    readyState: number = MockWebSocket.CONNECTING;

    constructor(url: string) {
        this.url = url;
        // Auto-open in next tick
        setTimeout(() => {
            this.readyState = MockWebSocket.OPEN;
            if (this.onopen) this.onopen();
        }, 10);
    }
}

describe('useConnectStore', () => {
    beforeEach(() => {
        vi.stubGlobal('WebSocket', MockWebSocket);
        vi.useFakeTimers();
        useConnectStore.setState({
            socket: null,
            status: 'disconnected',
            reconnectAttempts: 0,
            _lastMessage: null,
        });
    });

    afterEach(() => {
        vi.useRealTimers();
        useConnectStore.getState().disconnect();
        vi.unstubAllGlobals();
    });

    it('should connect to websocket', async () => {
        const { connect } = useConnectStore.getState();

        act(() => {
            connect('ws://test');
        });

        expect(useConnectStore.getState().status).toBe('connecting');

        await act(async () => {
            await vi.advanceTimersByTimeAsync(100);
        });

        expect(useConnectStore.getState().status).toBe('connected');
        expect(useConnectStore.getState().socket).not.toBeNull();
    });

    it('should handle incoming messages', async () => {
        const { connect } = useConnectStore.getState();

        act(() => {
            connect('ws://test');
        });

        await act(async () => {
            await vi.advanceTimersByTimeAsync(100);
        });

        const socket = useConnectStore.getState().socket as unknown as MockWebSocket;
        expect(socket).toBeDefined();

        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const msg: WSMessage = { event: 'token_stream', data: { chunk: 'hello' } as any };

        act(() => {
            socket.onmessage({ data: JSON.stringify(msg) });
        });

        expect(useConnectStore.getState()._lastMessage).toEqual(msg);
    });
});
