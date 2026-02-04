import React from 'react';
import { Server, Key, Shield, Clock, Terminal, Wifi, ArrowDown, ArrowUp } from 'lucide-react';
import { Peer } from '../../types';
import { formatBytes } from '../../utils';

interface PeerCardProps {
    peer: Peer;
    onCopyKey: (key: string) => void;
    onSSH: (name: string) => void;
}

export const PeerCard: React.FC<PeerCardProps> = ({ peer, onCopyKey, onSSH }) => {
    const isOnline = peer.live?.online ?? false;
    const rx = peer.live?.rx ?? 0;
    const tx = peer.live?.tx ?? 0;

    return (
        <div className={`bg-slate-900/40 p-6 rounded-3xl border transition-all group relative overflow-hidden ${isOnline ? 'border-white/5 hover:border-blue-500/20' : 'border-rose-500/10 grayscale-[0.8]'}`}>
            <div className={`absolute top-0 right-0 w-32 h-32 rounded-full -mr-16 -mt-16 blur-3xl transition-all ${isOnline ? 'bg-blue-600/5 group-hover:bg-blue-600/10' : 'bg-rose-600/5'}`} />

            <div className="relative z-10">
                <div className="flex justify-between items-start mb-6">
                    <div className={`w-12 h-12 bg-white/5 rounded-2xl flex items-center justify-center transition-colors ${isOnline ? 'text-slate-400 group-hover:text-blue-400' : 'text-slate-600'}`}>
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
                        <div className={`flex items-center gap-2 text-xs font-bold ${isOnline ? 'text-emerald-500' : 'text-rose-500'}`}>
                            <div className={`w-1.5 h-1.5 rounded-full ${isOnline ? 'bg-emerald-500 animate-pulse' : 'bg-rose-500'}`} />
                            {isOnline ? 'ACTIVE TUNNEL' : 'OFFLINE'}
                        </div>
                    </div>
                    {isOnline && (
                        <div className="flex items-center gap-1.5 text-slate-500 text-[10px] font-bold uppercase tracking-tighter bg-white/5 px-2 py-1 rounded-lg">
                            <Wifi size={10} className="text-emerald-500" />
                            STABLE
                        </div>
                    )}
                </div>

                <div className="grid grid-cols-2 gap-4 mb-4">
                    <div className="bg-black/20 p-3 rounded-2xl border border-white/5">
                        <div className="flex items-center gap-2 text-[10px] font-bold text-slate-500 uppercase mb-1">
                            <ArrowDown size={10} className="text-blue-400" /> Received
                        </div>
                        <div className="text-sm font-bold text-white font-mono">{formatBytes(rx)}</div>
                    </div>
                    <div className="bg-black/20 p-3 rounded-2xl border border-white/5">
                        <div className="flex items-center gap-2 text-[10px] font-bold text-slate-500 uppercase mb-1">
                            <ArrowUp size={10} className="text-indigo-400" /> Sent
                        </div>
                        <div className="text-sm font-bold text-white font-mono">{formatBytes(tx)}</div>
                    </div>
                </div>

                <div className="space-y-2 pt-3 border-t border-white/5">
                    <div className="flex items-center gap-3 text-xs text-slate-400 font-medium">
                        <Shield size={12} className="text-slate-600" />
                        {peer.wg_ip}
                    </div>
                    {peer.live?.handshake ? (
                        <div className="flex items-center gap-3 text-xs text-slate-400 font-medium">
                            <Clock size={12} className="text-slate-600" />
                            Handshake: {Math.floor(peer.live.handshake / 60)}m ago
                        </div>
                    ) : null}
                </div>
            </div>
        </div>
    );
};
