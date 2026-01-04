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
        getState: () => ({ connect: mockConnect, status: 'connected' }),
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

    it('should handle full 3-step wizard flow and start session', async () => {
        // 1. Mock Templates Data
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useTemplates as any).mockReturnValue({
            data: [{
                template_uuid: 't1',
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
        (useSessionStore as any).mockImplementation((selector: any) => {
            if (selector && selector.name === 'initSession') return mockInitSession;
            return mockInitSession;
        });
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useSessionStore as any).mockReturnValue(mockInitSession);

        // 3. Mock API Response
        globalFetch.mockResolvedValue({
            ok: true,
            json: async () => ({ session_id: 'sess_123', status: 'started' })
        });

        // 4. Render
        render(<SessionStarter onStarted={mockOnStarted} />);

        // === Step 1: Template Selection ===
        expect(screen.getByText('Select Workflow Template')).toBeInTheDocument();
        fireEvent.click(screen.getByText('Council Debate'));

        const nextBtn1 = screen.getByText('Next Step');
        await waitFor(() => expect(nextBtn1).not.toBeDisabled());
        fireEvent.click(nextBtn1);

        // === Step 2: Input ===
        await waitFor(() => expect(screen.getByText('Upload Document')).toBeInTheDocument());

        // Enter Objective
        const objectiveInput = screen.getByPlaceholderText(/Optimize for clarity/i);
        fireEvent.change(objectiveInput, { target: { value: 'Fix grammar' } });

        const nextBtn2 = screen.getByText('Next Step');
        fireEvent.click(nextBtn2);

        // === Step 3: Confirmation ===
        await waitFor(() => expect(screen.getByText('Session Summary')).toBeInTheDocument());
        expect(screen.getByText('Fix grammar')).toBeInTheDocument();

        // Click Start
        const startBtn = screen.getByText('Start Council Session');
        fireEvent.click(startBtn);

        // === Verification ===
        await waitFor(() => {
            // Check API payload
            expect(globalFetch).toHaveBeenCalledWith('/api/v1/workflows/execute', expect.objectContaining({
                body: expect.stringContaining('"optimization_objective":"Fix grammar"')
            }));

            // Check API payload structure for document_content (empty in this case)
            expect(globalFetch).toHaveBeenCalledWith('/api/v1/workflows/execute', expect.objectContaining({
                body: expect.stringContaining('"document_content":""')
            }));

            // Check Navigation
            expect(mockNavigate).toHaveBeenCalledWith('/meeting');

            // Check WebSocket connection
            expect(mockConnect).toHaveBeenCalledWith(expect.stringMatching(/ws:.*\/ws/));
        });
    });

    it('should navigate back and forth between steps', async () => {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useTemplates as any).mockReturnValue({
            data: [{ template_uuid: 't1', name: 'Test', is_system: true, graph: {} }],
            isLoading: false
        });

        render(<SessionStarter onStarted={mockOnStarted} />);

        // Select template and go to Input
        fireEvent.click(screen.getByText('Test'));
        fireEvent.click(screen.getByText('Next Step'));

        await waitFor(() => expect(screen.getByText('Upload Document')).toBeInTheDocument());

        // Go Back to Template
        fireEvent.click(screen.getByText('Back'));
        await waitFor(() => expect(screen.getByText('Select Workflow Template')).toBeInTheDocument());

        // Go Forward again
        fireEvent.click(screen.getByText('Next Step')); // Template is still selected in state
        await waitFor(() => expect(screen.getByText('Upload Document')).toBeInTheDocument());

        // Go Forward to Confirm
        fireEvent.click(screen.getByText('Next Step'));
        await waitFor(() => expect(screen.getByText('Session Summary')).toBeInTheDocument());

        // Go Back from Confirm
        fireEvent.click(screen.getByText('Back'));
        await waitFor(() => expect(screen.getByText('Upload Document')).toBeInTheDocument());
    });
});
