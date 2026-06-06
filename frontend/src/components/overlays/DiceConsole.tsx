import React, { useEffect } from 'react';
import { Dices, PlaneLanding } from 'lucide-react';
import { useViewStore } from '../../store/useViewStore';
import { useGameStore } from '../../store/useGameStore';
import { useUserStore } from '../../store/useUserStore';
import { useMapStore } from '../../store/useMapStore';

const DIRECTION_DEG: Record<string, number> = {
  '北': 0, '东北': 45, '东': 90, '东南': 135,
  '南': 180, '西南': 225, '西': 270, '西北': 315,
};

function CompassSVG({ direction }: { direction?: string }) {
  const deg = direction ? (DIRECTION_DEG[direction] ?? 0) : 0;
  return (
    <svg width="80" height="80" viewBox="0 0 80 80">
      <circle cx="40" cy="40" r="36" fill="none" stroke="rgba(43,58,54,0.16)" strokeWidth="1" />
      <circle cx="40" cy="40" r="34" fill="rgba(194,159,96,0.08)" />
      {(['N', 'E', 'S', 'W'] as const).map((label, i) => {
        const positions = [{ x: 40, y: 9 }, { x: 72, y: 44 }, { x: 40, y: 74 }, { x: 9, y: 44 }];
        const p = positions[i];
        return (
          <text key={label} x={p.x} y={p.y} textAnchor="middle" dominantBaseline="middle"
            fill={label === 'N' ? 'rgba(194,159,96,0.85)' : 'rgba(43,58,54,0.35)'}
            fontSize="8" fontWeight="bold" fontFamily="sans-serif">{label}</text>
        );
      })}
      <g style={{ transform: `rotate(${deg}deg)`, transformOrigin: '40px 40px', transition: 'transform 1.4s cubic-bezier(0.34,1.56,0.64,1)' }}>
        <polygon points="40,12 36.5,40 40,37 43.5,40" fill="#C29F60" />
        <polygon points="40,68 36.5,40 40,43 43.5,40" fill="rgba(43,58,54,0.45)" />
        <circle cx="40" cy="40" r="3.5" fill="#2B3A36" opacity="0.9" />
      </g>
    </svg>
  );
}

function FoxMascot({ rolling, landed }: { rolling?: boolean; landed?: boolean }) {
  return (
    <div
      className={`relative mx-auto h-24 w-24 ${rolling ? 'dice-spin' : ''}`}
      style={{ filter: 'drop-shadow(0 18px 28px rgba(43,58,54,0.18))' }}
    >
      <svg viewBox="0 0 96 96" className="h-full w-full">
        <path d="M18 28 L33 15 L37 35 Z" fill="#D47A3C" />
        <path d="M78 28 L63 15 L59 35 Z" fill="#D47A3C" />
        <path d="M24 31 C28 15 68 15 72 31 C84 48 75 78 48 82 C21 78 12 48 24 31Z" fill="#E48743" />
        <path d="M30 32 C37 43 42 55 48 80 C54 55 59 43 66 32 C60 68 36 68 30 32Z" fill="#FFF2D7" opacity="0.96" />
        <path d="M20 31 L32 20 L35 34 Z" fill="#8D3F28" opacity="0.55" />
        <path d="M76 31 L64 20 L61 34 Z" fill="#8D3F28" opacity="0.55" />
        <circle cx="36" cy="47" r="3.8" fill="#22302C" />
        <circle cx="60" cy="47" r="3.8" fill="#22302C" />
        <path d="M45 58 Q48 61 51 58" fill="none" stroke="#22302C" strokeWidth="2.4" strokeLinecap="round" />
        <path d="M48 54 L43 59 L53 59 Z" fill="#22302C" />
        <path d="M17 65 C7 64 7 48 21 46 C34 44 39 57 28 65 C35 67 42 72 48 80 C35 79 24 74 17 65Z" fill="#D87538" />
        <path d="M18 63 C12 60 14 52 22 52 C28 52 30 58 24 62 C31 64 37 70 41 75 C31 73 23 69 18 63Z" fill="#FFF2D7" opacity="0.92" />
        {landed && <circle cx="48" cy="48" r="42" fill="none" stroke="#C29F60" strokeWidth="2" strokeDasharray="6 6" opacity="0.55" />}
      </svg>
    </div>
  );
}

