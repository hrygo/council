import { describe, it, expect, beforeEach } from 'vitest';
import { useLayoutStore } from '../useLayoutStore';

describe('useLayoutStore', () => {
    beforeEach(() => {
        // Reset state properties only, not the actions
        const state = useLayoutStore.getState();
        state.setPanelSizes([20, 50, 30]);
        useLayoutStore.setState({
            leftCollapsed: false,
            rightCollapsed: false,
            maximizedPanel: null,
        });
    });

    it('should have correct initial state', () => {
        const state = useLayoutStore.getState();
        expect(state.panelSizes).toEqual([20, 50, 30]);
        expect(state.leftCollapsed).toBe(false);
        expect(state.rightCollapsed).toBe(false);
        expect(state.maximizedPanel).toBeNull();
    });

    it('should set panel sizes', () => {
        useLayoutStore.getState().setPanelSizes([30, 40, 30]);
        expect(useLayoutStore.getState().panelSizes).toEqual([30, 40, 30]);
    });

    it('should toggle left panel', () => {
        useLayoutStore.getState().toggleLeftPanel();
        expect(useLayoutStore.getState().leftCollapsed).toBe(true);

        useLayoutStore.getState().toggleLeftPanel();
        expect(useLayoutStore.getState().leftCollapsed).toBe(false);
    });

    it('should toggle right panel', () => {
        useLayoutStore.getState().toggleRightPanel();
        expect(useLayoutStore.getState().rightCollapsed).toBe(true);

        useLayoutStore.getState().toggleRightPanel();
        expect(useLayoutStore.getState().rightCollapsed).toBe(false);
    });

    it('should maximize panel', () => {
        useLayoutStore.getState().maximizePanel('center');
        expect(useLayoutStore.getState().maximizedPanel).toBe('center');

        useLayoutStore.getState().maximizePanel(null);
        expect(useLayoutStore.getState().maximizedPanel).toBeNull();
    });
});
