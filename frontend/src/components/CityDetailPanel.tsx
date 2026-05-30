import React, { useEffect, useState } from 'react';
import { City } from '../store/useCityStore';
import { useMapStore } from '../store/useMapStore';
import { useViewStore } from '../store/useViewStore';

interface Props {
  city: City | null;
  onBack: () => void;
}

export const CityDetailPanel: React.FC<Props> = ({ city, onBack }) => {
  const { resetView } = useMapStore();
  const { openDrawer, setCanvasMode } = useViewStore();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    if (city) {
      requestAnimationFrame(() => setMounted(true));
    } else {
      setMounted(false);
    }
  }, [city]);

  const handleBack = () => {
    setMounted(false);
    resetView();
    // 退出时也确保退出街景模式
    setCanvasMode('map');
    setTimeout(() => {
      onBack();
    }, 300); // 等待退场动画完成
  };

  if (!city) return null;

  return (
    <div 
      className={`absolute inset-0 bg-background flex flex-col z-30 transition-transform duration-300 ease-in-out ${mounted ? 'translate-x-0' : 'translate-x-full'}`}
    >
      <div className="p-6 pb-4 border-b border-border/50 sticky top-0 bg-background/95 backdrop-blur z-10 shadow-sm">
        <button 
          onClick={handleBack}
          className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors flex items-center gap-1 mb-4"
        >
          <span>&larr;</span> 返回大地图
        </button>
        <h2 className="text-3xl font-bold text-primary tracking-tight mb-2">{city.name}</h2>
        <p className="text-sm text-foreground/80">{city.province} · {city.intro}</p>
      </div>

      <div className="flex-1 overflow-y-auto p-6 space-y-8 pb-12">
        {/* 风云人物 */}
        {city.figures && city.figures.length > 0 && (
          <section>
            <h3 className="text-lg font-bold text-foreground mb-3 flex items-center gap-2">
              <span className="w-1.5 h-4 bg-primary rounded-full"></span> 历史风云
            </h3>
            <div className="flex gap-4 overflow-x-auto pb-4 snap-x hide-scrollbar">
              {city.figures.map(fig => (
                <div 
                  key={fig.name} 
                  onClick={() => openDrawer('chat', fig)}
                  className="flex-none w-64 bg-card rounded-2xl p-4 shadow-sm border border-border snap-start hover:border-primary/40 hover:shadow-md transition-all cursor-pointer"
                >
                  <div className="text-xs text-primary font-bold bg-primary/10 w-fit px-2 py-0.5 rounded-full mb-2">{fig.dynasty}</div>
                  <h4 className="text-lg font-bold mb-1">{fig.name}</h4>
                  <p className="text-sm text-muted-foreground line-clamp-3">{fig.desc}</p>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* 地方美食 */}
        {city.foods && city.foods.length > 0 && (
          <section>
            <h3 className="text-lg font-bold text-foreground mb-3 flex items-center gap-2">
              <span className="w-1.5 h-4 bg-primary rounded-full"></span> 风味人间
            </h3>
            <div className="grid grid-cols-2 gap-3">
              {city.foods.map(food => (
                <div 
                  key={food.name} 
                  onClick={() => openDrawer('gallery', food)}
                  className="bg-card rounded-xl p-3 shadow-sm border border-border hover:border-primary/40 hover:shadow-md transition-all cursor-pointer"
                >
                  <h4 className="font-bold text-sm mb-1">{food.name}</h4>
                  <p className="text-xs text-muted-foreground line-clamp-3">{food.desc}</p>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* 3D街景入口 */}
        <section>
          <div 
            onClick={() => setCanvasMode('street', city)}
            className="w-full h-32 rounded-2xl bg-gradient-to-br from-primary/10 to-accent/10 border border-primary/20 flex flex-col items-center justify-center cursor-pointer hover:bg-primary/20 hover:scale-[1.02] shadow-sm transition-all group"
          >
            <span className="text-4xl mb-2 group-hover:scale-110 transition-transform">🛣️</span>
            <span className="text-sm font-bold text-primary">进入 3D 街景</span>
          </div>
        </section>
      </div>
    </div>
  );
};
