import React from 'react';
import { AnimatePresence } from 'framer-motion';
import { Tab } from '../types';
import { useGatewayContext } from '../providers/GatewayProvider';

// Views
import { DashboardView } from '../features/dashboard/DashboardView';
import { NodesView } from '../features/nodes/NodesView';
import { ServicesView } from '../features/services/ServicesView';
import { LogsView } from '../features/logs/LogsView';
import { SettingsView } from '../features/settings/SettingsView';

interface AppRoutesProps {
    activeTab: Tab;
    setActiveTab: (tab: Tab) => void;
    onShowPeerModal: () => void;
    onShowServiceModal: () => void;
    onDeleteService: (domain: string) => void;
    onUpdateConfig: (e: React.FormEvent) => void;
}

export const AppRoutes: React.FC<AppRoutesProps> = ({
    activeTab,
    setActiveTab,
    onShowPeerModal,
    onShowServiceModal,
    onDeleteService,
    onUpdateConfig
}) => {
    const { status, peers, services, config, refresh, showNotify } = useGatewayContext();

    return (
        <AnimatePresence mode="wait">
            {activeTab === 'dashboard' && (
                <DashboardView
                    key="dashboard"
                    status={status}
                    peers={peers}
                    services={services}
                    config={config}
                    onUpdateConfig={onUpdateConfig}
                    setActiveTab={setActiveTab}
                />
            )}

            {activeTab === 'nodes' && (
                <NodesView
                    key="nodes"
                    peers={peers}
                    onAdd={onShowPeerModal}
                    onCopyKey={(key) => {
                        navigator.clipboard.writeText(key);
                        showNotify('Public key copied');
                    }}
                />
            )}

            {activeTab === 'services' && (
                <ServicesView
                    key="services"
                    services={services}
                    onAddService={onShowServiceModal}
                    onDeleteService={onDeleteService}
                />
            )}

            {activeTab === 'logs' && (
                <LogsView key="logs" status={status} peers={peers} services={services} />
            )}

            {activeTab === 'settings' && (
                <SettingsView
                    key="settings"
                    status={status}
                    config={config}
                    onUpdateConfig={onUpdateConfig}
                />
            )}
        </AnimatePresence>
    );
};
