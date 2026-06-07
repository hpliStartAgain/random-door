import { create } from 'zustand';

interface DrawerState {
  isOpen: boolean;
  type: 'chat' | 'gallery' | null;
  data: any | null;
}

export interface StreetTarget {
  id: number;
  city_id?: number;
  name: string;
  province?: string;
  lat?: number;
  lng?: number;
  cover_image_url?: string;
  soundscape_url?: string;
  tags?: string[];
}

export type ViewName = 'HOME' | 'EXPLORE' | 'ACHIEVEMENT' | 'ASSETS' | 'CITY_DETAIL' | 'FREE_EXPLORE' | 'GAME_DICE';
export type RollPhase = 'idle' | 'rolling' | 'revealing' | 'flying' | 'landed';

interface ViewState {
  currentView: ViewName;
  activeCityId: number | null;
  canvasMode: 'map' | 'street';
  streetTarget: StreetTarget | null;
  drawer: DrawerState;
  hasEntered: boolean;
  rollPhase: RollPhase;
  setView: (view: ViewName, cityId?: number) => void;
  setCanvasMode: (mode: 'map' | 'street', target?: StreetTarget) => void;
  openDrawer: (type: 'chat' | 'gallery', data: any) => void;
  closeDrawer: () => void;
  enter: () => void;
  setRollPhase: (phase: RollPhase) => void;
  profileOpen: boolean;
  setProfileOpen: (open: boolean) => void;
}

export const useViewStore = create<ViewState>((set) => ({
  currentView: 'HOME',
  activeCityId: null,
  canvasMode: 'map',
  streetTarget: null,
  drawer: { isOpen: false, type: null, data: null },
  hasEntered: false,
  rollPhase: 'idle',
  profileOpen: false,
  setView: (view, cityId) => set({ currentView: view, activeCityId: cityId ?? null }),
  setCanvasMode: (mode, target) => set((state) => ({ 
    canvasMode: mode, 
    streetTarget: target !== undefined ? target : state.streetTarget 
  })),
  openDrawer: (type, data) => set({ drawer: { isOpen: true, type, data } }),
  closeDrawer: () => set((state) => ({ drawer: { ...state.drawer, isOpen: false } })),
  enter: () => set({ hasEntered: true }),
  setRollPhase: (phase) => set({ rollPhase: phase }),
  setProfileOpen: (open) => set({ profileOpen: open }),
}));
