import { GroupList } from '../components/GroupList';

export function GroupsPage() {
    return (
        <div className="h-full overflow-y-auto bg-gray-50 dark:bg-gray-950 p-6 md:p-8">
            <div className="max-w-7xl mx-auto space-y-8">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 tracking-tight">
                        Groups
                    </h1>
                    <p className="text-gray-500 dark:text-gray-400 mt-2 text-lg">
                        Manage your collaboration spaces and default agent assignments.
                    </p>
                </div>

                <GroupList />
            </div>
        </div>
    );
}
