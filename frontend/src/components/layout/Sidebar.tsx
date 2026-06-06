import React, { useState, useEffect } from 'react';
import { Dices } from 'lucide-react';
import { api } from '../../api';
import { useCityStore } from '../../store/useCityStore';
import { useGameStore } from '../../store/useGameStore';
import { useMapStore } from '../../store/useMapStore';
import { useUserStore } from '../../store/useUserStore';
import { useViewStore } from '../../store/useViewStore';
import { CityDetailPanel } from '../CityDetailPanel';

function SidebarFox() {
  return (
    <svg viewBox="0 0 96 96" className="h-20 w-20 drop-shadow-lg">
      <path d="M18 28 L33 15 L37 35 Z" fill="#D47A3C" />
      <path d="M78 28 L63 15 L59 35 Z" fill="#D47A3C" />
      <path d="M24 31 C28 15 68 15 72 31 C84 48 75 78 48 82 C21 78 12 48 24 31Z" fill="#E48743" />
      <path d="M30 32 C37 43 42 55 48 80 C54 55 59 43 66 32 C60 68 36 68 30 32Z" fill="#FFF2D7" opacity="0.96" />
      <path d="M20 31 L32 20 L35 34 Z" fill="#8D3F28" opacity="0.55" />
      <path d="M76 31 L64 20 L61 34 Z" fill="#8D3F28" opacity="0.55" />
      <circle cx="36" cy="47" r="3.8" fill="#22302C" />
      <circle cx="60" cy="47" r="3.8" fill="#22302C" />
      <path d="M45 58 Q48 61 51 58" fill="none" stroke="#22302C" strokeWidth="2.4" strokeLinecap="round" />
      <path d="M48 54 L43 59 L53 59 Z" fill="#22302C" />
      <path d="M17 65 C7 64 7 48 21 46 C34 44 39 57 28 65 C35 67 42 72 48 80 C35 79 24 74 17 65Z" fill="#D87538" />
      <path d="M18 63 C12 60 14 52 22 52 C28 52 30 58 24 62 C31 64 37 70 41 75 C31 73 23 69 18 63Z" fill="#FFF2D7" opacity="0.92" />
    </svg>
  );
}

