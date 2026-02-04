import React from 'react';
import { motion } from 'framer-motion';
import { Users, Layers, Activity, ShieldCheck, Globe, ChevronRight, Shield, Settings, BarChart3, ArrowDown, ArrowUp } from 'lucide-react';
import { Status, Peer, Service } from '../../types';
import { StatCard } from '../../components/ui/Cards';
import { UsageChart } from './UsageChart';
import { formatBytes } from '../../utils';

interface DashboardViewProps {
    status: Status | null;
    peers: Peer[];
    services: Service[];
    config: any;
    onUpdateConfig: (e: React.FormEvent) => void;
    setActiveTab: (tab: any) => void;
}

export const DashboardView: React.FC<DashboardViewProps> = ({
    status, peers, services, setActiveTab
}) => {
    return (
        <motion.div initial={{ opacity: 0, scale: 0.98 }} animate={{ opacity: 1, scale: 1 }} exit={{ opacity: 0, scale: 0.98 }} transition={{ duration: 0.2 }}>
            {/* Stats Overview */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
                <StatCard title="Network Peers" value={status?.peer_count ?? 0} sub="Authorized Nodes" icon={<Users className="text-blue-400" />} />
                <StatCard title="Active Proxies" value={status?.service_count ?? 0} sub="Public Gateways" icon={<Layers className="text-indigo-400" />} />
                <StatCard title="Tunnel Entry" value={status?.vps_ip || '---'} sub="VPS Hub Interface" icon={<Activity className="text-cyan-400" />} highlight="cyan" />
                <StatCard title="Compliance" value={status?.ready ? 'Secure' : 'Warning'} sub="SSL/DNS Health" icon={<ShieldCheck className="text-emerald-400" />} highlight="emerald" />
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div className="lg:col-span-2 space-y-8">
                    {/* Analytics Overview */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="glass-card p-8 flex items-center gap-6">
                            <div className="w-14 h-14 bg-blue-600/10 rounded-2xl flex items-center justify-center text-blue-400 shadow-inner">
                                <ArrowDown size={32} />
                            </div>
                            <div>
                                <h4 className="text-slate-500 text-xs font-bold uppercase tracking-widest mb-1">Global Received</h4>
                                <div className="text-3xl font-black text-white font-mono tracking-tighter">
                                    {formatBytes(status?.total_rx || 0)}
                                </div>
                            </div>
                        </div>
                        <div className="glass-card p-8 flex items-center gap-6">
                            <div className="w-14 h-14 bg-indigo-600/10 rounded-2xl flex items-center justify-center text-indigo-400 shadow-inner">
                                <ArrowUp size={32} />
                            </div>
                            <div>
                                <h4 className="text-slate-500 text-xs font-bold uppercase tracking-widest mb-1">Global Transmitted</h4>
                                <div className="text-3xl font-black text-white font-mono tracking-tighter">
                                    {formatBytes(status?.total_tx || 0)}
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Analytics Chart */}
                    <div className="glass-card">
                        <div className="px-8 py-6 border-b border-white/5 flex justify-between items-center">
                            <div className="flex items-center gap-3">
                                <BarChart3 size={20} className="text-blue-400" />
                                <h2 className="text-xl font-semibold">Tunnel Throughput (Live Trend)</h2>
                            </div>
                        </div>
                        <div className="p-8">
                            <UsageChart />
                        </div>
                    </div>
                </div>

                <div className="glass-card h-fit">
                    <div className="px-8 py-6 border-b border-white/5">
                        <h2 className="text-xl font-semibold flex items-center gap-3">
                            <Shield size={20} className="text-emerald-500" />
                            Network Status
                        </h2>
                    </div>
                    <div className="p-8 space-y-8">
                        <div className="flex justify-between items-center">
                            <span className="text-slate-400 text-sm font-medium">Project ID</span>
                            <span className="text-white font-bold">{status?.project}</span>
                        </div>
                        <div className="flex justify-between items-center">
                            <span className="text-slate-400 text-sm font-medium">Node Health</span>
                            <span className="text-emerald-400 font-bold">{peers.filter(p => p.live?.online).length} / {peers.length} Online</span>
                        </div>
                        <div className="flex justify-between items-center">
                            <span className="text-slate-400 text-sm font-medium">Tunnel Strategy</span>
                            <span className="text-blue-400 font-bold">Encapsulated</span>
                        </div>
                        <hr className="border-white/5" />
                        <div className="space-y-4">
                            <h4 className="text-xs font-bold text-slate-500 uppercase tracking-widest">Active Services</h4>
                            <div className="space-y-2">
                                {services.slice(0, 3).map(s => (
                                    <div key={s.domain} className="flex items-center justify-between text-sm py-2 px-3 bg-white/5 rounded-xl border border-white/5">
                                        <span className="text-slate-300 font-medium truncate max-w-[150px]">{s.domain}</span>
                                        <span className="text-blue-500 font-bold font-mono">{s.port}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                        <button onClick={() => setActiveTab('settings')} className="w-full flex items-center justify-center gap-2 py-3.5 bg-white/5 hover:bg-white/10 rounded-2xl text-sm font-bold border border-white/5 transition-all text-slate-300">
                            <Settings size={18} /> Cluster Settings
                        </button>
                    </div>
                </div>
            </div>
        </motion.div>
    );
};
