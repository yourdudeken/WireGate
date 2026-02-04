import React, { createContext, useContext, ReactNode, useState } from 'react';
import { useGateway } from '../hooks/useGateway';
import { Status, Peer, Service, Notification } from '../types';

interface GatewayContextType {
    status: Status | null;
    peers: Peer[];
    services: Service[];
    config: any;
    loading: boolean;
    authError: boolean;
    refresh: () => Promise<void>;
    notification: Notification | null;
    showNotify: (msg: string, type?: 'success' | 'error') => void;
}

const GatewayContext = createContext<GatewayContextType | undefined>(undefined);

export const GatewayProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const gateway = useGateway();
    const [notification, setNotification] = useState<Notification | null>(null);

    const showNotify = (msg: string, type: 'success' | 'error' = 'success') => {
        setNotification({ msg, type });
        setTimeout(() => setNotification(null), 4000);
    };

    const value = {
        ...gateway,
        notification,
        showNotify,
    };

    return (
        <GatewayContext.Provider value={value}>
            {children}
        </GatewayContext.Provider>
    );
};

export const useGatewayContext = () => {
    const context = useContext(GatewayContext);
    if (context === undefined) {
        throw new Error('useGatewayContext must be used within a GatewayProvider');
    }
    return context;
};
