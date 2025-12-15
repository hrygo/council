import { create } from 'zustand';

interface SessionState {
    user: {
        id: string;
        name: string;
    } | null;
    setUser: (user: { id: string; name: string }) => void;
    logout: () => void;
}

export const useSessionStore = create<SessionState>((set) => ({
    user: null,
    setUser: (user) => set({ user }),
    logout: () => set({ user: null }),
}));
