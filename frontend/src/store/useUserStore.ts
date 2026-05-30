import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { api } from '../api';

interface UserState {
  userId: number | null;
  anonymousId: string | null;
  currentCityId: number | null;
  initUser: () => Promise<void>;
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
          anonId = crypto.randomUUID();
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
