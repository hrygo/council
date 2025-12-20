import { type FC } from 'react';
import { useNavigate } from 'react-router-dom';
import { Play, Boxes, Users, Network, ArrowRight, Zap, Shield, Cpu } from 'lucide-react';

export const HomePage: FC = () => {
    const navigate = useNavigate();

    return (
        <div className="h-full w-full bg-gray-50 dark:bg-gray-900 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-blue-100 via-gray-50 to-white dark:from-slate-800 dark:via-gray-900 dark:to-black text-gray-900 dark:text-gray-100 overflow-y-auto pb-20">

            {/* Hero Section */}
            <div className="relative pt-20 pb-8 px-6 lg:px-8 max-w-7xl mx-auto flex flex-col items-center text-center">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[600px] h-[400px] bg-blue-500/20 blur-[120px] rounded-full pointer-events-none" />

                <div className="relative z-10 space-y-6">
                    <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 text-sm font-medium border border-blue-100 dark:border-blue-800 animate-in fade-in slide-in-from-bottom-4 duration-700">
                        <span className="relative flex h-2 w-2">
                            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                            <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
                        </span>
                        Council AI v0.13.0 Active
                    </div>

                    <h1 className="text-5xl md:text-7xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-gray-900 via-blue-800 to-gray-900 dark:from-white dark:via-blue-200 dark:to-gray-400 animate-in fade-in slide-in-from-bottom-6 duration-700 delay-100">
                        Orchestrate Intelligence <br />
                        <span className="text-blue-600 dark:text-blue-500">Decide with Confidence</span>
                    </h1>

                    <p className="max-w-2xl mx-auto text-lg text-gray-600 dark:text-gray-400 animate-in fade-in slide-in-from-bottom-8 duration-700 delay-200">
                        The advanced multi-agent consensus system. Assemble your council of AI experts to debate, analyze, and execute complex workflows with human-in-the-loop precision.
                    </p>

                    <div className="flex items-center justify-center gap-4 pt-4 animate-in fade-in slide-in-from-bottom-10 duration-700 delay-300">
                        <button
                            onClick={() => navigate('/meeting')}
                            className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-xl font-semibold text-lg shadow-lg shadow-blue-500/20 flex items-center gap-2 transition-all hover:scale-105 active:scale-95"
                        >
                            <Play className="fill-current" size={20} />
                            Start Session
                        </button>
                        <button
                            onClick={() => navigate('/editor')}
                            className="bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-900 dark:text-gray-200 border border-gray-200 dark:border-gray-700 px-8 py-3 rounded-xl font-semibold text-lg shadow-sm flex items-center gap-2 transition-all"
                        >
                            <Boxes size={20} />
                            Builder
                        </button>
                    </div>
                </div>
            </div>

            {/* Feature Grid */}
            <div className="py-6 px-6 lg:px-8 max-w-7xl mx-auto">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    <FeatureCard
                        icon={<Users size={24} />}
                        title="Agent Groups"
                        desc="Curate specialized teams of agents. Define personas, assign prompts, and organize your council structure."
                        action={() => navigate('/groups')}
                        color="text-purple-500"
                        bg="bg-purple-50 dark:bg-purple-900/20"
                    />
                    <FeatureCard
                        icon={<Boxes size={24} />}
                        title="Workflow Engine"
                        desc="Visual node-based editor. Design complex logic loops, voting mechanisms, and human-in-the-loop checkpoints."
                        action={() => navigate('/editor')}
                        color="text-blue-500"
                        bg="bg-blue-50 dark:bg-blue-900/20"
                        delay={100}
                    />
                    <FeatureCard
                        icon={<Network size={24} />}
                        title="Model Hub"
                        desc="Connect to any LLM. Manage models, context windows, and API configurations centrally."
                        action={() => navigate('/agents')}
                        color="text-emerald-500"
                        bg="bg-emerald-50 dark:bg-emerald-900/20"
                        delay={200}
                    />
                </div>
            </div>

            {/* Stats / Info Section */}
            <div className="py-8 border-t border-gray-200 dark:border-gray-800 bg-white/50 dark:bg-gray-800/20 px-6">
                <div className="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
                    <StatItem icon={<Zap />} label="Real-time Debate" value="< 50ms" sub="Latency" />
                    <StatItem icon={<Shield />} label="Consensus Engine" value="RAFT" sub="Protocol" />
                    <StatItem icon={<Cpu />} label="Agent Capacity" value="Unlimited" sub="Scalability" />
                </div>
            </div>

            {/* Footer */}
            <footer className="py-8 text-center text-gray-400 text-sm">
                <p>© 2025 Council AI. All systems nominal.</p>
            </footer>
        </div>
    );
};

const FeatureCard: FC<{
    icon: React.ReactNode;
    title: string;
    desc: string;
    action: () => void;
    color: string;
    bg: string;
    delay?: number;
}> = ({ icon, title, desc, action, color, bg, delay = 0 }) => (
    <div
        onClick={action}
        className="group relative bg-white dark:bg-gray-800 p-6 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-xl hover:border-blue-500/30 transition-all cursor-pointer overflow-hidden animate-in fade-in slide-in-from-bottom-8 duration-700 fill-mode-backwards"
        style={{ animationDelay: `${delay}ms` }}
    >
        <div className={`absolute top-0 right-0 p-3 opacity-10 group-hover:opacity-20 transition-opacity ${color}`}>
            <div className="scale-150 transform">{icon}</div>
        </div>

        <div className={`w-12 h-12 rounded-xl flex items-center justify-center mb-4 ${bg} ${color}`}>
            {icon}
        </div>

        <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2 flex items-center gap-2 group-hover:text-blue-500 transition-colors">
            {title}
            <ArrowRight size={16} className="opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
        </h3>
        <p className="text-gray-500 dark:text-gray-400 leading-relaxed text-sm">
            {desc}
        </p>
    </div>
);

const StatItem: FC<{ icon: React.ReactNode; label: string; value: string; sub: string }> = ({ icon, label, value, sub }) => (
    <div className="flex flex-col items-center gap-2">
        <div className="text-gray-400 mb-2">{icon}</div>
        <div className="text-3xl font-bold text-gray-900 dark:text-gray-100">{value}</div>
        <div className="text-sm text-gray-500 font-medium uppercase tracking-wider">{label}</div>
        <div className="text-xs text-gray-400">{sub}</div>
    </div>
);
