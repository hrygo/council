import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ChatInput } from '../ChatInput';
import { useConnectStore } from '../../../stores/useConnectStore';

// Mock useConnectStore
vi.mock('../../../stores/useConnectStore', () => ({
    useConnectStore: vi.fn(),
}));

// Mock react-i18next
vi.mock('react-i18next', () => ({
    useTranslation: () => ({
        t: (key: string) => key,
    }),
}));

describe('ChatInput', () => {
    it('should update input value on change', () => {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useConnectStore as any).mockReturnValue(vi.fn());
        render(<ChatInput sessionId="s1" />);

        const textarea = screen.getByPlaceholderText('meeting.inputPlaceholder');
        fireEvent.change(textarea, { target: { value: 'Hello' } });
        expect((textarea as HTMLTextAreaElement).value).toBe('Hello');
    });

    it('should call send and clear input on send button click', () => {
        const sendMock = vi.fn();
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useConnectStore as any).mockImplementation((selector: any) => selector({ send: sendMock }));

        render(<ChatInput sessionId="s1" />);
        const textarea = screen.getByPlaceholderText('meeting.inputPlaceholder');
        const button = screen.getByRole('button');

        fireEvent.change(textarea, { target: { value: 'Hello World' } });
        fireEvent.click(button);

        expect(sendMock).toHaveBeenCalledWith(expect.objectContaining({
            cmd: 'user_input',
            data: { content: 'Hello World', session_id: 's1' }
        }));
        expect((textarea as HTMLTextAreaElement).value).toBe('');
    });

    it('should call send on Enter key press', () => {
        const sendMock = vi.fn();
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useConnectStore as any).mockImplementation((selector: any) => selector({ send: sendMock }));

        render(<ChatInput sessionId="s1" />);
        const textarea = screen.getByPlaceholderText('meeting.inputPlaceholder');

        fireEvent.change(textarea, { target: { value: 'Enter test' } });
        fireEvent.keyDown(textarea, { key: 'Enter', code: 'Enter', shiftKey: false });

        expect(sendMock).toHaveBeenCalled();
    });

    it('should not call send on Shift + Enter', () => {
        const sendMock = vi.fn();
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (useConnectStore as any).mockImplementation((selector: any) => selector({ send: sendMock }));

        render(<ChatInput sessionId="s1" />);
        const textarea = screen.getByPlaceholderText('meeting.inputPlaceholder');

        fireEvent.change(textarea, { target: { value: 'Shift Enter test' } });
        fireEvent.keyDown(textarea, { key: 'Enter', code: 'Enter', shiftKey: true });

        expect(sendMock).not.toHaveBeenCalled();
    });
});
