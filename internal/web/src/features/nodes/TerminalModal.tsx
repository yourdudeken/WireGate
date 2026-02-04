import React, { useEffect, useRef } from 'react';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import '@xterm/xterm/css/xterm.css';
import { X } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface TerminalModalProps {
    isOpen: boolean;
    onClose: () => void;
    peerName: string;
}

export const TerminalModal: React.FC<TerminalModalProps> = ({ isOpen, onClose, peerName }) => {
    const terminalRef = useRef<HTMLDivElement>(null);
    const xtermRef = useRef<Terminal | null>(null);
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        if (isOpen && terminalRef.current && !xtermRef.current) {
            const term = new Terminal({
                cursorBlink: true,
                theme: {
                    background: '#0f172a',
                    foreground: '#f8fafc',
                    cursor: '#38bdf8',
                },
                fontSize: 14,
                fontFamily: 'Menlo, Monaco, "Courier New", monospace',
            });

            const fitAddon = new FitAddon();
            term.loadAddon(fitAddon);
            term.loadAddon(new WebLinksAddon());

            term.open(terminalRef.current);
            fitAddon.fit();
            xtermRef.current = term;

            // Connect to WebSocket
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/api/ssh?peer=${encodeURIComponent(peerName)}`;

            const ws = new WebSocket(wsUrl);
            wsRef.current = ws;

            ws.onopen = () => {
                term.write('\x1b[1;32m Connected to ' + peerName + ' via WireGuard\x1b[0m\r\n');
            };

            ws.onmessage = (event) => {
                if (event.data instanceof Blob) {
                    event.data.arrayBuffer().then(buf => {
                        term.write(new Uint8Array(buf));
                    });
                } else {
                    term.write(event.data);
                }
            };

            ws.onclose = () => {
                term.write('\r\n\x1b[1;31m Connection closed\x1b[0m\r\n');
            };

            term.onData(data => {
                if (ws.readyState === WebSocket.OPEN) {
                    ws.send(data);
                }
            });

            window.addEventListener('resize', () => fitAddon.fit());
        }

        return () => {
            if (!isOpen) {
                wsRef.current?.close();
                xtermRef.current?.dispose();
                xtermRef.current = null;
            }
        };
    }, [isOpen, peerName]);

    return (
        <AnimatePresence>
            {isOpen && (
                <div className="fixed inset-0 z-[250] flex items-center justify-center p-4">
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={onClose}
                        className="absolute inset-0 bg-black/80 backdrop-blur-xl"
                    />
                    <motion.div
                        initial={{ opacity: 0, scale: 0.9, y: 20 }}
                        animate={{ opacity: 1, scale: 1, y: 0 }}
                        exit={{ opacity: 0, scale: 0.9, y: 10 }}
                        className="relative bg-slate-950 border border-white/10 rounded-3xl w-full max-w-5xl h-[70vh] shadow-[0_50px_100px_-20px_rgba(0,0,0,0.5)] overflow-hidden flex flex-col"
                    >
                        <div className="px-8 py-4 border-b border-white/5 flex justify-between items-center bg-slate-900/50">
                            <div className="flex items-center gap-3">
                                <div className="w-3 h-3 rounded-full bg-rose-500" />
                                <div className="w-3 h-3 rounded-full bg-amber-500" />
                                <div className="w-3 h-3 rounded-full bg-emerald-500" />
                                <span className="ml-4 text-sm font-bold text-slate-400 font-mono">ssh {peerName}</span>
                            </div>
                            <button onClick={onClose} className="p-2 text-slate-500 hover:text-white transition-colors">
                                <X size={20} />
                            </button>
                        </div>
                        <div className="flex-1 p-4 overflow-hidden" ref={terminalRef} />
                    </motion.div>
                </div>
            )}
        </AnimatePresence>
    );
};
