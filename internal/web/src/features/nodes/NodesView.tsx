import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Plus } from 'lucide-react';
import { Peer } from '../../types';
import { PeerCard } from './PeerCard';
import { TerminalModal } from './TerminalModal';

interface NodesViewProps {
    peers: Peer[];
    onAdd: () => void;
    onCopyKey: (key: string) => void;
}

export const NodesView: React.FC<NodesViewProps> = ({ peers, onAdd, onCopyKey }) => {
    const [selectedPeer, setSelectedPeer] = useState<string | null>(null);

    return (
        <motion.div initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} exit={{ opacity: 0, x: -20 }}>
            <div className="glass-card">
                <div className="px-10 py-8 border-b border-white/5 flex justify-between items-center">
                    <div>
                        <h2 className="text-2xl font-bold text-white mb-1 tracking-tight">Infrastructure Nodes</h2>
                        <p className="text-slate-500 text-sm font-medium">Authorized home servers and edge nodes in the WireGuard cluster.</p>
                    </div>
                    <button onClick={onAdd} className="bg-blue-600 hover:bg-blue-500 text-white px-6 py-3 rounded-2xl font-bold flex items-center gap-2 shadow-lg transition-all hover:-translate-y-0.5">
                        <Plus size={18} /> Authorize Node
                    </button>
                </div>
                <div className="p-10">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {peers.map(p => (
                            <PeerCard
                                key={p.name}
                                peer={p}
                                onCopyKey={onCopyKey}
                                onSSH={(name) => setSelectedPeer(name)}
                            />
                        ))}
                        {peers.length === 0 && <div className="col-span-full py-24 text-center text-slate-500 italic">No nodes detected.</div>}
                    </div>
                </div>
            </div>

            <TerminalModal
                isOpen={!!selectedPeer}
                onClose={() => setSelectedPeer(null)}
                peerName={selectedPeer || ''}
            />
        </motion.div>
    );
};
