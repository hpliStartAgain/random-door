import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCityStore } from '../store/useCityStore';

export const FreeExplorePage: React.FC = () => {
  const { cities, loadCities } = useCityStore();
  const navigate = useNavigate();

  useEffect(() => {
    loadCities();
  }, [loadCities]);

  return (
    <div className="min-h-screen bg-background relative flex flex-col">
      {/* 假装这是一个全屏地图底板 */}
      <div className="absolute inset-0 bg-[url('https://webapi.amap.com/theme/v1.3/markers/n/mark_b.png')] bg-repeat opacity-5 pointer-events-none" />
      
      <header className="p-4 z-10">
        <div className="glass-panel inline-block px-6 py-2 rounded-full font-bold">
          <span className="text-primary cursor-pointer" onClick={() => navigate('/mode')}>&lt; 返回</span>
          <span className="mx-4 text-foreground">自由探索地图</span>
        </div>
      </header>

      <main className="flex-1 p-4 z-10 flex flex-wrap gap-4 items-start justify-center content-start mt-8">
        {cities.map((city) => (
          <div 
            key={city.id} 
            onClick={() => navigate(`/city/${city.id}`)}
            className="glass-panel p-4 rounded-xl cursor-pointer hover:scale-105 transition-transform text-center w-32"
          >
            <h4 className="font-bold text-lg">{city.name}</h4>
            <p className="text-xs text-muted-foreground">{city.province}</p>
          </div>
        ))}
        {cities.length === 0 && (
          <div className="glass-panel p-8 rounded-xl text-muted-foreground">
            正在加载城市节点...
          </div>
        )}
      </main>
    </div>
  );
};
