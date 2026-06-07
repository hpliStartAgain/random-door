import React, { useState, useEffect } from 'react';
import { Sparkles } from 'lucide-react';
import { api } from '../../api';
import { foxImages } from '../../assets/foxImages';
import { useCityStore } from '../../store/useCityStore';
import { useGameStore } from '../../store/useGameStore';
import { useMapStore } from '../../store/useMapStore';
import { useUserStore } from '../../store/useUserStore';
import { useViewStore } from '../../store/useViewStore';
import { CityDetailPanel } from '../CityDetailPanel';
import { AchievementUnlock } from '../overlays/AchievementUnlock';
import type { Achievement } from '../overlays/AchievementUnlock';
import { getRegionOptions, getAllTags } from '../../lib/cityFilters';

function SidebarFoxDoor() {
  return (
    <img
      src={foxImages.magicDoor}
      alt="开启任意门"
      className="h-28 w-28 object-contain fox-float drop-shadow-lg"
      onError={(e) => { (e.currentTarget as HTMLImageElement).style.opacity = '0'; }}
    />
  );
}

export const Sidebar: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'explore' | 'dice'>('explore');
  const [selectedCity, setSelectedCity] = useState<any | null>(null);
  const [selectedVisitId, setSelectedVisitId] = useState<number | undefined>(undefined);
  const [unlockedAchievements, setUnlockedAchievements] = useState<Achievement[]>([]);
  
  const {
    cities,
    cityCache,
    loadCities,
    loadCity,
    searchQuery,
    setSearchQuery,
    filteredCities,
    activeRegion,
    activeTag,
    setActiveRegion,
    setActiveTag,
    resetFilters,
  } = useCityStore();
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
    setSelectedCity((prev: any | null) => (
      prev?.id === activeCityId ? prev : cities.find((city) => city.id === activeCityId) ?? null
    ));
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
  }, [currentView, activeCityId, cities, lastRoll?.target_city.id, lastRoll?.visit_id, loadCity, setCurrentCityId]);

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
        if (visit.unlocked_achievements?.length) {
          setUnlockedAchievements(visit.unlocked_achievements);
        }
      }
      const detail = await loadCity(city.id);
      setSelectedCity(detail);
    } catch (e) {
      console.error(e);
    }
  };

  const displayCities = filteredCities();
  const regionOptions = getRegionOptions(cities);
  const tagOptions = getAllTags(cities).slice(0, 8);
  const hasActiveFilter = !!(searchQuery || activeRegion || activeTag);

  const totalLandmarks = displayCities.reduce((s, c) => s + (c.landmark_count ?? 0), 0);
  const totalFoods = displayCities.reduce((s, c) => s + (c.food_count ?? 0), 0);
  const totalChars = displayCities.reduce((s, c) => s + (c.character_count ?? 0), 0);

  return (
    <aside className="w-[380px] h-full sidebar-container flex flex-col z-20 overflow-hidden relative border-r border-border/50 shadow-sm shrink-0 bg-background">
      {/* 城市详情面板（通过绝对定位覆盖在上方，滑动出入） */}
      <CityDetailPanel city={selectedCity} visitId={selectedVisitId} onBack={closeCityDetail} />
      {unlockedAchievements.length > 0 && (
        <AchievementUnlock
          achievements={unlockedAchievements}
          onClose={() => setUnlockedAchievements([])}
        />
      )}
      
      {/* 头部标题区 */}
      <div className="p-6 pb-2">
        <h2 className="text-2xl font-bold text-primary tracking-tight mb-2">
          {activeTab === 'explore' ? '探索风物' : '任意门'}
        </h2>
        <p className="text-sm text-muted-foreground">
          {activeTab === 'explore'
            ? (
              <>
                {cities.length} 座城市 · 筛选中 {displayCities.length} 个
                {Object.keys(cityCache).length > 0 && (
                  <span className="ml-2 text-primary/70">· 已探索 {Object.keys(cityCache).length} 座</span>
                )}
                {(totalLandmarks > 0 || totalFoods > 0 || totalChars > 0) && (
                  <span className="block mt-0.5 text-xs text-muted-foreground/70">
                    {totalLandmarks > 0 && `${totalLandmarks} 地标`}
                    {totalFoods > 0 && ` · ${totalFoods} 美食`}
                    {totalChars > 0 && ` · ${totalChars} 人物`}
                  </span>
                )}
              </>
            )
            : '让狐狸为你开启任意门，寻找下一座城市'}
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
          任意门
        </button>
      </div>

      {/* 搜索框 + 筛选 chips */}
      {activeTab === 'explore' && (
        <div className="px-6 pb-3 space-y-2">
          <input 
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            aria-label="搜索城市"
            placeholder="搜索城市、省份、标签..."
            className="w-full px-4 py-2 bg-card border border-border/60 rounded-xl text-sm outline-none focus:border-primary/50 transition-colors"
          />

          {/* 区域 chips */}
          {regionOptions.length > 0 && (
            <div className="flex gap-1.5 overflow-x-auto pb-1" style={{ scrollbarWidth: 'none' }}>
              {regionOptions.map((region) => (
                <button
                  key={region.key}
                  type="button"
                  onClick={() => setActiveRegion(activeRegion === region.key ? null : region.key)}
                  className={`shrink-0 px-3 py-1 rounded-full text-xs font-medium border transition-colors whitespace-nowrap ${
                    activeRegion === region.key
                      ? 'bg-primary text-primary-foreground border-primary'
                      : 'bg-secondary text-foreground border-border hover:bg-secondary/80'
                  }`}
                >
                  {region.label}
                </button>
              ))}
            </div>
          )}

          {/* Tag chips */}
          {tagOptions.length > 0 && (
            <div className="flex gap-1.5 overflow-x-auto pb-1" style={{ scrollbarWidth: 'none' }}>
              {tagOptions.map((tag) => (
                <button
                  key={tag}
                  type="button"
                  onClick={() => setActiveTag(activeTag === tag ? null : tag)}
                  className={`shrink-0 px-3 py-1 rounded-full text-xs font-medium border transition-colors whitespace-nowrap ${
                    activeTag === tag
                      ? 'bg-primary text-primary-foreground border-primary'
                      : 'bg-secondary text-foreground border-border hover:bg-secondary/80'
                  }`}
                >
                  #{tag}
                </button>
              ))}
              {hasActiveFilter && (
                <button
                  type="button"
                  onClick={resetFilters}
                  className="shrink-0 px-3 py-1 rounded-full text-xs font-medium border border-destructive/50 text-destructive bg-destructive/5 hover:bg-destructive/10 transition-colors whitespace-nowrap"
                >
                  清空
                </button>
              )}
            </div>
          )}
        </div>
      )}

      {/* 滚动内容区 */}
      <div className="flex-1 overflow-y-auto px-6 pb-8 space-y-4">
        {activeTab === 'explore' ? (
          displayCities.length > 0 ? displayCities.map((city) => (
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
            <div className="text-center text-muted-foreground py-10 space-y-3">
              <p className="text-sm">
                {hasActiveFilter
                  ? `未找到符合条件的城市`
                  : '正在加载地标档案...'}
              </p>
              {hasActiveFilter && (
                <div className="flex flex-col items-center gap-2">
                  {searchQuery && (
                    <button
                      type="button"
                      onClick={() => setSearchQuery('')}
                      className="text-xs text-primary underline underline-offset-2"
                    >
                      清空搜索
                    </button>
                  )}
                  <button
                    type="button"
                    onClick={resetFilters}
                    className="text-xs text-primary underline underline-offset-2"
                  >
                    清空所有筛选
                  </button>
                </div>
              )}
            </div>
          )
        ) : (
          <div className="h-full flex flex-col items-center justify-center space-y-8 p-4">
            <SidebarFoxDoor />

            <div className="text-center space-y-1 px-2">
              <div className="text-sm font-semibold text-foreground">听天由命</div>
              <div className="text-xs text-muted-foreground leading-relaxed">
                抽取下一站，让命运为你选择
              </div>
            </div>

            <button
              onClick={() => setView('GAME_DICE')}
              className="w-full py-3.5 text-primary-foreground font-bold rounded-2xl flex items-center justify-center gap-2 hover:brightness-95 transition-all"
              style={{
                background: 'linear-gradient(135deg,hsl(var(--primary)),hsl(var(--accent)))',
                boxShadow: '0 14px 32px rgba(43,58,54,0.18)',
              }}
            >
              <Sparkles className="h-4 w-4" />
              开启任意门
            </button>

            {lastRoll && (
              <p className="text-xs text-center text-muted-foreground">
                上次目的地：{lastRoll.target_city.name} · {lastRoll.direction} {lastRoll.distance_km}km
              </p>
            )}
          </div>
        )}
      </div>
    </aside>
  );
};
