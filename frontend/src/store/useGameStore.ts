import { create } from 'zustand';
import { api } from '../api';
import type { GameRollResponse, NearestCity } from '../api/types';

interface GameState {
  rolling: boolean;
  nearestCity: NearestCity | null;
  targetCity: NearestCity | null;
  lastRoll: GameRollResponse | null;
  initGame: (userId: number, lat?: number, lng?: number) => Promise<NearestCity>;
  roll: (userId: number, fromCityId: number, lat: number, lng: number) => Promise<GameRollResponse>;
  reset: () => void;
}

export const useGameStore = create<GameState>((set) => ({
  rolling: false,
  nearestCity: null,
  targetCity: null,
  lastRoll: null,
  initGame: async (userId, lat, lng) => {
    const res = await api.gameInit(userId, lat, lng);
    set({ nearestCity: res.nearest_city });
    return res.nearest_city;
  },
  roll: async (userId, fromCityId, lat, lng) => {
    set({ rolling: true });
    try {
      const res = await api.gameRoll(userId, fromCityId, lat, lng);
      set({ rolling: false, targetCity: res.target_city, lastRoll: res });
      return res;
    } catch (e) {
      set({ rolling: false });
      throw e;
    }
  },
  reset: () => set({ targetCity: null, lastRoll: null, rolling: false }),
}));
