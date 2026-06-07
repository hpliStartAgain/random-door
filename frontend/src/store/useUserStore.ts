import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { api } from '../api';

interface UserState {
  userId: number | null;
  anonymousId: string | null;
  username: string | null;
  nickname: string | null;
  currentCityId: number | null;
  initUser: () => Promise<void>;
  register: (username: string, password: string, nickname?: string) => Promise<void>;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  setCurrentCityId: (cityId: number | null) => void;
}

/** 兼容 HTTP 非安全上下文的 UUID v4 生成器（回退到 Math.random） */
function generateUUID(): string {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID();
  }
  // 回退：RFC 4122 UUID v4
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      userId: null,
      anonymousId: null,
      username: null,
      nickname: null,
      currentCityId: null,
      setCurrentCityId: (cityId) => set({ currentCityId: cityId }),
      initUser: async () => {
        let anonId = get().anonymousId;
        if (!anonId) {
          anonId = generateUUID();
          set({ anonymousId: anonId });
        }
        try {
          const res = await api.createAnonymousUser(anonId);
          set({ userId: res.user_id, currentCityId: res.current_city_id ?? get().currentCityId });
        } catch (error) {
          console.error('Failed to init user:', error);
        }
      },
      register: async (username, password, nickname) => {
        const res = await api.register({
          user_id: get().userId,
          username,
          password,
          nickname: nickname?.trim() || undefined,
        });
        set({
          userId: res.user_id,
          anonymousId: res.anonymous_id,
          username: res.username,
          nickname: res.nickname ?? null,
          currentCityId: res.current_city_id ?? null,
        });
      },
      login: async (username, password) => {
        const res = await api.login({ username, password });
        set({
          userId: res.user_id,
          anonymousId: res.anonymous_id,
          username: res.username,
          nickname: res.nickname ?? null,
          currentCityId: res.current_city_id ?? null,
        });
      },
      logout: async () => {
        const anonId = generateUUID();
        set({ userId: null, anonymousId: anonId, username: null, nickname: null, currentCityId: null });
        try {
          const res = await api.createAnonymousUser(anonId);
          set({ userId: res.user_id, currentCityId: res.current_city_id ?? null });
        } catch (error) {
          console.error('Failed to init anonymous user after logout:', error);
        }
      },
    }),
    { name: 'city-roam-user' }
  )
);
