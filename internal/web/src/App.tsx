import React, { useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { Lock, AlertCircle, CheckCircle2 } from 'lucide-react';

// Types & Config
import { Tab } from './types';
import { gatewayApi } from './api/gateway';

// Layout
import { Sidebar } from './components/layout/Sidebar';
import { Header } from './components/layout/Header';

// Components
import { Modal } from './components/ui/Modal';
import { Input } from './components/ui/Cards';

// Routing
import { AppRoutes } from './routes/AppRoutes';

// Providers
import { useGatewayContext } from './providers/GatewayProvider';

export default function App() {
    const [activeTab, setActiveTab] = useState<Tab>('dashboard');
    const [showPeerModal, setShowPeerModal] = useState(false);
    const [showServiceModal, setShowServiceModal] = useState(false);

    const {
        status,
        peers,
        authError,
        refresh,
        notification,
        showNotify
    } = useGatewayContext();

    const handleUpdateConfig = async (e: React.FormEvent) => {
        e.preventDefault();
        const form = e.target as HTMLFormElement;
        const updates = [
            { key: 'vps.ip', val: (form.querySelector('#vps-ip') as HTMLInputElement).value },
            { key: 'vps.user', val: (form.querySelector('#vps-user') as HTMLInputElement).value },
            { key: 'proxy.email', val: (form.querySelector('#proxy-email') as HTMLInputElement).value },
            { key: 'project', val: (form.querySelector('#project-name') as HTMLInputElement).value },
        ];

        try {
            await Promise.all(updates.map(u => gatewayApi.updateConfig(u.key, u.val)));
            showNotify('Cluster configuration synced');
            refresh();
        } catch (err) {
            showNotify('Failed to update config', 'error');
        }
    };

    const handleDeleteService = async (domain: string) => {
        if (!confirm(`Delete route ${domain}?`)) return;
        try {
            await gatewayApi.deleteService(domain);
            showNotify('Route terminated');
            refresh();
        } catch (err) {
            showNotify('Delete failed', 'error');
        }
    };

    const handleAddPeer = async (e: React.FormEvent) => {
        e.preventDefault();
        const form = e.target as HTMLFormElement;
        const name = (form.querySelector('#peer-name-input') as HTMLInputElement).value;
        const ip = (form.querySelector('#peer-ip-input') as HTMLInputElement).value;

        try {
            await gatewayApi.addPeer(name, ip);
            showNotify('Node authorized');
            setShowPeerModal(false);
            refresh();
        } catch (err: any) {
            showNotify(err.response?.data || 'Failed to add node', 'error');
        }
    };

    const handleAddService = async (e: React.FormEvent) => {
        e.preventDefault();
        const form = e.target as HTMLFormElement;
        const domain = (form.querySelector('#svc-domain') as HTMLInputElement).value;
        const port = parseInt((form.querySelector('#svc-port') as HTMLInputElement).value);
        const peer_name = (form.querySelector('#svc-peer') as HTMLSelectElement).value;

        try {
            await gatewayApi.addService(domain, port, peer_name);
            showNotify('Route established');
            setShowServiceModal(false);
            refresh();
        } catch (err: any) {
            showNotify(err.response?.data || 'Failed to add route', 'error');
        }
    };

    if (authError) {
        return (
            <div className="min-h-screen bg-[#090c12] flex items-center justify-center p-4">
                <div className="glass-card p-12 max-w-md w-full text-center">
                    <div className="w-20 h-20 bg-blue-600/20 text-blue-500 rounded-full flex items-center justify-center mx-auto mb-6">
                        <Lock size={40} />
                    </div>
                    <h1 className="text-2xl font-bold text-white mb-2">Authentication Required</h1>
                    <p className="text-slate-400 mb-8">Please refresh and provide your management credentials to access the gateway.</p>
                    <button onClick={() => window.location.reload()} className="w-full bg-blue-600 hover:bg-blue-500 text-white font-bold py-4 rounded-2xl transition-all">
                        Refresh Dashboard
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="flex min-h-screen bg-[#090c12] selection:bg-blue-500/30">
            <Sidebar activeTab={activeTab} setActiveTab={setActiveTab} />

            <main className="flex-1 p-8 md:p-12 max-w-7xl mx-auto w-full">
                <Header activeTab={activeTab} status={status} onRefresh={refresh} />

                <AppRoutes
                    activeTab={activeTab}
                    setActiveTab={setActiveTab}
                    onShowPeerModal={() => setShowPeerModal(true)}
                    onShowServiceModal={() => setShowServiceModal(true)}
                    onDeleteService={handleDeleteService}
                    onUpdateConfig={handleUpdateConfig}
                />
            </main>

            {/* Global Notifications */}
            <AnimatePresence>
                {notification && (
                    <motion.div
                        initial={{ opacity: 0, x: 100, scale: 0.9 }}
                        animate={{ opacity: 1, x: 0, scale: 1 }}
                        exit={{ opacity: 0, scale: 0.9 }}
                        className={`fixed bottom-8 right-8 z-[300] px-6 py-4 rounded-2xl shadow-2xl backdrop-blur-xl border border-white/10 flex items-center gap-4 ${notification.type === 'error' ? 'bg-rose-950/80 text-rose-200' : 'bg-slate-900/80 text-emerald-100'}`}
                    >
                        {notification.type === 'success' ? <CheckCircle2 className="text-emerald-400" /> : <AlertCircle className="text-rose-400" />}
                        <span className="font-bold text-sm">{notification.msg}</span>
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Shared Modals */}
            <Modal isOpen={showPeerModal} onClose={() => setShowPeerModal(false)} title="Authorize New Infrastructure Node">
                <form onSubmit={handleAddPeer} className="space-y-6">
                    <Input id="peer-name-input" label="Node Identifier" placeholder="e.g. warehouse-lab" />
                    <Input id="peer-ip-input" label="Assigned Tunnel IP" placeholder="e.g. 10.0.0.5" />
                    <div className="pt-4 flex gap-4">
                        <button type="button" onClick={() => setShowPeerModal(false)} className="flex-1 bg-white/5 text-slate-300 font-bold py-4 rounded-2xl border border-white/5 hover:bg-white/10 transition-all">Cancel</button>
                        <button type="submit" className="flex-1 bg-blue-600 text-white font-bold py-4 rounded-2xl shadow-lg hover:bg-blue-500 hover:-translate-y-0.5 transition-all">Generate Config</button>
                    </div>
                </form>
            </Modal>

            <Modal isOpen={showServiceModal} onClose={() => setShowServiceModal(false)} title="Establish New Routing Policy">
                <form onSubmit={handleAddService} className="space-y-6">
                    <Input id="svc-domain" label="Public Gateway Domain" placeholder="e.g. app.domain.com" />
                    <Input id="svc-port" label="Standard/Custom Service Port" type="number" placeholder="8080" />
                    <div className="flex flex-col gap-2">
                        <label className="text-xs font-bold uppercase tracking-widest text-slate-500">Destination Infrastructure Node</label>
                        <select id="svc-peer" className="bg-slate-950/50 border border-white/5 rounded-xl p-3.5 text-white focus:outline-none focus:border-blue-500/50 focus:ring-4 focus:ring-blue-500/5 transition-all w-full text-sm font-bold">
                            <option value="">Select target node...</option>
                            {peers.map(p => <option key={p.name} value={p.name}>{p.name} ({p.wg_ip})</option>)}
                        </select>
                    </div>
                    <div className="pt-4 flex gap-4">
                        <button type="button" onClick={() => setShowServiceModal(false)} className="flex-1 bg-white/5 text-slate-300 font-bold py-4 rounded-2xl border border-white/5 hover:bg-white/10 transition-all">Cancel</button>
                        <button type="submit" className="flex-1 bg-emerald-600 text-white font-bold py-4 rounded-2xl shadow-lg hover:bg-emerald-500 hover:-translate-y-0.5 transition-all">Establish Policy</button>
                    </div>
                </form>
            </Modal>
        </div>
    );
}
