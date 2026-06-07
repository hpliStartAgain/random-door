import React, { useEffect, useState } from 'react';
import { RefreshCw, ArrowRight, X } from 'lucide-react';
import { useViewStore } from '../../store/useViewStore';
import { useGameStore } from '../../store/useGameStore';
import { useUserStore } from '../../store/useUserStore';
import { useMapStore } from '../../store/useMapStore';
import { useCityStore } from '../../store/useCityStore';
import { foxImages } from '../../assets/foxImages';
import { AchievementUnlock } from './AchievementUnlock';
import type { Achievement } from './AchievementUnlock';

type VisualPhase = 'idle' | 'loading' | 'revealing' | 'result' | 'flying';

function LoadingCityFlash({ cities }: { cities: { id: number; name: string }[] }) {
  if (cities.length === 0) return null;
  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none select-none">
      {cities.map((city, i) => (
        <span
          key={city.id}
          className="absolute font-serif-display font-bold text-primary/20"
          style={{
            fontSize: `${12 + (i % 4) * 4}px`,
            left: `${(i * 21 + 6) % 84}%`,
            top: `${(i * 17 + 10) % 78}%`,
            animation: `foxCityFloat ${2.2 + i * 0.35}s ease-in-out infinite ${i * 0.22}s`,
          }}
        >
          {city.name}
        </span>
      ))}
    </div>
  );
}

