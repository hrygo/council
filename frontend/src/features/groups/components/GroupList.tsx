import { useState } from 'react';
import { Plus, Users, Search } from 'lucide-react';
import { useGroups, useDeleteGroup, useCreateGroup, useUpdateGroup } from '../../../hooks/useGroups';
import { GroupCard } from './GroupCard';
import { CreateGroupModal } from './CreateGroupModal';
import type { Group, CreateGroupInput } from '../../../types/group';

export function GroupList() {
    const { data: groups, isLoading, error } = useGroups();
    const createGroup = useCreateGroup();
    const updateGroup = useUpdateGroup();
    const deleteGroup = useDeleteGroup();

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingGroup, setEditingGroup] = useState<Group | null>(null);
    const [searchQuery, setSearchQuery] = useState('');

    const handleCreate = (data: CreateGroupInput) => {
        createGroup.mutate(data, {
            onSuccess: () => setIsModalOpen(false)
        });
    };

    const handleUpdate = (data: CreateGroupInput) => {
        if (!editingGroup) return;
        updateGroup.mutate({ ...editingGroup, ...data }, {
            onSuccess: () => {
                setIsModalOpen(false);
                setEditingGroup(null);
            }
        });
    };

    const handleDelete = (group: Group) => {
        if (confirm(`Are you sure you want to delete group "${group.name}"?`)) {
            deleteGroup.mutate(group.id);
        }
    };

    const openCreateModal = () => {
        setEditingGroup(null);
        setIsModalOpen(true);
    };

    const openEditModal = (group: Group) => {
        setEditingGroup(group);
        setIsModalOpen(true);
    };

    const filteredGroups = groups?.filter(g =>
        g.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        g.system_prompt?.toLowerCase().includes(searchQuery.toLowerCase())
    );

    if (isLoading) {
        return (
            <div className="flex items-center justify-center p-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-4 bg-red-50 text-red-600 rounded-lg">
                Error loading groups: {error.message}
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between gap-4">
                <div className="relative flex-1 max-w-md">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                    <input
                        type="text"
                        placeholder="Search groups..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="w-full pl-10 pr-4 py-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50"
                    />
                </div>
                <button
                    onClick={openCreateModal}
                    className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
                >
                    <Plus size={18} />
                    New Group
                </button>
            </div>

            {filteredGroups?.length === 0 ? (
                <div className="text-center py-16 bg-white dark:bg-gray-800/50 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
                    <div className="mx-auto w-16 h-16 bg-gray-100 dark:bg-gray-800 rounded-full flex items-center justify-center mb-4 text-gray-400">
                        <Users size={32} />
                    </div>
                    <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100">No groups found</h3>
                    <p className="text-gray-500 mt-1 mb-6">Get started by creating your first collaboration group.</p>
                    <button
                        onClick={openCreateModal}
                        className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
                    >
                        Create Group
                    </button>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {filteredGroups?.map(group => (
                        <GroupCard
                            key={group.id}
                            group={group}
                            onEdit={openEditModal}
                            onDelete={handleDelete}
                        />
                    ))}
                </div>
            )}

            <CreateGroupModal
                open={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onSubmit={editingGroup ? handleUpdate : handleCreate}
                isLoading={createGroup.isPending || updateGroup.isPending}
                initialData={editingGroup}
            />
        </div>
    );
}
