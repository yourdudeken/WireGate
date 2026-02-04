import React, { useState, useEffect } from 'react';
import {
    LayoutDashboard,
    Server,
    Globe,
    Terminal,
    Settings,
    RefreshCw,
    Users,
    Layers,
    Activity,
    ShieldCheck,
    Save,
    Plus,
    PlusCircle,
    X,
    Trash2,
    Key,
    Clock,
    Shield
} from 'lucide-react';
import axios from 'axios';
import { motion, AnimatePresence } from 'framer-motion';

const API_BASE = '';

interface Status {
    project: string;
    vps_ip: string;
    vps_user: string;
    ready: boolean;
    peer_count: number;
    service_count: number;
}

interface Peer {
    name: string;
    wg_ip: string;
    keepalive: number;
    public_key: string;
}

interface Service {
    domain: string;
    port: number;
    peer_name: string;
}

export default function App() {
    const [status, setStatus] = useState<Status | null>(null);
    const [peers, setPeers] = useState<Peer[]>([]);
    const [services, setServices] = useState<Service[]>([]);
    const [config, setConfig] = useState<any>(null);
    const [loading, setLoading] = useState(true);

    const [showPeerModal, setShowPeerModal] = useState(false);
    const [showServiceModal, setShowServiceModal] = useState(false);

    useEffect(() => {
        refreshAll();
        const interval = setInterval(refreshAll, 10000);
        return () => clearInterval(interval);
    }, []);

    const refreshAll = async () => {
        try {
            const [sRes, pRes, svRes, cRes] = await Promise.all([
                axios.get(`${API_BASE}/api/status`),
                axios.get(`${API_BASE}/api/peers`),
                axios.get(`${API_BASE}/api/services`),
                axios.get(`${API_BASE}/api/config`)
            ]);
            setStatus(sRes.data);
            setPeers(pRes.data || []);
            setServices(svRes.data || []);
            setConfig(cRes.data);
            setLoading(false);
        } catch (err) {
            console.error('Fetch error:', err);
        }
    };

    const handleUpdateConfig = async (e: React.FormEvent) => {
        e.preventDefault();
        const form = e.target as HTMLFormElement;
        const updates = [
            { key: 'vps.ip', value: (form.querySelector('#vps-ip') as HTMLInputElement).value },
            { key: 'vps.user', value: (form.querySelector('#vps-user') as HTMLInputElement).value },
            { key: 'proxy.email', value: (form.querySelector('#proxy-email') as HTMLInputElement).value },
            { key: 'project', value: (form.querySelector('#project-name') as HTMLInputElement).value },
        ];

        try {
            for (const u of updates) {
                await axios.post(`${API_BASE}/api/config/update`, u);
            }
            refreshAll();
        } catch (err) {
            alert('Failed to update config');
        }
    };

    const handleDeleteService = async (domain: string) => {
        if (!confirm(`Delete route ${domain}?`)) return;
        try {
            await axios.delete(`${API_BASE}/api/services/delete?domain=${encodeURIComponent(domain)}`);
            refreshAll();
        } catch (err) {
            alert('Delete failed');
        }
    };

    return (
        <div className="flex min-h-screen bg-[#090c12]">
            {/* Sidebar */}
            <aside className="w-20 bg-slate-950/80 backdrop-blur-2xl border-r border-white/5 flex flex-col items-center py-8 sticky top-0 h-screen z-50">
                <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center font-bold text-xl shadow-[0_0_20px_rgba(59,130,246,0.5)] mb-12">
                    WG
                </div>
                <nav className="flex flex-col gap-6 flex-1">
                    <NavItem active icon={<LayoutDashboard size={24} />} />
                    <NavItem icon={<Server size={24} />} />
                    <NavItem icon={<Globe size={24} />} />
                    <NavItem icon={<Terminal size={24} />} />
                    <div className="flex-1" />
                    <NavItem icon={<Settings size={24} />} />
                </nav>
            </aside>

            {/* Main */}
            <main className="flex-1 p-8 md:p-12 max-w-7xl mx-auto w-full">
                <header className="flex justify-between items-end mb-12">
                    <div>
                        <h1 className="text-4xl font-bold tracking-tight text-white mb-2">W-G Gateway</h1>
                        <p className="text-slate-400 font-medium">Cloud-to-Home Infrastructure Command Center</p>
                    </div>
                    <div className="flex items-center gap-4">
                        <div className="bg-slate-900/50 px-4 py-2 rounded-full border border-white/5 flex items-center gap-3">
                            <div className={`w-2 h-2 rounded-full ${status?.ready ? 'bg-emerald-500 shadow-[0_0_10px_#10b981]' : 'bg-amber-500'} animate-pulse`} />
                            <span className="text-sm font-semibold">{status?.ready ? 'Gateway Online' : 'Pending Setup'}</span>
                        </div>
                        <button onClick={refreshAll} className="p-2.5 rounded-xl bg-white/5 hover:bg-white/10 transition-colors border border-white/5">
                            <RefreshCw size={20} className="text-slate-300" />
                        </button>
                    </div>
                </header>

                {/* Stats */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
                    <StatCard title="Network Peers" value={status?.peer_count || 0} sub="Connected Nodes" icon={<Users className="text-blue-400" />} />
                    <StatCard title="Active Proxies" value={status?.service_count || 0} sub="Routed Services" icon={<Layers className="text-indigo-400" />} />
                    <StatCard title="Gateway Entry" value={status?.vps_ip || '-'} sub="VPS IP Address" icon={<Activity className="text-cyan-400" />} highlight="cyan" />
                    <StatCard title="System Health" value={status?.ready ? 'Secure' : 'Alert'} sub="Security Audit" icon={<ShieldCheck className="text-emerald-400" />} highlight="emerald" />
                </div>

                {/* Content Row 1 */}
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
                    {/* Config */}
                    <div className="glass-card">
                        <div className="px-8 py-6 border-b border-white/5 flex items-center gap-3">
                            <Settings size={20} className="text-blue-500" />
                            <h2 className="text-xl font-semibold">Global Configuration</h2>
                        </div>
                        <div className="p-8">
                            <form onSubmit={handleUpdateConfig} className="space-y-6">
                                <div className="grid grid-cols-2 gap-6">
                                    <Input id="vps-ip" label="VPS IP" defaultValue={status?.vps_ip} />
                                    <Input id="vps-user" label="SSH User" defaultValue={status?.vps_user} />
                                    <Input id="proxy-email" label="Admin Email" defaultValue={config?.proxy?.email} type="email" />
                                    <Input id="project-name" label="Project Name" defaultValue={status?.project} />
                                </div>
                                <button type="submit" className="w-full bg-blue-600 hover:bg-blue-500 text-white font-bold py-3 rounded-2xl shadow-xl transition-all hover:-translate-y-0.5 active:translate-y-0 flex items-center justify-center gap-2">
                                    <Save size={20} /> Save Cluster State
                                </button>
                            </form>
                        </div>
                    </div>

                    {/* Nodes */}
                    <div className="glass-card">
                        <div className="px-8 py-6 border-b border-white/5 flex justify-between items-center">
                            <div className="flex items-center gap-3">
                                <Server size={20} className="text-indigo-500" />
                                <h2 className="text-xl font-semibold">Network Nodes</h2>
                            </div>
                            <button disabled className="bg-white/5 px-4 py-2 rounded-xl text-sm font-semibold text-slate-400 border border-white/5 flex items-center gap-2">
                                <Plus size={16} /> New Peer
                            </button>
                        </div>
                        <div className="p-8">
                            <div className="space-y-4 max-h-[400px] overflow-y-auto">
                                {peers.map(p => (
                                    <div key={p.name} className="bg-slate-950/40 p-4 rounded-2xl border border-white/5 flex items-center justify-between hover:border-white/10 transition-colors">
                                        <div>
                                            <h4 className="font-bold text-white mb-1">{p.name}</h4>
                                            <div className="flex gap-4 text-xs font-semibold text-slate-500">
                                                <span className="flex items-center gap-1"><Shield size={12} /> {p.wg_ip}</span>
                                                <span className="flex items-center gap-1"><Clock size={12} /> {p.keepalive}s</span>
                                            </div>
                                        </div>
                                        <button onClick={() => {
                                            navigator.clipboard.writeText(p.public_key);
                                            alert('Key copied');
                                        }} className="p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors">
                                            <Key size={16} className="text-slate-400" />
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>

                {/* Services Row */}
                <div className="glass-card">
                    <div className="px-8 py-6 border-b border-white/5 flex justify-between items-center">
                        <div className="flex items-center gap-3">
                            <Globe size={20} className="text-blue-400" />
                            <h2 className="text-xl font-semibold">Traffic Routing Rules</h2>
                        </div>
                        <button disabled className="bg-blue-600/10 text-blue-400 px-5 py-2.5 rounded-xl font-bold border border-blue-500/20 flex items-center gap-2">
                            <PlusCircle size={18} /> Deploy Service
                        </button>
                    </div>
                    <div className="p-8">
                        <div className="overflow-x-auto">
                            <table className="w-full text-left">
                                <thead>
                                    <tr className="text-slate-500 text-xs font-bold uppercase tracking-widest">
                                        <th className="pb-6 px-4">Access Domain</th>
                                        <th className="pb-6 px-4 text-center">Port</th>
                                        <th className="pb-6 px-4">Destination Node</th>
                                        <th className="pb-6 px-4 text-right">Action</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-white/5">
                                    {services.map(s => (
                                        <tr key={s.domain} className="group hover:bg-white/[0.02] transition-colors">
                                            <td className="py-6 px-4 font-bold text-blue-400">{s.domain}</td>
                                            <td className="py-6 px-4 text-center">
                                                <span className="bg-slate-950/50 px-2.5 py-1 rounded-md border border-white/5 font-mono text-sm">
                                                    {s.port}
                                                </span>
                                            </td>
                                            <td className="py-6 px-4">
                                                <div className="flex items-center gap-2 text-slate-300 font-semibold">
                                                    <Server size={14} className="text-slate-500" />
                                                    {s.peer_name}
                                                </div>
                                            </td>
                                            <td className="py-6 px-4 text-right">
                                                <button onClick={() => handleDeleteService(s.domain)} className="p-2.5 rounded-xl text-rose-500 hover:bg-rose-500/10 hover:shadow-[0_0_15px_rgba(244,63,94,0.2)] transition-all">
                                                    <Trash2 size={20} />
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
}

function NavItem({ icon, active = false }: { icon: React.ReactNode, active?: boolean }) {
    return (
        <a href="#" className={`w-12 h-12 flex items-center justify-center rounded-2xl transition-all duration-300 ${active ? 'bg-blue-600/10 text-blue-500 shadow-[inset_0_0_0_1px_rgba(37,99,235,0.2)]' : 'text-slate-500 hover:text-white hover:bg-white/5'}`}>
            {icon}
        </a>
    );
}

function StatCard({ title, value, sub, icon, highlight }: { title: string, value: string | number, sub: string, icon: React.ReactNode, highlight?: string }) {
    return (
        <div className="stat-card">
            <div className="flex items-center gap-3 text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">
                {icon}
                {title}
            </div>
            <h2 className={`text-3xl font-bold mb-2 ${highlight === 'cyan' ? 'text-cyan-400' : highlight === 'emerald' ? 'text-emerald-400' : 'text-white'}`}>
                {value}
            </h2>
            <p className="text-sm text-slate-500 font-medium">{sub}</p>
        </div>
    );
}

function Input({ id, label, type = 'text', defaultValue }: { id: string, label: string, type?: string, defaultValue?: string }) {
    return (
        <div className="flex flex-col gap-2">
            <label htmlFor={id} className="text-xs font-bold uppercase tracking-widest text-slate-500">{label}</label>
            <input id={id} type={type} defaultValue={defaultValue} className="bg-slate-950/50 border border-white/5 rounded-xl p-3.5 text-white focus:outline-none focus:border-blue-500/50 focus:ring-4 focus:ring-blue-500/5 transition-all" />
        </div>
    );
}
