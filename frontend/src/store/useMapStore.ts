import { create } from 'zustand';

interface MapStore {
  mapInstance: any | null;
  AMap: any | null;
  setMapContext: (map: any, AMap: any) => void;
  flyTo: (lng: number, lat: number, zoom?: number) => void;
  resetView: () => void;
}

export const useMapStore = create<MapStore>((set, get) => ({
  mapInstance: null,
  AMap: null,
  setMapContext: (map, AMap) => set({ mapInstance: map, AMap }),
  flyTo: (lng, lat, zoom = 14) => {
    const map = get().mapInstance;
    if (map) {
      map.setPitch(65, true, 800);
      setTimeout(() => {
        map.setZoomAndCenter(zoom, [lng, lat], false, 1500);
      }, 400);
    }
  },
  resetView: () => {
    const map = get().mapInstance;
    if (map) {
      map.setPitch(35, true, 800);
      setTimeout(() => {
        map.setZoomAndCenter(4.5, [104.195397, 35.86166], false, 1500);
      }, 400);
    }
  },
}));
