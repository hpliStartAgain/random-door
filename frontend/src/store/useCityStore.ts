import { create } from 'zustand';
import { api } from '../api';
import type { CityListItem, CityDetail } from '../api/types';

export type City = CityListItem;

interface CityState {
  cities: City[];
  cityCache: Record<number, CityDetail>;
  loadCities: () => Promise<void>;
  loadCity: (id: number) => Promise<CityDetail>;
  reloadCity: (id: number) => Promise<CityDetail>;
  invalidateCity: (id?: number) => void;
}

export const useCityStore = create<CityState>((set, get) => ({
  cities: [],
  cityCache: {},
  loadCities: async () => {
    const res = await api.getCities();
    set({ cities: res.cities });
  },
  loadCity: async (id: number) => {
    const cached = get().cityCache[id];
    if (cached) return cached;
    const city = await api.getCityDetail(id);
    set((state) => ({ cityCache: { ...state.cityCache, [id]: city } }));
    return city;
  },
  reloadCity: async (id: number) => {
    const city = await api.getCityDetail(id);
    set((state) => ({ cityCache: { ...state.cityCache, [id]: city } }));
    return city;
  },
  invalidateCity: (id?: number) => set((state) => {
    if (id === undefined) return { cityCache: {} };
    const next = { ...state.cityCache };
    delete next[id];
    return { cityCache: next };
  }),
}));
