import { render, screen, fireEvent } from '@testing-library/react';
import { PropertyPanel } from './PropertyPanel';
import type { WorkflowNode } from '../../../../types/workflow';
import { vi } from 'vitest';

describe('PropertyPanel', () => {
    const mockOnUpdate = vi.fn();
    const mockOnDelete = vi.fn();
    const mockOnClose = vi.fn();

    const baseNode: WorkflowNode = {
        id: '1',
        type: 'vote',
        position: { x: 0, y: 0 },
        data: {
            label: 'Test Vote Node',
            threshold: 0.6,
            vote_type: 'yes_no',
            agent_ids: []
        }
    };

    it('renders correctly for Vote Node', () => {
        render(
            <PropertyPanel
                node={baseNode}
                onUpdate={mockOnUpdate}
                onDelete={mockOnDelete}
                onClose={mockOnClose}
            />
        );

        expect(screen.getByText('Test Vote Node')).toBeDefined();
        expect(screen.getByText('Approval Threshold')).toBeDefined();
        // Check for specific form elements
        expect(screen.getByText('Yes/No')).toBeDefined();
    });

    it('calls onUpdate when label changes', () => {
        render(
            <PropertyPanel
                node={baseNode}
                onUpdate={mockOnUpdate}
                onDelete={mockOnDelete}
                onClose={mockOnClose}
            />
        );

        const input = screen.getByDisplayValue('Test Vote Node');
        fireEvent.change(input, { target: { value: 'New Name' } });
        expect(mockOnUpdate).toHaveBeenCalledWith('1', { label: 'New Name' });
    });

    it('calls onDelete when delete button clicked', () => {
        render(
            <PropertyPanel
                node={baseNode}
                onUpdate={mockOnUpdate}
                onDelete={mockOnDelete}
                onClose={mockOnClose}
            />
        );

        fireEvent.click(screen.getByText('Delete Node'));
        expect(mockOnDelete).toHaveBeenCalledWith('1');
    });

    it('renders Loop Node form correctly', () => {
        const loopNode: WorkflowNode = {
            ...baseNode,
            type: 'loop',
            data: {
                label: 'Loop Node',
                max_rounds: 5,
                exit_condition: 'max_rounds',
                agent_pairs: []
            }
        };

        render(
            <PropertyPanel
                node={loopNode}
                onUpdate={mockOnUpdate}
                onDelete={mockOnDelete}
                onClose={mockOnClose}
            />
        );

        expect(screen.getByText('Max Rounds')).toBeDefined();
        expect(screen.getByDisplayValue('5')).toBeDefined();
    });
});
