import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { apiClient } from '../api/client';

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
          const res = await apiClient.post('/users/anonymous', { anonymous_id: anonId });
          set({ userId: (res as any).user_id });
        } catch (error) {
          console.error('Failed to init user:', error);
        }
      },
    }),
    {
      name: 'city-roam-user',
    }
  )
);
