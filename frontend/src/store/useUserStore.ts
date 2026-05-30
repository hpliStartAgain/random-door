import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { api } from '../api';

interface UserState {
  userId: number | null;
  anonymousId: string | null;
  currentCityId: number | null;
  initUser: () => Promise<void>;
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
      currentCityId: null,
      initUser: async () => {
        let anonId = get().anonymousId;
        if (!anonId) {
          anonId = generateUUID();
          set({ anonymousId: anonId });
        }
        try {
          const res = await api.createAnonymousUser(anonId);
          set({ userId: res.user_id });
        } catch (error) {
          console.error('Failed to init user:', error);
        }
      },
    }),
    { name: 'city-roam-user' }
  )
);
