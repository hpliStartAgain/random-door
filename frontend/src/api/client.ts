import axios from 'axios';
import { useUserStore } from '../store/useUserStore';

export const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000,
});

apiClient.interceptors.request.use((config) => {
  const userId = useUserStore.getState().userId;
  if (userId) {
    config.headers['X-User-Id'] = String(userId);
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    // 统一错误处理，按 api-contract 格式
    const err = error.response?.data?.error || {
      code: 'UNKNOWN_ERROR',
      message: error.message || '网络请求错误',
    };
    return Promise.reject(err);
  }
);
