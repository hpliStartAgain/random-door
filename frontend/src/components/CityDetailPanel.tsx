import React, { useEffect, useState } from 'react';
import type { CityDetail } from '../api/types';
import { useMapStore } from '../store/useMapStore';
import { useViewStore } from '../store/useViewStore';
import { useUserStore } from '../store/useUserStore';
import { api } from '../api';

interface Props {
  city: CityDetail | null;
  onBack: () => void;
}

export const CityDetailPanel: React.FC<Props> = ({ city, onBack }) => {
  const { resetView } = useMapStore();
  const { openDrawer, setCanvasMode } = useViewStore();
  const { userId } = useUserStore();
  const [mounted, setMounted] = useState(false);
  const [checkingIn, setCheckingIn] = useState(false);
  const [checkinResult, setCheckinResult] = useState<any>(null);

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
    setCanvasMode('map');
    setCheckinResult(null);
    setTimeout(() => {
      onBack();
    }, 300);
  };

  const handleCheckin = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || e.target.files.length === 0 || !userId || !city) return;
    const file = e.target.files[0];
    
    setCheckingIn(true);
    setCheckinResult(null);
    try {
      const formData = new FormData();
      formData.append('selfie_file', file);
      formData.append('user_id', userId.toString());
      formData.append('city_id', city.id.toString());
      // Find the first landmark to use as reference if any
      const firstLandmarkId = (city as any).landmarks?.[0]?.id || 1; // Default to 1 if no landmarks are passed in the detail (should be passed by API)
      formData.append('landmark_id', firstLandmarkId.toString());
      
      const imgRes = await api.generateImage(formData);
      const chkRes = await api.createCheckin(userId, city.id, firstLandmarkId, undefined, imgRes.generated_image_url);
      
      setCheckinResult({
        imageUrl: imgRes.generated_image_url,
        achievements: chkRes.unlocked_achievements
      });
    } catch (err) {
      console.error(err);
      alert('打卡失败，请重试');
    } finally {
      setCheckingIn(false);
    }
  };

  if (!city) return null;

  const mediaUrl = (url?: string) => {
    if (!url) return '';
    // 绝对 URL 原样返回；相对路径（/uploads、/static）交给当前域名
    // 开发环境由 Vite 代理转发，生产环境由 Caddy 托管。
    return url;
  };

  return (
    <div 
      className={`absolute inset-0 bg-background flex flex-col z-30 transition-transform duration-300 ease-in-out ${mounted ? 'translate-x-0' : 'translate-x-full'}`}
    >
      {/* 顶部英雄头图 */}
      <div className="relative shrink-0">
        {city.cover_image_url ? (
          <img src={mediaUrl(city.cover_image_url)} alt={city.name} className="w-full h-40 object-cover" />
        ) : (
          <div className="w-full h-40 bg-gradient-to-br from-primary/30 to-accent/40" />
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-black/30" />
        <button
          onClick={handleBack}
          className="absolute top-4 left-4 px-3 py-1.5 rounded-full bg-white/90 backdrop-blur text-xs font-semibold text-foreground hover:bg-white transition-colors flex items-center gap-1 shadow-sm"
        >
          <span>&larr;</span> 返回大地图
        </button>
        <div className="absolute bottom-4 left-5 right-5">
          <h2 className="font-serif-display text-3xl font-bold text-white tracking-tight drop-shadow-lg">{city.name}</h2>
          <p className="text-sm text-white/90 mt-0.5">{city.province}</p>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-5 py-6 space-y-6 pb-12">
        {/* 城市简介 */}
        {city.intro && (
          <p className="text-sm leading-relaxed text-foreground/80">{city.intro}</p>
        )}

        {/* 标签 */}
        {city.tags && city.tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {city.tags.map(tag => (
              <span key={tag} className="pill bg-accent/15 text-accent">{tag}</span>
            ))}
          </div>
        )}

        {/* 地标巡礼 */}
        {city.landmarks && city.landmarks.length > 0 && (
          <section className="panel">
            <div className="panel-header"><span>🏛️</span> 地标巡礼</div>
            <div className="panel-body space-y-3">
              {city.landmarks.map(lm => (
                <div
                  key={lm.id}
                  className="flex gap-3 rounded-xl border border-border/60 overflow-hidden hover:border-accent/50 hover:shadow-sm transition-all"
                >
                  {lm.image_url ? (
                    <img src={mediaUrl(lm.image_url)} alt={lm.name} className="w-24 h-24 object-cover shrink-0" />
                  ) : (
                    <div className="w-24 h-24 shrink-0 bg-accent/10 flex items-center justify-center text-2xl">🏯</div>
                  )}
                  <div className="py-2.5 pr-3 min-w-0">
                    <h4 className="font-serif-display font-semibold text-base text-foreground mb-1">{lm.name}</h4>
                    <p className="text-xs text-muted-foreground leading-relaxed line-clamp-3">{lm.description}</p>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* 历史人物 */}
        {city.characters && city.characters.length > 0 && (
          <section className="panel">
            <div className="panel-header"><span>👤</span> 历史人物</div>
            <div className="panel-body space-y-3">
              {city.characters.map(char => (
                <div
                  key={char.id}
                  onClick={() => openDrawer('chat', char)}
                  className="flex items-center gap-3 rounded-xl border border-border/60 p-3 hover:border-primary/50 hover:shadow-sm transition-all cursor-pointer group"
                >
                  <div className="w-12 h-12 rounded-full shrink-0 bg-gradient-to-br from-primary/20 to-accent/30 flex items-center justify-center overflow-hidden">
                    {char.avatar_url ? (
                      <img src={mediaUrl(char.avatar_url)} alt={char.name} className="w-full h-full object-cover" />
                    ) : (
                      <span className="text-xl">🎎</span>
                    )}
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <h4 className="font-serif-display font-semibold text-base">{char.name}</h4>
                      <span className="pill bg-primary/10 text-primary">{char.character_type === 'culture' ? '文化符号' : char.character_type}</span>
                    </div>
                    <p className="text-xs text-muted-foreground line-clamp-2 mt-0.5">{char.dialect_style}</p>
                  </div>
                  <span className="text-muted-foreground group-hover:text-primary transition-colors text-sm">&rarr;</span>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* 风味人间 */}
        {city.foods && city.foods.length > 0 && (
          <section className="panel">
            <div className="panel-header"><span>🍜</span> 风味人间</div>
            <div className="panel-body grid grid-cols-2 gap-3">
              {city.foods.map(food => (
                <div
                  key={food.id}
                  onClick={() => openDrawer('gallery', food)}
                  className="rounded-xl border border-border/60 overflow-hidden hover:border-accent/50 hover:shadow-sm transition-all cursor-pointer"
                >
                  {food.image_url ? (
                    <img src={mediaUrl(food.image_url)} alt={food.name} className="w-full h-20 object-cover" />
                  ) : (
                    <div className="w-full h-20 bg-accent/10 flex items-center justify-center text-2xl">🥢</div>
                  )}
                  <div className="p-2.5">
                    <h4 className="font-semibold text-sm mb-0.5">{food.name}</h4>
                    <p className="text-xs text-muted-foreground line-clamp-2">{food.description}</p>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* 沉浸体验 */}
        <section className="panel">
          <div className="panel-header"><span>✨</span> 沉浸体验</div>
          <div className="panel-body space-y-3">
            <div
              onClick={() => setCanvasMode('street', city)}
              className="w-full h-28 rounded-xl bg-gradient-to-br from-primary/10 to-accent/15 border border-primary/20 flex flex-col items-center justify-center cursor-pointer hover:scale-[1.01] hover:shadow-sm transition-all group"
            >
              <span className="text-3xl mb-1.5 group-hover:scale-110 transition-transform">🛣️</span>
              <span className="text-sm font-semibold text-primary">进入 3D 街景</span>
            </div>

            <label className="w-full py-3 bg-accent/15 text-accent font-semibold rounded-xl hover:bg-accent/25 transition-colors flex items-center justify-center cursor-pointer relative overflow-hidden">
              {checkingIn ? (
                <span className="animate-pulse">正在穿梭时空生成大片...</span>
              ) : (
                <>
                  <span className="mr-2">📸</span> 生成赛博打卡
                  <input
                    type="file"
                    accept="image/*"
                    capture="user"
                    className="absolute inset-0 opacity-0 cursor-pointer"
                    onChange={handleCheckin}
                    disabled={checkingIn}
                  />
                </>
              )}
            </label>

            {checkinResult && (
              <div className="p-4 rounded-2xl bg-card border border-border shadow-sm animate-in fade-in slide-in-from-bottom-4">
                <h4 className="font-semibold text-primary mb-2">打卡成功！</h4>
                <img
                  src={mediaUrl(checkinResult.imageUrl)}
                  alt="Cyber Checkin"
                  className="w-full rounded-xl object-cover mb-4 shadow-sm"
                />
                {checkinResult.achievements?.length > 0 && (
                  <div className="space-y-2">
                    <p className="text-xs text-muted-foreground font-semibold">🎉 解锁新成就</p>
                    {checkinResult.achievements.map((ach: any) => (
                      <div key={ach.code} className="flex items-center gap-2 bg-primary/5 p-2 rounded-lg border border-primary/10">
                        <span className="text-2xl">🏅</span>
                        <div>
                          <p className="text-sm font-semibold text-foreground">{ach.name}</p>
                          <p className="text-xs text-muted-foreground">{ach.description}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        </section>
      </div>
    </div>
  );
};