export const Sidebar: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'explore' | 'dice'>('explore');
  const [selectedCity, setSelectedCity] = useState<any | null>(null);
  const [selectedVisitId, setSelectedVisitId] = useState<number | undefined>(undefined);
  const [searchQuery, setSearchQuery] = useState('');
  
  const { cities, cityCache, loadCities, loadCity } = useCityStore();
  const { lastRoll } = useGameStore();
  const { flyTo } = useMapStore();
  const { userId, setCurrentCityId } = useUserStore();
  const { currentView, activeCityId, setView } = useViewStore();

  useEffect(() => {
    loadCities();
  }, [loadCities]);

  useEffect(() => {
    if (currentView !== 'CITY_DETAIL' || !activeCityId) return;
    let cancelled = false;
    loadCity(activeCityId)
      .then((detail) => {
        if (cancelled) return;
        setSelectedCity(detail);
        setCurrentCityId(activeCityId);
        if (lastRoll?.target_city.id === activeCityId) {
          setSelectedVisitId(lastRoll.visit_id);
        } else {
          setSelectedVisitId(undefined);
        }
      })
      .catch(console.error);
    return () => { cancelled = true; };
  }, [currentView, activeCityId, lastRoll?.target_city.id, lastRoll?.visit_id, loadCity, setCurrentCityId]);

  useEffect(() => {
    if (currentView !== 'CITY_DETAIL' || !activeCityId) return;
    const cached = cityCache[activeCityId];
    if (cached) setSelectedCity(cached);
  }, [currentView, activeCityId, cityCache]);

  const closeCityDetail = () => {
    setSelectedCity(null);
    setSelectedVisitId(undefined);
    if (currentView === 'CITY_DETAIL') {
      setView('HOME');
    }
  };

  const handleCityClick = async (city: any) => {
    flyTo(city.lng, city.lat);
    setSelectedCity(city);
    setSelectedVisitId(undefined);
    setCurrentCityId(city.id);
    setView('CITY_DETAIL', city.id);
    try {
      if (userId) {
        const visit = await api.createFreeVisit(userId, city.id);
        setSelectedVisitId(visit.visit_id);
      }
      const detail = await loadCity(city.id);
      setSelectedCity(detail);
    } catch (e) {
      console.error(e);
    }
  };

  const filteredCities = searchQuery
    ? cities.filter(c =>
        c.name.includes(searchQuery) ||
        c.province.includes(searchQuery) ||
        c.tags?.some(t => t.includes(searchQuery))
      )
    : cities;

  return (
    <aside className="w-[380px] h-full sidebar-container flex flex-col z-20 overflow-hidden relative border-r border-border/50 shadow-sm shrink-0 bg-background">
      {/* 城市详情面板（通过绝对定位覆盖在上方，滑动出入） */}
      <CityDetailPanel city={selectedCity} visitId={selectedVisitId} onBack={closeCityDetail} />
      
      {/* 头部标题区 */}
      <div className="p-6 pb-2">
        <h2 className="text-2xl font-bold text-primary tracking-tight mb-2">
          {activeTab === 'explore' ? '探索风物' : '听天由命'}
        </h2>
        <p className="text-sm text-muted-foreground">
          {activeTab === 'explore' 
            ? `${cities.length || 0} 处名胜 · 穿越历史的长河`
            : '掷出命运的骰子，让风带你去往未知的远方'}
        </p>
      </div>

      {/* Tabs */}
      <div className="px-6 py-4 flex gap-2">
        <button 
          onClick={() => setActiveTab('explore')}
          className={`px-4 py-1.5 rounded-full text-sm font-medium border transition-colors ${activeTab === 'explore' ? 'bg-primary text-primary-foreground border-primary' : 'bg-transparent text-foreground border-border hover:bg-secondary'}`}
        >
          自由探索
        </button>
        <button 
          onClick={() => { setActiveTab('dice'); setView('GAME_DICE'); }}
          className={`px-4 py-1.5 rounded-full text-sm font-medium border transition-colors ${activeTab === 'dice' ? 'bg-primary text-primary-foreground border-primary' : 'bg-transparent text-foreground border-border hover:bg-secondary'}`}
        >
          随机漫游
        </button>
      </div>

      {/* 搜索框 */}
      {activeTab === 'explore' && (
        <div className="px-6 pb-4">
          <input 
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="搜索城市、省份、标签..."
            className="w-full px-4 py-2 bg-card border border-border/60 rounded-xl text-sm outline-none focus:border-primary/50 transition-colors"
          />
        </div>
      )}

      {/* 滚动内容区 */}
      <div className="flex-1 overflow-y-auto px-6 pb-8 space-y-4">
        {activeTab === 'explore' ? (
          filteredCities.length > 0 ? filteredCities.map((city) => (
            <div
              key={city.id}
              onClick={() => handleCityClick(city)}
              className="group relative h-48 rounded-2xl overflow-hidden cursor-pointer shadow-sm hover:shadow-md transition-shadow"
            >
              {city.cover_image_url ? (
                <img
                  src={city.cover_image_url}
                  alt={city.name}
                  className="absolute inset-0 w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
                />
              ) : (
                <div className="absolute inset-0 bg-gradient-to-br from-primary/80 to-accent/80 transition-transform duration-700 group-hover:scale-105" />
              )}
              <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent group-hover:from-black/60 transition-colors" />

              <div className="absolute top-3 right-3 bg-white/90 backdrop-blur text-primary text-xs font-bold px-3 py-1 rounded-full flex items-center gap-1 shadow-sm">
                ✓ 必看地标
              </div>

              <div className="absolute bottom-4 left-4 right-4">
                <h3 className="text-white text-xl font-bold leading-tight">{city.name}</h3>
                <p className="text-white/80 text-sm mt-0.5">{city.province}</p>
                <div className="flex gap-2 mt-2 flex-wrap">
                  {city.tags?.map(tag => (
                    <span key={tag} className="text-[10px] px-2 py-0.5 bg-black/30 text-white rounded-full backdrop-blur-sm border border-white/20">{tag}</span>
                  ))}
                </div>
              </div>
            </div>
          )) : (
            <div className="text-center text-muted-foreground py-10">
              {searchQuery ? `未找到"${searchQuery}"相关城市` : '正在加载地标档案...'}
            </div>
          )
        ) : (
          <div className="h-full flex flex-col items-center justify-center space-y-8 p-4">
            <div className="w-32 h-32 rounded-full bg-secondary border border-border flex items-center justify-center shadow-inner relative">
              <SidebarFox />
            </div>
            
            <button 
              onClick={() => setView('GAME_DICE')}
              className="w-full py-3 bg-primary text-primary-foreground font-bold rounded-xl shadow-sm hover:opacity-90 transition-opacity flex items-center justify-center gap-2"
            >
              <Dices className="h-4 w-4" />
              打开命运骰台
            </button>
            
            {lastRoll && (
              <p className="text-sm text-center text-muted-foreground">
                上次目的地：{lastRoll.target_city.name} · {lastRoll.direction} {lastRoll.distance_km}km
              </p>
            )}
          </div>
        )}
      </div>
    </aside>
  );
};
