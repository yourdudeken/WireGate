import axios from 'axios';
import { API_BASE_URL } from '../config/constants';

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Interceptor for handling global auth errors
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Custom event or global state can be triggered here
            window.dispatchEvent(new CustomEvent('auth-unauthorized'));
        }
        return Promise.reject(error);
    }
);
