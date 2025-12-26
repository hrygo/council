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
        expect(screen.getByText('meeting.waitingForStart')).toBeInTheDocument();
    });

    it('should render message groups with headers', () => {
        const store = useSessionStore.getState();
        store.initSession({
            session_uuid: 'test-session',
            workflow_id: 'test-workflow',
            group_uuid: 'test-group',
            nodes: [] // Providing empty nodes array
        });

        // Manually mocking state update since we don't have the full engine running
        useSessionStore.setState({
            messageGroups: [
                {
                    node_id: 'node-1',
                    nodeName: 'Analyst',
                    nodeType: 'agent',
                    status: 'completed',
                    isParallel: false,
                    messages: [
                        {
                            message_uuid: 'msg-1',
                            node_id: 'node-1', // Add missing nodeId
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

    it('should render parallel messages', () => {
        useSessionStore.setState({
            messageGroups: [
                {
                    node_id: 'node-parallel-1',
                    nodeName: 'ParallelReview',
                    nodeType: 'parallel',
                    status: 'running',
                    isParallel: true,
                    messages: [
                        {
                            message_uuid: 'msg-p1',
                            node_id: 'node-parallel-1',
                            role: 'agent',
                            agentName: 'Security',
                            content: 'Security Check',
                            timestamp: new Date(),
                            isStreaming: false
                        },
                        {
                            message_uuid: 'msg-p2',
                            node_id: 'node-parallel-1',
                            role: 'agent',
                            agentName: 'Performance',
                            content: 'Performance Check',
                            timestamp: new Date(),
                            isStreaming: false
                        }
                    ]
                }
            ]
        });

        render(<ChatPanel sessionId="test-session" />);

        expect(screen.getByText('Security')).toBeInTheDocument();
        expect(screen.getByText('Performance')).toBeInTheDocument();
        expect(screen.getByText('Security Check')).toBeInTheDocument();
    });

    it('should show active status for current node group', () => {
        // const store = useSessionStore.getState();
        useSessionStore.setState({
            currentSession: {
                session_uuid: 'sess-1',
                workflow_id: 'wf-1',
                group_uuid: 'grp-1',
                status: 'running',
                nodes: new Map(),
                active_node_ids: ['node-1'],
                totalTokens: 0,
                totalCostUsd: 0,
                startedAt: new Date()
            },
            messageGroups: [
                {
                    node_id: 'node-1',
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
        const groupStatus = screen.getByText('meeting.status.processing');
        expect(groupStatus).toBeInTheDocument();
    });
});
