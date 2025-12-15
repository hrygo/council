import { create } from 'zustand';

interface LayoutState {
    sidebarOpen: boolean;
    rightPanelOpen: boolean;
    toggleSidebar: () => void;
    toggleRightPanel: () => void;
}

export const useLayoutStore = create<LayoutState>((set) => ({
    sidebarOpen: true,
    rightPanelOpen: true,
    toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
    toggleRightPanel: () => set((state) => ({ rightPanelOpen: !state.rightPanelOpen })),
}));
