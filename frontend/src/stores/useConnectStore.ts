import { create } from 'zustand';

interface ConnectState {
    socket: WebSocket | null;
    isConnected: boolean;
    connect: (url: string) => void;
    disconnect: () => void;
    lastMessage: any;
    sendMessage: (msg: any) => void;
}

export const useConnectStore = create<ConnectState>((set, get) => ({
    socket: null,
    isConnected: false,
    lastMessage: null,

    connect: (url: string) => {
        if (get().socket) return;
        console.log('Connecting to WS:', url);
        const ws = new WebSocket(url);

        ws.onopen = () => {
            console.log('WS Connected');
            set({ isConnected: true });
        };
        ws.onclose = () => {
            console.log('WS Disconnected');
            set({ isConnected: false, socket: null });
        };
        ws.onerror = (err) => {
            console.error('WS Error', err);
        };
        ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                set({ lastMessage: data });
            } catch (e) {
                console.error('WS Parse Error', e);
            }
        };

        set({ socket: ws });
    },

    disconnect: () => {
        get().socket?.close();
        set({ socket: null, isConnected: false });
    },

    sendMessage: (msg: any) => {
        const { socket, isConnected } = get();
        if (socket && isConnected) {
            socket.send(JSON.stringify(msg));
        } else {
            console.warn('Cannot send message: Socket not connected');
        }
    }
}));
