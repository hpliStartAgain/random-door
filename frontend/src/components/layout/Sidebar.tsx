import React, { useState, useEffect } from 'react';
import { useCityStore } from '../../store/useCityStore';
import { useGameStore } from '../../store/useGameStore';
import { useMapStore } from '../../store/useMapStore';
import { useUserStore } from '../../store/useUserStore';
import { CityDetailPanel } from '../CityDetailPanel';

export const Sidebar: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'explore' | 'dice'>('explore');
  const [selectedCity, setSelectedCity] = useState<any | null>(null);
  
  const { cities, loadCities, loadCity } = useCityStore();
  const { roll, rolling, targetCity } = useGameStore();
  const { flyTo, mapInstance } = useMapStore();
  const { userId, currentCityId } = useUserStore();

  useEffect(() => {
    loadCities();
  }, [loadCities]);

  const handleRoll = async () => {
    if (!userId) return;
    
    // Get center from map if available
    let lat = 39.9042;
    let lng = 116.4074;
    if (mapInstance) {
      const center = mapInstance.getCenter();
      lat = center.lat;
      lng = center.lng;
    }
    
    try {
      const res = await roll(userId, currentCityId || 1, lat, lng);
      if (res && res.target_city) {
        flyTo(res.target_city.lng, res.target_city.lat);
        const detail = await loadCity(res.target_city.id);
        setSelectedCity(detail);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const handleCityClick = async (city: any) => {
    flyTo(city.lng, city.lat);
    setSelectedCity(city);
    try {
      const detail = await loadCity(city.id);
      setSelectedCity(detail);
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <aside className="w-[380px] h-full sidebar-container flex flex-col z-20 overflow-hidden relative border-r border-border/50 shadow-sm shrink-0 bg-background">
      {/* 城市详情面板（通过绝对定位覆盖在上方，滑动出入） */}
      <CityDetailPanel city={selectedCity} onBack={() => setSelectedCity(null)} />
      
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
          onClick={() => setActiveTab('dice')}
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
            placeholder="搜索地标、朝代、风云人物..." 
            className="w-full px-4 py-2 bg-card border border-border/60 rounded-xl text-sm outline-none focus:border-primary/50 transition-colors"
          />
        </div>
      )}

      {/* 滚动内容区 */}
      <div className="flex-1 overflow-y-auto px-6 pb-8 space-y-4">
        {activeTab === 'explore' ? (
          cities.length > 0 ? cities.map((city) => (
            <div 
              key={city.id} 
              onClick={() => handleCityClick(city)}
              className="group relative h-48 rounded-2xl overflow-hidden cursor-pointer shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="absolute inset-0 bg-gradient-to-br from-primary/80 to-accent/80 mix-blend-multiply transition-transform duration-700 group-hover:scale-105" />
              <div className="absolute inset-0 bg-black/20 group-hover:bg-black/10 transition-colors" />
              
              <div className="absolute top-3 right-3 bg-white/90 backdrop-blur text-primary text-xs font-bold px-3 py-1 rounded-full flex items-center gap-1 shadow-sm">
                ✓ 必看地标
              </div>

              <div className="absolute bottom-4 left-4 right-4">
                <h3 className="text-white text-xl font-bold leading-tight">{city.name}</h3>
                <p className="text-white/90 text-sm mt-1">{city.province}</p>
                <div className="flex gap-2 mt-2">
                  {city.tags?.map(tag => (
                    <span key={tag} className="text-[10px] px-2 py-0.5 bg-black/30 text-white rounded-full backdrop-blur-sm border border-white/20">{tag}</span>
                  ))}
                </div>
              </div>
            </div>
          )) : (
            <div className="text-center text-muted-foreground py-10">正在加载地标档案...</div>
          )
        ) : (
          <div className="h-full flex flex-col items-center justify-center space-y-8 p-4">
            <div className="w-32 h-32 rounded-full bg-secondary border border-border flex items-center justify-center shadow-inner relative">
               {rolling ? (
                 <span className="text-5xl animate-spin">🎲</span>
               ) : targetCity ? (
                 <span className="text-3xl font-bold text-primary">{targetCity.name}</span>
               ) : (
                 <span className="text-5xl opacity-50">🎲</span>
               )}
            </div>
            
            <button 
              onClick={handleRoll}
              disabled={rolling}
              className="w-full py-3 bg-primary text-primary-foreground font-bold rounded-xl shadow-sm hover:opacity-90 disabled:opacity-50 transition-opacity"
            >
              {rolling ? '抛掷中...' : '抛掷骰子'}
            </button>
            
            {targetCity && !rolling && (
              <p className="text-sm text-center text-muted-foreground">
                目标确认！飞往 {targetCity.name}。
              </p>
            )}
          </div>
        )}
      </div>
    </aside>
  );
};
