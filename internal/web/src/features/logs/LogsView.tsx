import React from 'react';
import { motion } from 'framer-motion';
import { Terminal } from 'lucide-react';
import { Status, Peer, Service } from '../../types';

interface LogsViewProps {
    status: Status | null;
    peers: Peer[];
    services: Service[];
}

export const LogsView: React.FC<LogsViewProps> = ({ status, peers, services }) => (
    <motion.div initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }} exit={{ opacity: 0, scale: 0.95 }} className="space-y-6">
        <div className="glass-card flex flex-col h-[70vh]">
            <div className="px-10 py-6 border-b border-white/5 bg-slate-950/20 flex justify-between items-center">
                <div className="flex items-center gap-3">
                    <Terminal size={20} className="text-emerald-500" />
                    <h2 className="text-lg font-bold">System Runtime Log</h2>
                </div>
                <div className="flex items-center gap-4">
                    <span className="flex items-center gap-2 text-xs font-bold text-emerald-500 uppercase tracking-widest animate-pulse">
                        <div className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                        Streaming Live
                    </span>
                </div>
            </div>
            <div className="flex-1 p-8 font-mono text-sm overflow-y-auto bg-black/40 text-slate-400 custom-scrollbar">
                <div className="space-y-1">
                    <p><span className="text-emerald-600 font-bold">[SYS]</span> Cluster controller initialized version 2.2.2</p>
                    <p><span className="text-blue-600 font-bold">[WEB]</span> Dashboard authenticated session active for admin</p>
                    <p><span className="text-slate-600 font-bold">[TUNNEL]</span> WireGuard interface wg0 is UP (10.0.0.1)</p>
                    <p><span className="text-slate-600 font-bold">[PROXY]</span> Traefik core engine ready at {status?.vps_ip}</p>
                    {peers.map(p => (
                        <p key={p.name}><span className="text-indigo-600 font-bold">[NODE]</span> Node "{p.name}" status: CONNECTED over {p.wg_ip}</p>
                    ))}
                    {services.map(s => (
                        <p key={s.domain}><span className="text-emerald-600 font-bold">[PROXY]</span> Routing established: https://{s.domain} -{'>'} {s.peer_name}:{s.port}</p>
                    ))}
                    <p className="animate-pulse">_</p>
                </div>
            </div>
        </div>
    </motion.div>
);
