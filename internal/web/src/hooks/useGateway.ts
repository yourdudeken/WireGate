import { useState, useEffect, useCallback } from 'react';
import { gatewayApi } from '../api/gateway';
import { Status, Peer, Service } from '../types';
import { REFRESH_INTERVAL } from '../config/constants';

export function useGateway() {
    const [data, setData] = useState<{
        status: Status | null;
        peers: Peer[];
        services: Service[];
        config: any;
    }>({
        status: null,
        peers: [],
        services: [],
        config: null,
    });

    const [loading, setLoading] = useState(true);
    const [authError, setAuthError] = useState(false);

    const refresh = useCallback(async () => {
        try {
            const [status, peers, services, config] = await Promise.all([
                gatewayApi.getStatus(),
                gatewayApi.getPeers(),
                gatewayApi.getServices(),
                gatewayApi.getConfig(),
            ]);

            setData({ status, peers, services, config });
            setAuthError(false);
        } catch (err: any) {
            if (err.response?.status === 401) {
                setAuthError(true);
            }
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        refresh();
        const interval = setInterval(refresh, REFRESH_INTERVAL);

        const handleAuthError = () => setAuthError(true);
        window.addEventListener('auth-unauthorized', handleAuthError);

        return () => {
            clearInterval(interval);
            window.removeEventListener('auth-unauthorized', handleAuthError);
        };
    }, [refresh]);

    return { ...data, loading, authError, refresh };
}
