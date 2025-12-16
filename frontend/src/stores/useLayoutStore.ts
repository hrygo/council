import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

interface LayoutState {
    panelSizes: number[]; // [left, center, right]
    leftCollapsed: boolean;
    rightCollapsed: boolean;
    maximizedPanel: 'left' | 'center' | 'right' | null;

    setPanelSizes: (sizes: number[]) => void;
    toggleLeftPanel: () => void;
    toggleRightPanel: () => void;
    maximizePanel: (panel: 'left' | 'center' | 'right' | null) => void;
}

export const useLayoutStore = create<LayoutState>()(
    persist(
        (set) => ({
            panelSizes: [20, 50, 30],
            leftCollapsed: false,
            rightCollapsed: false,
            maximizedPanel: null,

            setPanelSizes: (sizes) => set({ panelSizes: sizes }),
            toggleLeftPanel: () => set((state) => ({ leftCollapsed: !state.leftCollapsed })),
            toggleRightPanel: () => set((state) => ({ rightCollapsed: !state.rightCollapsed })),
            maximizePanel: (panel) => set({ maximizedPanel: panel }),
        }),
        {
            name: 'council-layout',
            storage: createJSONStorage(() => localStorage),
        }
    )
);
