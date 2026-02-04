import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X } from 'lucide-react';

interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    title: string;
    children: React.ReactNode;
}

export const Modal: React.FC<ModalProps> = ({ isOpen, onClose, title, children }) => {
    return (
        <AnimatePresence>
            {isOpen && (
                <div className="fixed inset-0 z-[200] flex items-center justify-center p-4">
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={onClose}
                        className="absolute inset-0 bg-black/70 backdrop-blur-md"
                    />
                    <motion.div
                        initial={{ opacity: 0, scale: 0.95, y: 30 }}
                        animate={{ opacity: 1, scale: 1, y: 0 }}
                        exit={{ opacity: 0, scale: 0.95, y: 20 }}
                        className="relative bg-[#0f172a] border border-white/10 rounded-[40px] w-full max-w-xl shadow-[0_40px_80px_-20px_rgba(0,0,0,0.8)] overflow-hidden"
                    >
                        <div className="px-12 pt-12 pb-6 flex justify-between items-center">
                            <h2 className="text-3xl font-bold text-white tracking-tight">{title}</h2>
                            <button onClick={onClose} className="p-2.5 text-slate-400 hover:text-white hover:bg-white/5 rounded-full transition-all">
                                <X size={28} />
                            </button>
                        </div>
                        <div className="px-12 pb-12">
                            {children}
                        </div>
                    </motion.div>
                </div>
            )}
        </AnimatePresence>
    );
};
