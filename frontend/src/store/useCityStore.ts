import { create } from 'zustand';

export interface Figure {
  name: string;
  dynasty: string;
  desc: string;
}

export interface Food {
  name: string;
  desc: string;
}

export interface City {
  id: number;
  name: string;
  province: string;
  lat: number;
  lng: number;
  intro?: string;
  tags?: string[];
  figures?: Figure[];
  foods?: Food[];
  panoramaUrl?: string;
}

const MOCK_CITIES: City[] = [
  { 
    id: 1, name: '北京', province: '北京市', lat: 39.9042, lng: 116.4074, 
    intro: '三朝古都，长城横亘，故宫恢弘，历史与现代在此激烈交锋。', tags: ['故宫', '长城'],
    panoramaUrl: 'https://pannellum.org/images/bma-0.jpg',
    figures: [
      { name: '朱棣', dynasty: '明朝', desc: '明成祖，营建北京紫禁城，迁都北京，奠定了北京作为明清两代帝都的基础。' },
      { name: '老舍', dynasty: '近现代', desc: '著名文学家，以京味文学著称，《茶馆》《骆驼祥子》刻画了深刻的老北京平民风貌。' }
    ],
    foods: [
      { name: '北京烤鸭', desc: '肉质肥而不腻，外脆里嫩，享誉世界的皇家御膳代表。' },
      { name: '炸酱面', desc: '老北京人的家常面食，菜码讲究，酱香浓郁。' }
    ]
  },
  { 
    id: 2, name: '西安', province: '陕西省', lat: 34.3416, lng: 108.9398, 
    intro: '十三朝古都，地下沉睡着秦皇大军，地上回荡着大唐胡旋。', tags: ['兵马俑', '大唐'],
    panoramaUrl: 'https://pannellum.org/images/alma.jpg',
    figures: [
      { name: '秦始皇', dynasty: '秦朝', desc: '扫平六国，一统天下，在西安附近的骊山留下了震撼世界的兵马俑坑。' },
      { name: '李白', dynasty: '唐朝', desc: '诗仙，在长安留下了无数浪漫不羁的诗篇，长安市上酒家眠。' },
      { name: '武则天', dynasty: '唐/武周', desc: '中国历史上唯一正统的女皇帝，长期以长安和洛阳为执政中心。' }
    ],
    foods: [
      { name: '羊肉泡馍', desc: '肉烂汤浓，香气四溢，自己掰馍的乐趣也是体验的一部分。' },
      { name: '肉夹馍', desc: '腊汁肉加上白吉馍，号称“中式汉堡”，历史悠久。' }
    ]
  },
  { id: 3, name: '洛阳', province: '河南省', lat: 34.6186, lng: 112.4540, intro: '若问古今兴废事，请君只看洛阳城。', tags: ['龙门石窟', '牡丹'] },
  { id: 4, name: '南京', province: '江苏省', lat: 32.0603, lng: 118.7969, intro: '六朝金粉地，十里秦淮河。这里充满了南方的柔媚与历史的沧桑。', tags: ['秦淮河', '六朝'] },
  { id: 5, name: '杭州', province: '浙江省', lat: 30.2741, lng: 120.1551, intro: '江南忆，最忆是杭州。山寺月中寻桂子，郡亭枕上看潮头。', tags: ['西湖', '南宋'] },
];

interface CityState {
  cities: City[];
  cityCache: Record<number, City>;
  loadCities: () => Promise<void>;
  loadCity: (id: number) => Promise<City>;
}

export const useCityStore = create<CityState>((set, get) => ({
  cities: [],
  cityCache: {},
  loadCities: async () => {
    // Mock 异步加载
    return new Promise((resolve) => {
      setTimeout(() => {
        set({ cities: MOCK_CITIES });
        resolve();
      }, 500);
    });
  },
  loadCity: async (id: number) => {
    return new Promise((resolve, reject) => {
      const city = MOCK_CITIES.find(c => c.id === id);
      if (city) {
        set((state) => ({ cityCache: { ...state.cityCache, [id]: city } }));
        resolve(city);
      } else {
        reject('City not found');
      }
    });
  },
}));
