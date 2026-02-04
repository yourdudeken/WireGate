import { apiClient } from './client';
import { Status, Peer, Service } from '../types';

export const gatewayApi = {
    getStatus: () => apiClient.get<Status>('/api/status').then(r => r.data),
    getPeers: () => apiClient.get<Peer[]>('/api/peers').then(r => r.data),
    getServices: () => apiClient.get<Service[]>('/api/services').then(r => r.data),
    getConfig: () => apiClient.get<any>('/api/config').then(r => r.data),

    updateConfig: (key: string, value: string) =>
        apiClient.post('/api/config/update', { key, value }),

    addPeer: (name: string, ip: string) =>
        apiClient.post('/api/peers/add', { name, ip }),

    addService: (domain: string, port: number, peer_name: string) =>
        apiClient.post('/api/services/add', { domain, port, peer_name }),

    deleteService: (domain: string) =>
        apiClient.delete(`/api/services/delete?domain=${encodeURIComponent(domain)}`),
};
