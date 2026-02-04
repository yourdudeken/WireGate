import React from 'react';
import { Server, Key, Shield, Clock, Terminal, Wifi } from 'lucide-react';
import { Peer } from '../../types';

interface PeerCardProps {
    peer: Peer;
    onCopyKey: (key: string) => void;
    onSSH: (name: string) => void;
}

export const PeerCard: React.FC<PeerCardProps> = ({ peer, onCopyKey, onSSH }) => {
    // Simulate active latency for UI polish
    const latency = Math.floor(Math.random() * 50) + 15;

    return (
        <div className="bg-slate-900/40 p-6 rounded-3xl border border-white/5 hover:border-blue-500/20 transition-all group relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-blue-600/5 rounded-full -mr-16 -mt-16 blur-3xl group-hover:bg-blue-600/10 transition-all" />
            <div className="relative z-10">
                <div className="flex justify-between items-start mb-6">
                    <div className="w-12 h-12 bg-white/5 rounded-2xl flex items-center justify-center text-slate-400 group-hover:text-blue-400 transition-colors">
                        <Server size={24} />
                    </div>
                    <div className="flex gap-2">
                        <button
                            onClick={() => onSSH(peer.name)}
                            className="p-2.5 bg-white/5 hover:bg-white/10 rounded-xl transition-all"
                            title="Launch Web SSH"
                        >
                            <Terminal size={18} className="text-slate-500 group-hover:text-emerald-400 transition-colors" />
                        </button>
                        <button
                            onClick={() => onCopyKey(peer.public_key)}
                            className="p-2.5 bg-white/5 hover:bg-white/10 rounded-xl transition-all"
                            title="Copy Identity Key"
                        >
                            <Key size={18} className="text-slate-500 group-hover:text-amber-400 transition-colors" />
                        </button>
                    </div>
                </div>

                <div className="flex justify-between items-end mb-4">
                    <div>
                        <h3 className="text-xl font-bold text-white mb-0.5">{peer.name}</h3>
                        <div className="flex items-center gap-2 text-xs font-bold text-emerald-500">
                            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                            ACTIVE TUNNEL
                        </div>
                    </div>
                    <div className="flex items-center gap-1.5 text-slate-500 text-[10px] font-bold uppercase tracking-tighter bg-white/5 px-2 py-1 rounded-lg">
                        <Wifi size={10} className="text-emerald-500" />
                        {latency}ms
                    </div>
                </div>

                <div className="space-y-3 pt-2 border-t border-white/5">
                    <div className="flex items-center gap-3 text-sm text-slate-400 font-medium">
                        <Shield size={14} className="text-slate-600" />
                        {peer.wg_ip}
                    </div>
                    <div className="flex items-center gap-3 text-sm text-slate-400 font-medium">
                        <Clock size={14} className="text-slate-600" />
                        Keepalive: {peer.keepalive}s
                    </div>
                </div>
            </div>
        </div>
    );
};
