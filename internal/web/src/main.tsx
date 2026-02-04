import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import { GatewayProvider } from './providers/GatewayProvider'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <GatewayProvider>
            <App />
        </GatewayProvider>
    </React.StrictMode>,
)
