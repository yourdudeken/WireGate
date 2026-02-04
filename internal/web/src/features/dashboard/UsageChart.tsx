import React from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

const data = [
    { time: '10:00', traffic: 400 },
    { time: '11:00', traffic: 300 },
    { time: '12:00', traffic: 200 },
    { time: '13:00', traffic: 278 },
    { time: '14:00', traffic: 189 },
    { time: '15:00', traffic: 239 },
    { time: '16:00', traffic: 349 },
];

export const UsageChart: React.FC = () => {
    return (
        <div className="h-[300px] w-full mt-4">
            <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={data}>
                    <defs>
                        <linearGradient id="colorTraffic" x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3} />
                            <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                        </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="#ffffff05" vertical={false} />
                    <XAxis
                        dataKey="time"
                        axisLine={false}
                        tickLine={false}
                        tick={{ fill: '#64748b', fontSize: 12 }}
                        dy={10}
                    />
                    <YAxis
                        hide
                    />
                    <Tooltip
                        contentStyle={{
                            backgroundColor: '#0f172a',
                            border: '1px solid rgba(255,255,255,0.1)',
                            borderRadius: '12px',
                            fontSize: '12px',
                            fontWeight: 'bold'
                        }}
                        itemStyle={{ color: '#38bdf8' }}
                    />
                    <Area
                        type="monotone"
                        dataKey="traffic"
                        stroke="#3b82f6"
                        strokeWidth={3}
                        fillOpacity={1}
                        fill="url(#colorTraffic)"
                    />
                </AreaChart>
            </ResponsiveContainer>
        </div>
    );
};
