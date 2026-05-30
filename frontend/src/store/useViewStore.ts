import { create } from 'zustand';
import { City } from './useCityStore';

interface DrawerState {
  isOpen: boolean;
  type: 'chat' | 'gallery' | null;
  data: any | null;
}

interface ViewState {
  currentView: 'HOME' | 'EXPLORE' | 'ACHIEVEMENT';
  canvasMode: 'map' | 'street';
  streetTarget: City | null;
  drawer: DrawerState;
  setView: (view: 'HOME' | 'EXPLORE' | 'ACHIEVEMENT') => void;
  setCanvasMode: (mode: 'map' | 'street', target?: City) => void;
  openDrawer: (type: 'chat' | 'gallery', data: any) => void;
  closeDrawer: () => void;
}

export const useViewStore = create<ViewState>((set) => ({
  currentView: 'HOME',
  canvasMode: 'map',
  streetTarget: null,
  drawer: { isOpen: false, type: null, data: null },
  setView: (view) => set({ currentView: view }),
  setCanvasMode: (mode, target) => set((state) => ({ 
    canvasMode: mode, 
    streetTarget: target !== undefined ? target : state.streetTarget 
  })),
  openDrawer: (type, data) => set({ drawer: { isOpen: true, type, data } }),
  closeDrawer: () => set((state) => ({ drawer: { ...state.drawer, isOpen: false } })),
}));
