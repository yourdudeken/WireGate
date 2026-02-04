import React from 'react';
import { motion } from 'framer-motion';
import { PlusCircle, Server, ExternalLink, Trash2 } from 'lucide-react';
import { Service } from '../../types';

interface ServicesViewProps {
    services: Service[];
    onAddService: () => void;
    onDeleteService: (domain: string) => void;
}

export const ServicesView: React.FC<ServicesViewProps> = ({ services, onAddService, onDeleteService }) => (
    <motion.div initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} exit={{ opacity: 0, x: -20 }}>
        <div className="glass-card">
            <div className="px-10 py-8 border-b border-white/5 flex justify-between items-center">
                <div>
                    <h2 className="text-2xl font-bold text-white mb-1 tracking-tight">Traffic Routing Rules</h2>
                    <p className="text-slate-500 text-sm font-medium">Mapping public domains to internal secure tunnel workloads.</p>
                </div>
                <button onClick={onAddService} className="bg-emerald-600 hover:bg-emerald-500 text-white px-6 py-3 rounded-2xl font-bold flex items-center gap-2 shadow-lg transition-all hover:-translate-y-0.5">
                    <PlusCircle size={18} /> Establish New Route
                </button>
            </div>
            <div className="p-0 overflow-hidden">
                <table className="w-full text-left">
                    <thead>
                        <tr className="bg-slate-950/20 text-slate-500 text-xs font-bold uppercase tracking-widest">
                            <th className="py-6 px-10">Access Domain</th>
                            <th className="py-6 px-10 text-center">Protocol Port</th>
                            <th className="py-6 px-10">Destination Node</th>
                            <th className="py-6 px-10 text-right">Cluster Action</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5">
                        {services.map(s => (
                            <tr key={s.domain} className="group hover:bg-white/[0.01] transition-all">
                                <td className="py-8 px-10 font-bold text-white tracking-tight flex items-center gap-4 group-hover:text-blue-400 transition-colors">
                                    <div className="w-2 h-2 rounded-full bg-emerald-500" />
                                    {s.domain}
                                </td>
                                <td className="py-8 px-10 text-center">
                                    <span className="inline-block bg-slate-900 border border-white/5 px-3 py-1.5 rounded-xl font-mono text-sm text-slate-400 font-bold">
                                        {s.port}
                                    </span>
                                </td>
                                <td className="py-8 px-10">
                                    <div className="flex items-center gap-3 text-slate-300 font-bold">
                                        <div className="w-8 h-8 bg-blue-600/10 rounded-lg flex items-center justify-center text-blue-500">
                                            <Server size={14} />
                                        </div>
                                        {s.peer_name}
                                    </div>
                                </td>
                                <td className="py-8 px-10 text-right">
                                    <div className="flex justify-end gap-3">
                                        <a href={`https://${s.domain}`} target="_blank" rel="noreferrer" className="p-3 bg-white/5 hover:bg-white/10 rounded-xl text-slate-400 transition-all">
                                            <ExternalLink size={18} />
                                        </a>
                                        <button onClick={() => onDeleteService(s.domain)} className="p-3 bg-rose-600/5 hover:bg-rose-600 hover:text-white rounded-xl text-rose-500 transition-all shadow-sm">
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                        {services.length === 0 && <tr><td colSpan={4} className="py-32 text-center text-slate-500 italic">No routing policies established.</td></tr>}
                    </tbody>
                </table>
            </div>
        </div>
    </motion.div>
);
