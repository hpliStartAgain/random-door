import { create } from 'zustand';
import { useCityStore } from './useCityStore';

interface GameState {
  rolling: boolean;
  targetCity: any | null;
  roll: () => Promise<any>;
  reset: () => void;
}

export const useGameStore = create<GameState>((set) => ({
  rolling: false,
  targetCity: null,
  roll: async () => {
    set({ rolling: true });
    return new Promise((resolve) => {
      setTimeout(() => {
        const cities = useCityStore.getState().cities;
        if (cities.length === 0) return resolve(null);
        
        // 随机抽取一个城市
        const randomCity = cities[Math.floor(Math.random() * cities.length)];
        set({ rolling: false, targetCity: randomCity });
        resolve(randomCity);
      }, 2000); // 骰子动画时长
    });
  },
  reset: () => set({ targetCity: null, rolling: false }),
}));
