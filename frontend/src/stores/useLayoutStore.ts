import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

const DEFAULT_PANEL_SIZES = [20, 50, 30];

interface LayoutState {
    panelSizes: number[]; // [left, center, right]
    leftCollapsed: boolean;
    rightCollapsed: boolean;
    maximizedPanel: 'left' | 'center' | 'right' | null;

    setPanelSizes: (sizes: number[]) => void;
    toggleLeftPanel: () => void;
    toggleRightPanel: () => void;
    maximizePanel: (panel: 'left' | 'center' | 'right' | null) => void;
    resetLayout: () => void;
}

export const useLayoutStore = create<LayoutState>()(
    persist(
        (set) => ({
            panelSizes: DEFAULT_PANEL_SIZES,
            leftCollapsed: false,
            rightCollapsed: false,
            maximizedPanel: null,

            setPanelSizes: (sizes) => set({ panelSizes: sizes }),
            toggleLeftPanel: () => set((state) => ({ leftCollapsed: !state.leftCollapsed })),
            toggleRightPanel: () => set((state) => ({ rightCollapsed: !state.rightCollapsed })),
            maximizePanel: (panel) => set({ maximizedPanel: panel }),
            resetLayout: () => set({
                panelSizes: DEFAULT_PANEL_SIZES,
                leftCollapsed: false,
                rightCollapsed: false,
                maximizedPanel: null,
            }),
        }),
        {
            name: 'council-layout',
            storage: createJSONStorage(() => localStorage),
        }
    )
);
