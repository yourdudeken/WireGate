import React from 'react';
import { motion } from 'framer-motion';
import { Users, Layers, Activity, ShieldCheck, Globe, ChevronRight, Shield, Settings, BarChart3 } from 'lucide-react';
import { Status, Peer, Service } from '../../types';
import { StatCard } from '../../components/ui/Cards';
import { UsageChart } from './UsageChart';

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
                    {/* Analytics Chart */}
                    <div className="glass-card">
                        <div className="px-8 py-6 border-b border-white/5 flex justify-between items-center">
                            <div className="flex items-center gap-3">
                                <BarChart3 size={20} className="text-blue-400" />
                                <h2 className="text-xl font-semibold">Tunnel Throughput (Last 24h)</h2>
                            </div>
                        </div>
                        <div className="p-8">
                            <UsageChart />
                        </div>
                    </div>

                    {/* Traffic Redistribution */}
                    <div className="glass-card">
                        <div className="px-8 py-6 border-b border-white/5 flex justify-between items-center">
                            <div className="flex items-center gap-3">
                                <Globe size={20} className="text-blue-400" />
                                <h2 className="text-xl font-semibold">Live Traffic Distribution</h2>
                            </div>
                        </div>
                        <div className="p-8">
                            <div className="space-y-4">
                                {services.slice(0, 5).map(s => (
                                    <div key={s.domain} className="flex items-center justify-between p-4 bg-slate-900/30 rounded-2xl border border-white/5 group hover:border-white/10 transition-all">
                                        <div className="flex items-center gap-4">
                                            <div className="w-10 h-10 rounded-xl bg-blue-600/10 flex items-center justify-center text-blue-400">
                                                <Activity size={18} />
                                            </div>
                                            <div>
                                                <h4 className="font-bold text-white text-sm tracking-tight">{s.domain}</h4>
                                                <p className="text-xs text-slate-500 font-medium">Mapped to {s.peer_name}:{s.port}</p>
                                            </div>
                                        </div>
                                        <ChevronRight size={16} className="text-slate-600 group-hover:text-blue-400 transition-colors" />
                                    </div>
                                ))}
                                {services.length > 5 && (
                                    <button onClick={() => setActiveTab('services')} className="w-full py-4 text-sm font-bold text-blue-400 hover:text-blue-300 transition-colors bg-white/5 rounded-xl border border-white/5">
                                        View all active routes
                                    </button>
                                )}
                                {services.length === 0 && <p className="text-center py-8 text-slate-500 italic">No traffic routing detected.</p>}
                            </div>
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
                            <span className="text-slate-400 text-sm font-medium">Core Nodes</span>
                            <span className="text-emerald-400 font-bold">{peers.length} Healthy</span>
                        </div>
                        <div className="flex justify-between items-center">
                            <span className="text-slate-400 text-sm font-medium">Uptime Strategy</span>
                            <span className="text-blue-400 font-bold">Encapsulated VPN</span>
                        </div>
                        <hr className="border-white/5" />
                        <button onClick={() => setActiveTab('settings')} className="w-full flex items-center justify-center gap-2 py-3.5 bg-white/5 hover:bg-white/10 rounded-2xl text-sm font-bold border border-white/5 transition-all text-slate-300">
                            <Settings size={18} /> Cluster Settings
                        </button>
                    </div>
                </div>
            </div>
        </motion.div>
    );
};
