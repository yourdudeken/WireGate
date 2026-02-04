import React from 'react';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
    label: string;
}

export const Input: React.FC<InputProps> = ({ id, label, className = '', ...props }) => (
    <div className="flex flex-col gap-2 w-full">
        <label htmlFor={id} className="text-xs font-bold uppercase tracking-widest text-slate-500">{label}</label>
        <input
            id={id}
            className={`bg-slate-950/50 border border-white/5 rounded-xl p-3.5 text-white focus:outline-none focus:border-blue-500/50 focus:ring-4 focus:ring-blue-500/5 transition-all text-sm font-bold ${className}`}
            {...props}
        />
    </div>
);

interface StatCardProps {
    title: string;
    value: string | number;
    sub: string;
    icon: React.ReactNode;
    highlight?: string;
}

export const StatCard: React.FC<StatCardProps> = ({ title, value, sub, icon, highlight }) => (
    <div className="bg-slate-900/40 border border-white/5 rounded-3xl p-6 transition-all duration-300 hover:-translate-y-1 hover:border-white/10 hover:shadow-2xl">
        <div className="flex items-center gap-3 text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">
            {icon}
            {title}
        </div>
        <h2 className={`text-3xl font-bold mb-1 tracking-tight ${highlight === 'cyan' ? 'text-cyan-400' : highlight === 'emerald' ? 'text-emerald-400' : 'text-white'
            }`}>
            {value}
        </h2>
        <p className="text-sm text-slate-500 font-medium">{sub}</p>
    </div>
);
