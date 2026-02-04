import React from 'react';
import { RefreshCw } from 'lucide-react';
import { Tab, Status } from '../../types';

interface HeaderProps {
    activeTab: Tab;
    status: Status | null;
    onRefresh: () => void;
}

export const Header: React.FC<HeaderProps> = ({ activeTab, status, onRefresh }) => {
    const getTitle = () => {
        switch (activeTab) {
            case 'dashboard': return 'W-G Gateway Dashboard';
            case 'nodes': return 'Infrastructure Nodes';
            case 'services': return 'Routing Policies';
            case 'logs': return 'System Event Log';
            case 'settings': return 'Cluster Configuration';
            default: return 'Management Console';
        }
    };

    return (
        <header className="flex flex-col md:flex-row justify-between items-start md:items-end gap-6 mb-12">
            <div>
                <h1 className="text-4xl font-bold tracking-tight text-white mb-2">{getTitle()}</h1>
                <p className="text-slate-400 font-medium capitalize">Manage your {activeTab} across the cluster</p>
            </div>
            <div className="flex items-center gap-4">
                <div className="bg-slate-900/50 px-4 py-2 rounded-full border border-white/5 flex items-center gap-3">
                    <div className={`w-2 h-2 rounded-full ${status?.ready ? 'bg-emerald-500 shadow-[0_0_10px_#10b981]' : 'bg-amber-500'} animate-pulse`} />
                    <span className="text-sm font-semibold">{status?.ready ? 'Gateway Online' : 'Pending DNS/SSL'}</span>
                </div>
                <button
                    onClick={onRefresh}
                    className="p-2.5 rounded-xl bg-white/5 hover:bg-white/10 transition-colors border border-white/5 group"
                >
                    <RefreshCw size={20} className="text-slate-300 group-active:rotate-180 transition-transform duration-500" />
                </button>
            </div>
        </header>
    );
};
