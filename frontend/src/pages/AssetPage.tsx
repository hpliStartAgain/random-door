import React, { useEffect, useState } from 'react';
import { api } from '../api';
import type { UserAssetsResponse } from '../api/types';
import { useUserStore } from '../store/useUserStore';
import { useViewStore } from '../store/useViewStore';
import { foxImages } from '../assets/foxImages';

export const AssetPage: React.FC = () => {
  const { setView } = useViewStore();
  const { userId } = useUserStore();
  const [data, setData] = useState<UserAssetsResponse | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!userId) return;
    setLoading(true);
    api.getUserAssets(userId)
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [userId]);

  return (
    <div className="min-h-screen bg-background p-4 max-w-5xl mx-auto pt-24">
      <header className="mb-6">
        <div className="glass-panel inline-flex items-center px-6 py-2 rounded-full font-bold shadow-sm border border-border">
          <button className="text-primary hover:underline" onClick={() => setView('HOME')}>&lt; 返回</button>
          <span className="mx-4">我的资产</span>
        </div>
      </header>

      {/* 页头主视觉 */}
      <div className="flex flex-col items-center mb-10 text-center">
        <img
          src={foxImages.postcardBoard}
          alt="旅行记忆拼贴板"
          className="w-48 h-36 object-contain fox-float drop-shadow-md mb-4"
          onError={(e) => { (e.currentTarget as HTMLImageElement).style.display = 'none'; }}
        />
        <h2 className="font-serif-display text-2xl font-black text-primary mb-1">我的旅行记忆</h2>
        <p className="text-sm text-muted-foreground">每一座城市，都会变成一张明信片</p>
      </div>

      {loading && <div className="text-sm text-muted-foreground text-center">正在加载资产...</div>}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <section className="glass-panel p-5 rounded-2xl border border-border">
          <h2 className="font-bold text-base mb-3">走过城市</h2>
          <div className="space-y-2 max-h-[520px] overflow-y-auto">
            {data?.visited_cities.length ? data.visited_cities.map(city => (
              <div key={`${city.id}-${city.visited_at}`} className="p-3 rounded-xl border border-border bg-background/70">
                <div className="font-semibold text-sm">{city.name}</div>
                <div className="text-xs text-muted-foreground">{city.province} · {new Date(city.visited_at).toLocaleDateString()}</div>
              </div>
            )) : (
              <div className="flex flex-col items-center py-8 gap-3">
                <img
                  src={foxImages.postcardBoard}
                  alt="暂无记录"
                  className="w-28 h-20 object-contain opacity-60"
                  onError={(e) => { (e.currentTarget as HTMLImageElement).style.display = 'none'; }}
                />
                <p className="text-sm text-muted-foreground text-center">
                  还没有城市收藏<br/>
                  <span className="text-xs">先开启一次任意门吧</span>
                </p>
              </div>
            )}
          </div>
        </section>

        <section className="glass-panel p-5 rounded-2xl border border-border lg:col-span-2">
          <h2 className="font-bold text-base mb-3">打卡海报</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            {data?.posters.length ? data.posters.map(poster => (
              <article key={poster.checkin_id} className="rounded-xl overflow-hidden border border-border bg-background/70">
                <img src={poster.generated_image_url} alt={`${poster.city_name}打卡海报`} className="w-full aspect-[4/3] object-cover" />
                <div className="p-3">
                  <div className="font-semibold text-sm">{poster.city_name}</div>
                  <div className="text-xs text-muted-foreground">
                    {poster.landmark_name || '城市打卡'} · {new Date(poster.created_at).toLocaleDateString()}
                  </div>
                </div>
              </article>
            )) : <div className="text-sm text-muted-foreground">还没有生成海报</div>}
          </div>
        </section>

        <section className="glass-panel p-5 rounded-2xl border border-border lg:col-span-3">
          <h2 className="font-bold text-base mb-3">成就进度</h2>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {data?.achievement_progress.length ? data.achievement_progress.map(item => (
              <div key={item.code} className="p-3 rounded-xl border border-border bg-background/70">
                <div className="flex justify-between text-xs mb-2">
                  <span className="font-semibold">{item.code}</span>
                  <span className="text-muted-foreground">{item.current} / {item.target}</span>
                </div>
                <div className="h-2 bg-border rounded-full overflow-hidden">
                  <div className="h-full bg-primary rounded-full" style={{ width: `${Math.min(100, (item.current / item.target) * 100)}%` }} />
                </div>
              </div>
            )) : <div className="text-sm text-muted-foreground">暂无进行中的成就</div>}
          </div>
        </section>
      </div>
    </div>
  );
};
