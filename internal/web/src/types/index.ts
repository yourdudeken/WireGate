export interface Status {
    project: string;
    vps_ip: string;
    vps_user: string;
    ready: boolean;
    peer_count: number;
    service_count: number;
    total_rx?: number;
    total_tx?: number;
}

export interface Peer {
    name: string;
    wg_ip: string;
    keepalive: number;
    public_key: string;
    live?: {
        handshake: number;
        rx: number;
        tx: number;
        online: boolean;
    };
}

export interface Service {
    domain: string;
    port: number;
    peer_name: string;
}

export type Tab = 'dashboard' | 'nodes' | 'services' | 'logs' | 'settings';

export interface Notification {
    msg: string;
    type: 'success' | 'error';
}
