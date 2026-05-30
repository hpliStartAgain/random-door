import React, { useEffect, useState } from 'react';
import { useViewStore } from '../../store/useViewStore';
import { useCityStore } from '../../store/useCityStore';

export const CityDrawer: React.FC = () => {
  const { currentView, activeCityId, setView } = useViewStore();
  const { loadCity } = useCityStore();
  const [city, setCity] = useState<any>(null);

  const isOpen = currentView === 'CITY_DETAIL' && activeCityId !== null;

  useEffect(() => {
    if (isOpen && activeCityId) {
      loadCity(activeCityId).then(setCity).catch(console.error);
    }
  }, [isOpen, activeCityId, loadCity]);

  return (
    <div 
      className={`absolute top-0 right-0 h-full w-[400px] z-20 p-4 transition-transform duration-500 ease-in-out ${isOpen ? 'translate-x-0' : 'translate-x-[110%]'}`}
    >
      <div className="glass-panel h-full w-full rounded-2xl shadow-2xl flex flex-col overflow-hidden">
        {/* 顶部控制栏 */}
        <div className="p-4 border-b border-border/50 flex justify-between items-center">
          <span className="font-bold">城市档案</span>
          <button 
            onClick={() => setView('FREE_EXPLORE')}
            className="w-8 h-8 flex items-center justify-center rounded-full bg-muted text-muted-foreground hover:text-foreground"
          >
            ✕
          </button>
        </div>

        {/* 内容区 */}
        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {city ? (
            <>
              <div className="text-center">
                <h2 className="text-3xl font-bold text-primary mb-1">{city.name}</h2>
                <p className="text-sm text-muted-foreground">{city.province} · {city.tags?.join(' / ')}</p>
              </div>
              <p className="text-sm leading-relaxed text-foreground/80">{city.intro}</p>
              
              <div className="pt-4 border-t border-border/50 space-y-4">
                <button className="w-full py-2 bg-primary/10 text-primary font-semibold rounded-lg hover:bg-primary/20 transition-colors">
                  与代表人物对话 (Mock)
                </button>
                <button className="w-full py-2 bg-destructive/10 text-destructive font-semibold rounded-lg hover:bg-destructive/20 transition-colors">
                  📸 生成赛博打卡
                </button>
              </div>
            </>
          ) : (
            <div className="text-center text-muted-foreground py-10">加载档案中...</div>
          )}
        </div>
      </div>
    </div>
  );
};
