import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { SessionStarter } from '../SessionStarter';
import { useTemplates } from '../../../hooks/useTemplates';
import { useSessionStore } from '../../../stores/useSessionStore';

// === Global Mocks ===
const mockNavigate = vi.fn();
const mockConnect = vi.fn();

vi.mock('react-router-dom', () => ({
    useNavigate: () => mockNavigate, // Always return our spy
}));

vi.mock('../../../hooks/useTemplates', () => ({
    useTemplates: vi.fn(),
}));

vi.mock('../../../stores/useSessionStore', () => ({
    useSessionStore: vi.fn(),
}));

vi.mock('../../../stores/useConnectStore', () => ({
    useConnectStore: {
        getState: () => ({ connect: mockConnect }),
    },
}));

// Mock fetch
const globalFetch = vi.fn();
vi.stubGlobal('fetch', globalFetch);


describe('SessionStarter', () => {
    const mockInitSession = vi.fn();
    const mockOnStarted = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should handle template with missing nodes gracefully', async () => {
        // 1. Mock Templates Data
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useTemplates as any).mockReturnValue({
            data: [{
                id: 't1',
                name: 'Test Template',
                description: 'A broken template',
                is_system: true,
                graph: {} // MISSING NODES - This caused the crash
            }],
            isLoading: false
        });

        // 2. Mock Session Store
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useSessionStore as any).mockImplementation((selector: any) => {

            if (selector && selector.name === 'initSession') return mockInitSession;
            return mockInitSession;
        });

        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useSessionStore as any).mockReturnValue(mockInitSession);

        // 3. Mock API Response
        globalFetch.mockResolvedValue({
            ok: true,
            json: async () => ({ session_id: 'sess_123' })
        });

        // 4. Render
        render(<SessionStarter onStarted={mockOnStarted} />);

        // 5. Select Template (Click on the item)
        // By default the first template might be selected or we click it
        fireEvent.click(screen.getByText('Test Template'));

        // 6. Enter Topic
        // The error shows the placeholder is "e.g. Should artificial intelligence..."
        // but we can use getByRole('textbox') for simplicity since there is only one textarea
        const textarea = screen.getByRole('textbox');
        fireEvent.change(textarea, { target: { value: 'Test Topic' } });

        // 7. Click Start
        const startBtn = screen.getByText('Start Council Session');
        // Wait for button to be enabled (it is disabled if topic is empty)
        await waitFor(() => expect(startBtn).not.toBeDisabled());
        fireEvent.click(startBtn);

        // 8. Verification
        await waitFor(() => {
            // Check if API was called
            expect(globalFetch).toHaveBeenCalledWith('/api/v1/workflows/execute', expect.anything());

            // Check if navigate was called (Success!)
            expect(mockNavigate).toHaveBeenCalledWith('/meeting');

            // Check if initSession was called with empty nodes (Sanity check for our fix)
            expect(mockInitSession).toHaveBeenCalledWith(expect.objectContaining({
                sessionId: 'sess_123',
                nodes: [] // Should be empty array, not crash
            }));
        });
    });

    it('should connect to WebSocket after successful API call', async () => {
        // 1. Mock Templates Data with valid graph
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useTemplates as any).mockReturnValue({
            data: [{
                id: 't1',
                name: 'Council Debate',
                description: 'Test template',
                is_system: true,
                graph: {
                    nodes: {
                        start: { id: 'start', type: 'start', name: 'Start' }
                    }
                }
            }],
            isLoading: false
        });

        // 2. Mock Session Store
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useSessionStore as any).mockReturnValue(mockInitSession);

        // 3. Mock API Response
        globalFetch.mockResolvedValue({
            ok: true,
            json: async () => ({ session_id: 'sess_456', status: 'started' })
        });

        // 4. Render
        render(<SessionStarter onStarted={mockOnStarted} />);

        // 5. Select Template
        fireEvent.click(screen.getByText('Council Debate'));

        // 6. Enter Topic
        const textarea = screen.getByRole('textbox');
        fireEvent.change(textarea, { target: { value: 'Should AI be regulated?' } });

        // 7. Click Start
        const startBtn = screen.getByText('Start Council Session');
        await waitFor(() => expect(startBtn).not.toBeDisabled());
        fireEvent.click(startBtn);

        // 8. Verification - WebSocket connect should be called
        await waitFor(() => {
            expect(mockConnect).toHaveBeenCalledWith(expect.stringMatching(/ws:.*\/ws/));
        });
    });
});
