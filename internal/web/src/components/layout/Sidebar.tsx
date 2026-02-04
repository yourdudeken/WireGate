import React from 'react';
import { LayoutDashboard, Server, Globe, Terminal, Settings } from 'lucide-react';
import { Tab } from '../../types';

interface SidebarProps {
    activeTab: Tab;
    setActiveTab: (tab: Tab) => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ activeTab, setActiveTab }) => {
    return (
        <aside className="w-20 bg-slate-950/80 backdrop-blur-2xl border-r border-white/5 flex flex-col items-center py-8 sticky top-0 h-screen z-50">
            <div
                className="w-12 h-12 flex items-center justify-center cursor-pointer transition-transform hover:scale-110 active:scale-95 mb-12"
                onClick={() => setActiveTab('dashboard')}
            >
                <img src="/logo.svg" alt="WG Gateway" className="w-10 h-10 drop-shadow-[0_0_10px_rgba(59,130,246,0.5)]" />
            </div>
            <nav className="flex flex-col gap-6 flex-1">
                <NavItem active={activeTab === 'dashboard'} onClick={() => setActiveTab('dashboard')} icon={<LayoutDashboard size={24} />} title="Dashboard" />
                <NavItem active={activeTab === 'nodes'} onClick={() => setActiveTab('nodes')} icon={<Server size={24} />} title="Infrastructure Nodes" />
                <NavItem active={activeTab === 'services'} onClick={() => setActiveTab('services')} icon={<Globe size={24} />} title="Traffic Routes" />
                <NavItem active={activeTab === 'logs'} onClick={() => setActiveTab('logs')} icon={<Terminal size={24} />} title="System Events" />
                <div className="flex-1" />
                <NavItem active={activeTab === 'settings'} onClick={() => setActiveTab('settings')} icon={<Settings size={24} />} title="Cluster Settings" />
            </nav>
        </aside>
    );
};

interface NavItemProps {
    icon: React.ReactNode;
    active?: boolean;
    onClick: () => void;
    title: string;
}

const NavItem: React.FC<NavItemProps> = ({ icon, active = false, onClick, title }) => (
    <div className="relative group">
        <button
            onClick={onClick}
            className={`w-12 h-12 flex items-center justify-center rounded-2xl transition-all duration-300 ${active ? 'bg-blue-600/10 text-blue-500 shadow-[inset_0_0_0_1px_rgba(37,99,235,0.2)]' : 'text-slate-500 hover:text-white hover:bg-white/5'
                }`}
        >
            {icon}
        </button>
        <div className="absolute left-[calc(100%+10px)] top-1/2 -translate-y-1/2 px-3 py-1.5 bg-slate-800 text-white text-xs font-bold rounded-lg opacity-0 pointer-events-none group-hover:opacity-100 transition-opacity whitespace-nowrap z-[100] shadow-xl border border-white/5">
            {title}
            <div className="absolute top-1/2 -left-1 -translate-y-1/2 w-2 h-2 bg-slate-800 rotate-45 border-l border-b border-white/5"></div>
        </div>
    </div>
);
