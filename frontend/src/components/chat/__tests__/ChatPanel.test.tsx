// Imports cleaned up
import { render, screen } from '@testing-library/react';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ChatPanel } from '../ChatPanel';
import { useSessionStore } from '../../../stores/useSessionStore';

// Mock scrollIntoView
window.HTMLElement.prototype.scrollIntoView = vi.fn();

describe('ChatPanel', () => {
    beforeEach(() => {
        useSessionStore.getState().clearSession();
    });

    it('should render empty state when no messages', () => {
        render(<ChatPanel sessionId="test-session" />);
        expect(screen.getByText('等待会议开始...')).toBeInTheDocument();
    });

    it('should render message groups with headers', () => {
        const store = useSessionStore.getState();
        store.initSession({
            sessionId: 'test-session',
            workflowId: 'test-workflow',
            groupId: 'test-group',
            nodes: [] // Providing empty nodes array
        });

        // Manually mocking state update since we don't have the full engine running
        useSessionStore.setState({
            messageGroups: [
                {
                    nodeId: 'node-1',
                    nodeName: 'Analyst',
                    nodeType: 'agent',
                    status: 'completed',
                    isParallel: false,
                    messages: [
                        {
                            id: 'msg-1',
                            nodeId: 'node-1', // Add missing nodeId
                            role: 'agent',
                            agentName: 'Analyst',
                            content: 'Test message',
                            timestamp: new Date(),
                            isStreaming: false
                        }
                    ]
                }
            ]
        });

        render(<ChatPanel sessionId="test-session" />);
        // "Analyst" appears in Header and MessageAvatar name
        const analystElements = screen.getAllByText('Analyst');
        expect(analystElements.length).toBeGreaterThan(0);

        expect(screen.getByText('Test message')).toBeInTheDocument();
    });

    it('should show active status for current node group', () => {
        // const store = useSessionStore.getState();
        useSessionStore.setState({
            currentSession: {
                id: 'sess-1',
                workflowId: 'wf-1',
                groupId: 'grp-1',
                status: 'running',
                nodes: new Map(),
                activeNodeIds: ['node-1'],
                totalTokens: 0,
                totalCostUsd: 0,
                startedAt: new Date()
            },
            messageGroups: [
                {
                    nodeId: 'node-1',
                    nodeName: 'Analyst',
                    nodeType: 'agent',
                    status: 'running',
                    isParallel: false,
                    messages: []
                }
            ]
        });

        render(<ChatPanel sessionId="test-session" />);
        // We need to check if the active class is applied.
        // Since we composed components, looking for a specific indicator or class on parent.
        // Detailed verification usually via snapshot or specific class check.
        const groupStatus = screen.getByText('进行中');
        expect(groupStatus).toBeInTheDocument();
    });
});
