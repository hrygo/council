import { useEffect } from 'react';
import { useLayoutStore } from '../stores/useLayoutStore';

/**
 * useFullscreenShortcuts
 * 
 * Provides global keyboard shortcuts for layout management.
 * - 'f': Toggle maximize for Center panel (Chat) by default, or cycle?
 *   Actually, PRD might specify logic. Let's assume 'f' maximizes Chat (center) if nothing maximized,
 *   or exits if something is maximized? 
 *   Better: 'Escape' always exits fullscreen.
 *   'Shift + f': Maximize Chat? 
 *   For now: 'Escape' to exit. 
 *   And let's map: 
 *   - Ctrl+1 / Cmd+1 -> Maximize Left (Workflow)
 *   - Ctrl+2 / Cmd+2 -> Maximize Center (Chat)
 *   - Ctrl+3 / Cmd+3 -> Maximize Right (Docs)
 */
export const useFullscreenShortcuts = () => {
    const { maximizePanel, maximizedPanel } = useLayoutStore();

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            // Ignore if input is active
            if (['INPUT', 'TEXTAREA'].includes((e.target as HTMLElement).tagName)) {
                return;
            }

            if (e.key === 'Escape') {
                if (maximizedPanel) {
                    maximizePanel(null);
                }
            }

            // Shortcuts with Meta (Cmd) or Ctrl
            if (e.metaKey || e.ctrlKey) {
                switch (e.key) {
                    case '1':
                        e.preventDefault();
                        maximizePanel(maximizedPanel === 'left' ? null : 'left');
                        break;
                    case '2':
                        e.preventDefault();
                        maximizePanel(maximizedPanel === 'center' ? null : 'center');
                        break;
                    case '3':
                        e.preventDefault();
                        maximizePanel(maximizedPanel === 'right' ? null : 'right');
                        break;
                }
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [maximizedPanel, maximizePanel]);
};
