import React, { useEffect, useState } from 'react';
import type { CityDetail } from '../api/types';
import { useMapStore } from '../store/useMapStore';
import { useViewStore } from '../store/useViewStore';
import { CheckinFlow } from './CheckinFlow';
import { AchievementUnlock } from './overlays/AchievementUnlock';
import type { Achievement } from './overlays/AchievementUnlock';

interface Props {
  city: CityDetail | null;
  visitId?: number;
  onBack: () => void;
}

export const CityDetailPanel: React.FC<Props> = ({ city, visitId, onBack }) => {
  const { resetView } = useMapStore();
  const { openDrawer, setCanvasMode } = useViewStore();
  const [mounted, setMounted] = useState(false);
  const [showCheckin, setShowCheckin] = useState(false);
  const [unlockedAchievements, setUnlockedAchievements] = useState<Achievement[]>([]);

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
    setShowCheckin(false);
    setTimeout(() => { onBack(); }, 300);
  };

  if (!city) return null;

  const mediaUrl = (url?: string) => url ?? '';

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

        {/* 方言卡 */}
        {(city.dialect_sample || city.dialect_explanation) && (
          <div className="bg-primary/5 rounded-2xl p-4 border border-primary/10">
            <div className="text-[10px] text-muted-foreground uppercase tracking-widest mb-2">方言速记</div>
            {city.dialect_sample && (
              <div className="text-2xl font-bold text-primary tracking-wide">{city.dialect_sample}</div>
            )}
            {city.dialect_explanation && (
              <div className="text-sm text-foreground/70 mt-1.5">{city.dialect_explanation}</div>
            )}
          </div>
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
                  onClick={() => openDrawer('gallery', { ...lm, target_type: 'landmark' })}
                  className="flex gap-3 rounded-xl border border-border/60 overflow-hidden hover:border-accent/50 hover:shadow-sm transition-all cursor-pointer"
                >
                  {lm.image_url ? (
                    <img src={mediaUrl(lm.image_url)} alt={lm.name} className="w-24 h-24 object-cover shrink-0" />
                  ) : (
                    <div className="w-24 h-24 shrink-0 bg-accent/10 flex items-center justify-center text-2xl">🏯</div>
                  )}
                  <div className="py-2.5 pr-3 min-w-0">
                    <h4 className="font-serif-display font-semibold text-base text-foreground mb-1">{lm.name}</h4>
                    <p className="text-xs text-muted-foreground leading-relaxed line-clamp-3">{lm.description}</p>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        setCanvasMode('street', {
                          id: lm.id,
                          city_id: city.id,
                          name: lm.name,
                          province: city.province,
                          lat: city.lat,
                          lng: city.lng,
                          cover_image_url: lm.image_url || city.cover_image_url,
                          tags: city.tags,
                        });
                      }}
                      className="mt-2 text-xs px-2.5 py-1 rounded-lg bg-primary/10 text-primary border border-primary/15 hover:bg-primary/15 transition-colors"
                    >
                      查看图片
                    </button>
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
                <button
                  key={char.id}
                  type="button"
                  onClick={() => openDrawer('chat', char)}
                  className="w-full flex items-center gap-3 rounded-xl border border-border/60 p-3 hover:border-primary/50 hover:shadow-sm transition-all cursor-pointer group text-left"
                  aria-label={`与${char.name}对话`}
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
                </button>
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
                <button
                  key={food.id}
                  type="button"
                  onClick={() => openDrawer('gallery', { ...food, target_type: 'food' })}
                  className="rounded-xl border border-border/60 overflow-hidden hover:border-accent/50 hover:shadow-sm transition-all cursor-pointer text-left"
                  aria-label={`查看${food.name}详情`}
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
                </button>
              ))}
            </div>
          </section>
        )}

        {/* 沉浸体验 */}
        <section className="panel">
          <div className="panel-header"><span>✨</span> 沉浸体验</div>
          <div className="panel-body space-y-3">
            <button
              type="button"
              onClick={() => setCanvasMode('street', city)}
              className="w-full h-28 rounded-xl bg-gradient-to-br from-primary/10 to-accent/15 border border-primary/20 flex flex-col items-center justify-center cursor-pointer hover:scale-[1.01] hover:shadow-sm transition-all group"
              aria-label={`查看${city.name}城市风光图`}
            >
              <span className="text-3xl mb-1.5 group-hover:scale-110 transition-transform">�️</span>
              <span className="text-sm font-semibold text-primary">查看城市风光</span>
            </button>

            <button
              onClick={() => setShowCheckin(true)}
              className="w-full py-3 bg-accent/15 text-accent font-semibold rounded-xl hover:bg-accent/25 transition-colors flex items-center justify-center gap-2"
            >
              <span>📸</span> 生成赛博打卡
            </button>
          </div>
        </section>
      </div>
      {/* CheckinFlow 滑入覆盖层 */}
      {showCheckin && (
        <CheckinFlow
          city={city}
          visitId={visitId}
          onClose={() => setShowCheckin(false)}
          onAchievementUnlocked={(ach) => setUnlockedAchievements(ach)}
        />
      )}

      {/* 成就解锁全屏庆祝 */}
      {unlockedAchievements.length > 0 && (
        <AchievementUnlock
          achievements={unlockedAchievements}
          onClose={() => setUnlockedAchievements([])}
        />
      )}
    </div>
  );
};