export const RandomCityModal: React.FC = () => {
  const { setView, rollPhase, setRollPhase } = useViewStore();
  const { roll, lastRoll, nearestCity, setFromPoint, initGame, reset } = useGameStore();
  const { userId } = useUserStore();
  const { mapInstance, flyTo } = useMapStore();
  const { cities } = useCityStore();

  const [showResult, setShowResult] = useState(false);
  const [flashCities, setFlashCities] = useState<typeof cities>([]);
  const [unlockedAchievements, setUnlockedAchievements] = useState<Achievement[]>([]);

  useEffect(() => {
    if (!userId || nearestCity) return;
    initGame(userId).catch(console.error);
  }, [userId, nearestCity, initGame]);

  // Auto-transition revealing → result after 850ms
  useEffect(() => {
    if (rollPhase === 'revealing') {
      setShowResult(false);
      const t = setTimeout(() => setShowResult(true), 850);
      return () => clearTimeout(t);
    }
    setShowResult(false);
  }, [rollPhase]);

  const vp: VisualPhase =
    rollPhase === 'rolling'   ? 'loading'   :
    rollPhase === 'revealing' && !showResult ? 'revealing' :
    rollPhase === 'revealing' && showResult  ? 'result'    :
    rollPhase === 'flying'    ? 'flying'    :
    'idle';

  const doRoll = async () => {
    if (!userId) return;
    let lat = 39.9042, lng = 116.4074;
    if (mapInstance) {
      const center = mapInstance.getCenter();
      lat = center.lat;
      lng = center.lng;
    }
    setFlashCities([...cities].sort(() => Math.random() - 0.5).slice(0, 12));
    setFromPoint({ lat, lng });
    setRollPhase('rolling');
    try {
      const result = await roll(userId, nearestCity?.id || 1, lat, lng);
      if (result.unlocked_achievements?.length) {
        setUnlockedAchievements(result.unlocked_achievements);
      }
      setRollPhase('revealing');
    } catch (e) {
      console.error(e);
      setRollPhase('idle');
    }
  };

  const handleDepart = () => {
    if (!lastRoll) return;
    setRollPhase('flying');
    flyTo(lastRoll.target_city.lng, lastRoll.target_city.lat);
    setTimeout(() => {
      setRollPhase('idle');
      setView('CITY_DETAIL', lastRoll.target_city.id);
    }, 1500);
  };

  const handleReroll = () => {
    reset();
    setRollPhase('idle');
    setTimeout(doRoll, 180);
  };

  const handleClose = () => {
    setRollPhase('idle');
    setView('HOME');
  };

  const isActive = vp !== 'idle' && vp !== 'flying';
  const foxSrc =
    vp === 'loading'   ? foxImages.compass      :
    vp === 'revealing' ? foxImages.cityCard     :
    vp === 'result'    ? foxImages.ticketReveal :
    foxImages.fortuneTable;

  const foxAnimClass =
    vp === 'loading'   ? 'fox-spin-compass' :
    vp === 'result'    ? 'fox-pop'          :
    vp === 'revealing' ? 'fox-pop'          :
    'fox-float';

  return (
    <>
      <style>{`
        @keyframes foxCityFloat {
          0%, 100% { opacity: 0.18; transform: translateY(0); }
          50%       { opacity: 0.55; transform: translateY(-12px); }
        }
        @keyframes revealSlide {
          from { opacity: 0; transform: translateY(16px); }
          to   { opacity: 1; transform: translateY(0); }
        }
        @keyframes cityReveal {
          from { opacity: 0; letter-spacing: 0.4em; }
          to   { opacity: 1; letter-spacing: 0.04em; }
        }
        .reveal-in { animation: revealSlide 0.45s ease-out both; }
        .city-in   { animation: cityReveal 0.7s ease-out 0.15s both; }
      `}</style>
      {unlockedAchievements.length > 0 && (
        <AchievementUnlock
          achievements={unlockedAchievements}
          onClose={() => setUnlockedAchievements([])}
        />
      )}

      {/* Backdrop blur overlay */}
      <div
        className={`absolute inset-0 transition-all duration-500 ${
          isActive ? 'opacity-100' : 'opacity-0 pointer-events-none'
        }`}
        style={{ background: 'rgba(43,58,54,0.26)', backdropFilter: 'blur(4px)' }}
      />

      {/* Floating city names during loading */}
      {vp === 'loading' && <LoadingCityFlash cities={flashCities} />}

      {/* Main card */}
      <div
        className={`absolute z-20 w-full max-w-sm px-4 transition-all duration-500 ${
          vp === 'flying'
            ? 'opacity-0 scale-95 pointer-events-none'
            : 'opacity-100 scale-100'
        } ${
          vp === 'idle'
            ? 'bottom-10 left-1/2 -translate-x-1/2'
            : 'top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2'
        }`}
      >
        <div
          className="rounded-3xl p-6 shadow-2xl border border-border/70"
          style={{ background: 'rgba(250,249,245,0.97)', backdropFilter: 'blur(24px)' }}
        >
          {/* Header */}
          <div className="flex items-center justify-between mb-4">
            <h2 className="font-serif-display font-bold text-primary tracking-wide text-base">
              {vp === 'idle'      && '命运占卜台'}
              {vp === 'loading'   && '命运罗盘转动中'}
              {vp === 'revealing' && '城市卡牌揭晓中'}
              {vp === 'result'    && '下一站'}
            </h2>
            <button
              onClick={handleClose}
              className="text-muted-foreground hover:text-primary transition-colors p-1"
            >
              <X className="h-4 w-4" />
            </button>
          </div>

          {/* Fox image */}
          <div className="flex justify-center my-4">
            <div className="relative">
              <img
                src={foxSrc}
                alt="fox"
                className={`w-44 h-44 object-contain rounded-2xl ${foxAnimClass}`}
                onError={(e) => {
                  (e.currentTarget as HTMLImageElement).style.opacity = '0';
                }}
              />
              {/* City name badge on ticket reveal image */}
              {vp === 'result' && lastRoll && (
                <div className="absolute -bottom-3 left-1/2 -translate-x-1/2 whitespace-nowrap">
                  <span
                    className="inline-block bg-accent text-accent-foreground font-black text-sm px-4 py-1 rounded-full shadow-lg"
                    style={{ boxShadow: '0 8px 20px rgba(194,159,96,0.35)' }}
                  >
                    {lastRoll.target_city.name}
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Subtitle / province info */}
          <div className="text-center text-muted-foreground text-sm mt-6 mb-3 min-h-[20px]">
            {vp === 'idle'      && '让狐狸为你寻找下一座城市'}
            {vp === 'loading'   && (
              <span className="animate-pulse">命运罗盘正在寻找下一站…</span>
            )}
            {vp === 'revealing' && '正在翻开命运选择的城市…'}
            {vp === 'result' && lastRoll && (
              <span className="text-foreground/80 font-medium">
                {lastRoll.target_city.province}
              </span>
            )}
          </div>

          {/* Big city name (result) */}
          {vp === 'result' && lastRoll && (
            <div className="city-in text-center mb-5">
              <div
                className="font-black text-5xl text-primary font-serif-display"
                style={{
                  textShadow: '0 14px 36px rgba(43,58,54,0.14), 0 0 54px rgba(194,159,96,0.20)',
                }}
              >
                {lastRoll.target_city.name}
              </div>
              <div className="text-muted-foreground text-xs mt-1.5 tracking-wide">
                {lastRoll.direction} · {lastRoll.distance_km.toLocaleString()} km
              </div>
            </div>
          )}

          {/* Action buttons */}
          {vp === 'idle' && (
            <button
              onClick={doRoll}
              className="w-full py-4 rounded-2xl text-primary-foreground font-bold text-base flex items-center justify-center gap-2 transition-all hover:brightness-95"
              style={{
                background: 'linear-gradient(135deg,hsl(var(--primary)),hsl(var(--accent)))',
                boxShadow: '0 16px 32px rgba(43,58,54,0.18)',
              }}
            >
              抽取下一站
            </button>
          )}

          {vp === 'result' && (
            <div className="flex gap-3 reveal-in">
              <button
                onClick={handleReroll}
                className="flex-1 py-3 rounded-2xl border border-border text-sm font-semibold hover:bg-secondary transition-colors flex items-center justify-center gap-1.5"
              >
                <RefreshCw className="h-3.5 w-3.5" />
                重新抽取
              </button>
              <button
                onClick={handleDepart}
                className="flex-1 py-3 rounded-2xl text-primary-foreground font-bold text-sm flex items-center justify-center gap-1.5 hover:brightness-95 transition-all"
                style={{
                  background: 'linear-gradient(135deg,hsl(var(--primary)),hsl(var(--accent)))',
                  boxShadow: '0 12px 26px rgba(194,159,96,0.24)',
                }}
              >
                出发去看看 <ArrowRight className="h-4 w-4" />
              </button>
            </div>
          )}
        </div>
      </div>
    </>
  );
};