export const DiceConsole: React.FC = () => {
  const { setView, rollPhase, setRollPhase } = useViewStore();
  const { roll, lastRoll, nearestCity, setFromPoint, initGame } = useGameStore();
  const { userId } = useUserStore();
  const { mapInstance } = useMapStore();

  useEffect(() => {
    if (!userId || nearestCity) return;
    initGame(userId).catch(console.error);
  }, [userId, nearestCity, initGame]);

  const handleRoll = async () => {
    if (!userId || rollPhase === 'rolling') return;
    let lat = 39.9042, lng = 116.4074;
    if (mapInstance) {
      const center = mapInstance.getCenter();
      lat = center.lat;
      lng = center.lng;
    }
    setFromPoint({ lat, lng });
    setRollPhase('rolling');
    try {
      await roll(userId, nearestCity?.id || 1, lat, lng);
      setRollPhase('revealing');
    } catch (e) {
      console.error(e);
      setRollPhase('idle');
    }
  };

  const handleLand = () => {
    if (!lastRoll) return;
    setRollPhase('flying');
    setTimeout(() => {
      setRollPhase('idle');
      setView('CITY_DETAIL', lastRoll.target_city.id);
    }, 3400);
  };

  const handleBack = () => {
    setRollPhase('idle');
    setView('HOME');
  };

  const isRolling = rollPhase === 'rolling';
  const isRevealing = rollPhase === 'revealing';
  const isFlying = rollPhase === 'flying';

  return (
    <>
      <style>{`
        @keyframes diceSpin {
          0%   { transform: rotateX(0deg)   rotateY(0deg)   rotateZ(0deg); }
          33%  { transform: rotateX(210deg) rotateY(120deg) rotateZ(60deg); }
          66%  { transform: rotateX(420deg) rotateY(240deg) rotateZ(160deg); }
          100% { transform: rotateX(720deg) rotateY(360deg) rotateZ(300deg); }
        }
        @keyframes revealSlide {
          from { opacity: 0; transform: translateY(16px); }
          to   { opacity: 1; transform: translateY(0); }
        }
        @keyframes cityReveal {
          from { opacity: 0; letter-spacing: 0.4em; }
          to   { opacity: 1; letter-spacing: 0.04em; }
        }
        .dice-spin { animation: diceSpin 0.55s ease-in-out infinite; display: inline-block; }
        .reveal-in { animation: revealSlide 0.45s ease-out both; }
        .city-in   { animation: cityReveal 0.7s ease-out 0.25s both; }
      `}</style>

      {rollPhase !== 'idle' && (
        <div className="absolute inset-0 z-10"
          style={{ background: 'rgba(43,58,54,0.22)', backdropFilter: 'blur(3px)' }} />
      )}

      <div className={`absolute z-20 transition-all duration-500 ${
        isFlying ? 'opacity-0 scale-95 pointer-events-none' : 'opacity-100 scale-100'
      } ${
        rollPhase === 'idle'
          ? 'bottom-8 left-1/2 -translate-x-1/2 w-full max-w-sm px-4'
          : 'top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md px-4'
      }`}>
        <div className="rounded-3xl p-6 shadow-2xl border border-border/80 bg-background/95"
          style={{ backdropFilter: 'blur(24px)' }}>

          <div className="flex items-center justify-between mb-5">
            <h2 className="text-base font-bold text-primary tracking-wide">
              {rollPhase === 'idle' && '命运掷骰'}
              {isRolling && '骰子滚动中…'}
              {isRevealing && '目的地已确定'}
              {isFlying && '飞行中…'}
            </h2>
            <button onClick={handleBack}
              className="text-muted-foreground hover:text-primary transition-colors text-sm leading-none">✕</button>
          </div>

          <div className="flex justify-center my-5" style={{ perspective: '600px' }}>
            <FoxMascot rolling={isRolling} landed={isRevealing} />
          </div>

          {isRevealing && lastRoll && (
            <div className="reveal-in space-y-4">
              <div className="flex items-center gap-5 py-1">
                <div className="flex-shrink-0"><CompassSVG direction={lastRoll.direction} /></div>
                <div className="flex-1 space-y-3">
                  <div>
                    <div className="text-muted-foreground text-[10px] uppercase tracking-widest mb-0.5">方向</div>
                    <div className="text-foreground font-bold text-3xl">{lastRoll.direction}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground text-[10px] uppercase tracking-widest mb-0.5">距离</div>
                    <div className="text-foreground font-semibold text-xl">
                      {lastRoll.distance_km.toLocaleString()}
                      <span className="text-muted-foreground text-sm ml-1">km</span>
                    </div>
                  </div>
                </div>
              </div>

              <div className="border-t border-border/70" />

              <div className="text-center py-3">
                <div className="text-muted-foreground text-[10px] tracking-[0.35em] uppercase mb-2">目标城市</div>
                <div className="city-in text-primary font-black text-5xl"
                  style={{ textShadow: '0 14px 36px rgba(43,58,54,0.16), 0 0 54px rgba(194,159,96,0.22)' }}>
                  {lastRoll.target_city.name}
                </div>
                <div className="text-muted-foreground text-sm mt-1.5">{lastRoll.target_city.province}</div>
              </div>

              <button onClick={handleLand}
                className="w-full py-4 rounded-2xl text-primary-foreground font-bold text-base transition-all duration-200 flex items-center justify-center gap-2 hover:brightness-95"
                style={{ background: 'linear-gradient(135deg,hsl(var(--primary)),hsl(var(--accent)))', boxShadow: '0 16px 36px rgba(194,159,96,0.26)' }}>
                <PlaneLanding className="h-5 w-5" /> 降落 {lastRoll.target_city.name}
              </button>
            </div>
          )}

          {rollPhase === 'idle' && (
            <button onClick={handleRoll}
              className="w-full py-4 rounded-2xl text-primary-foreground font-bold text-base transition-all duration-200 flex items-center justify-center gap-2 hover:brightness-95"
              style={{ background: 'linear-gradient(135deg,hsl(var(--primary)),hsl(var(--accent)))', boxShadow: '0 16px 32px rgba(43,58,54,0.18)' }}>
              <Dices className="h-5 w-5" /> 抛掷骰子
            </button>
          )}

          {isRolling && (
            <div className="text-center text-muted-foreground text-sm py-1 tracking-wide">
              正在计算命运轨迹…
            </div>
          )}
        </div>
      </div>
    </>
  );
};
