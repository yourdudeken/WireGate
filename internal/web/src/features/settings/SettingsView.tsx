import React from 'react';
import { motion } from 'framer-motion';
import { Save } from 'lucide-react';
import { Status } from '../../types';
import { Input } from '../../components/ui/Cards';

interface SettingsViewProps {
    status: Status | null;
    config: any;
    onUpdateConfig: (e: React.FormEvent) => void;
}

export const SettingsView: React.FC<SettingsViewProps> = ({ status, config, onUpdateConfig }) => (
    <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0, y: -20 }}>
        <div className="max-w-3xl">
            <div className="glass-card mb-8">
                <div className="px-10 py-8 border-b border-white/5">
                    <h2 className="text-2xl font-bold text-white mb-2">Core Infrastructure</h2>
                    <p className="text-slate-500 text-sm font-medium">Global gateway parameters and project metadata.</p>
                </div>
                <div className="p-10">
                    <form onSubmit={onUpdateConfig} className="space-y-8">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                            <div className="space-y-6">
                                <Input id="vps-ip" label="VPS Public IP" defaultValue={status?.vps_ip} />
                                <Input id="vps-user" label="Gateway SSH User" defaultValue={status?.vps_user} />
                            </div>
                            <div className="space-y-6">
                                <Input id="proxy-email" label="Let's Encrypt Email" defaultValue={config?.proxy?.email} type="email" />
                                <Input id="project-name" label="Cluster Project ID" defaultValue={status?.project} />
                            </div>
                        </div>
                        <div className="pt-4">
                            <button type="submit" className="px-8 py-4 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded-2xl shadow-xl transition-all hover:shadow-blue-500/20 active:translate-y-0.5 flex items-center justify-center gap-3">
                                <Save size={20} /> Update Global Policy
                            </button>
                        </div>
                    </form>
                </div>
            </div>

            <div className="glass-card border-rose-500/10">
                <div className="px-10 py-8 border-b border-rose-500/5">
                    <h2 className="text-xl font-bold text-rose-500 mb-1">Danger Zone</h2>
                    <p className="text-slate-500 text-sm font-medium">Destructive cluster-wide operations.</p>
                </div>
                <div className="p-10 flex flex-col md:flex-row justify-between items-center gap-6">
                    <div>
                        <h4 className="text-white font-bold mb-1">Rotate WireGuard Keys</h4>
                        <p className="text-slate-500 text-sm">Regenerate all cryptographic identity keys. This will break current connections.</p>
                    </div>
                    <button className="px-6 py-3 bg-rose-600/10 hover:bg-rose-600 text-rose-500 hover:text-white font-bold rounded-xl border border-rose-500/20 transition-all shrink-0">
                        Initialize Rotation
                    </button>
                </div>
            </div>
        </div>
    </motion.div>
);
