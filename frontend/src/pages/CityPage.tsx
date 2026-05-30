import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useCityStore } from '../store/useCityStore';

export const CityPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { loadCity } = useCityStore();
  const [city, setCity] = useState<any>(null);

  useEffect(() => {
    if (id) {
      loadCity(Number(id)).then(setCity).catch(console.error);
    }
  }, [id, loadCity]);

  if (!city) return <div className="p-8 text-center">加载中...</div>;

  return (
    <div className="min-h-screen bg-background p-4 relative">
      <header className="mb-6">
        <div className="glass-panel inline-block px-6 py-2 rounded-full font-bold">
          <span className="text-primary cursor-pointer" onClick={() => navigate(-1)}>&lt; 返回</span>
        </div>
      </header>

      <div className="glass-panel p-8 rounded-2xl max-w-4xl mx-auto space-y-8">
        <div className="text-center">
          <h2 className="text-4xl font-bold text-foreground mb-2">{city.name}</h2>
          <p className="text-muted-foreground">{city.province} · {city.tags?.join(' / ')}</p>
        </div>
        
        <p className="text-lg leading-relaxed">{city.intro}</p>

        <div className="flex gap-4 justify-center mt-8">
          <button 
            onClick={() => navigate(`/city/${city.id}/chat/8`)} // 假装人物id=8
            className="px-6 py-2 bg-primary text-primary-foreground rounded-lg shadow hover:opacity-90"
          >
            与人物对话
          </button>
          <button 
            onClick={() => navigate(`/city/${city.id}/checkin`)}
            className="px-6 py-2 bg-destructive text-destructive-foreground rounded-lg shadow hover:opacity-90"
          >
            赛博打卡
          </button>
          <button 
            onClick={() => navigate('/achievements')}
            className="px-6 py-2 bg-secondary text-secondary-foreground rounded-lg shadow border"
          >
            成就墙
          </button>
        </div>
      </div>
    </div>
  );
};
