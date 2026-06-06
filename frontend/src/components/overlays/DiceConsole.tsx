import React, { useEffect } from 'react';
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
      <circle cx="40" cy="40" r="36" fill="none" stroke="rgba(255,255,255,0.1)" strokeWidth="1" />
      <circle cx="40" cy="40" r="34" fill="rgba(255,255,255,0.03)" />
      {(['N', 'E', 'S', 'W'] as const).map((label, i) => {
        const positions = [{ x: 40, y: 9 }, { x: 72, y: 44 }, { x: 40, y: 74 }, { x: 9, y: 44 }];
        const p = positions[i];
        return (
          <text key={label} x={p.x} y={p.y} textAnchor="middle" dominantBaseline="middle"
            fill={label === 'N' ? 'rgba(255,255,255,0.5)' : 'rgba(255,255,255,0.2)'}
            fontSize="8" fontWeight="bold" fontFamily="sans-serif">{label}</text>
        );
      })}
      <g style={{ transform: `rotate(${deg}deg)`, transformOrigin: '40px 40px', transition: 'transform 1.4s cubic-bezier(0.34,1.56,0.64,1)' }}>
        <polygon points="40,12 36.5,40 40,37 43.5,40" fill="#ef4444" />
        <polygon points="40,68 36.5,40 40,43 43.5,40" fill="rgba(255,255,255,0.25)" />
        <circle cx="40" cy="40" r="3.5" fill="white" opacity="0.9" />
      </g>
    </svg>
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
          style={{ background: 'rgba(4,7,18,0.6)', backdropFilter: 'blur(3px)' }} />
      )}

      <div className={`absolute z-20 transition-all duration-500 ${
        isFlying ? 'opacity-0 scale-95 pointer-events-none' : 'opacity-100 scale-100'
      } ${
        rollPhase === 'idle'
          ? 'bottom-8 left-1/2 -translate-x-1/2 w-full max-w-sm px-4'
          : 'top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-md px-4'
      }`}>
        <div className="rounded-3xl p-6 shadow-2xl border border-white/10"
          style={{ background: 'rgba(8,12,28,0.97)', backdropFilter: 'blur(24px)' }}>

          <div className="flex items-center justify-between mb-5">
            <h2 className="text-base font-bold text-white/90 tracking-wide">
              {rollPhase === 'idle' && '命运掷骰'}
              {isRolling && '骰子滚动中…'}
              {isRevealing && '目的地已确定'}
              {isFlying && '飞行中…'}
            </h2>
            <button onClick={handleBack}
              className="text-white/25 hover:text-white/55 transition-colors text-sm leading-none">✕</button>
          </div>

          <div className="flex justify-center my-5" style={{ perspective: '600px' }}>
            <span className={isRolling ? 'dice-spin' : ''}
              style={{ fontSize: 80, lineHeight: 1, display: 'inline-block' }}>
              {isRevealing ? '🎯' : '🎲'}
            </span>
          </div>

          {isRevealing && lastRoll && (
            <div className="reveal-in space-y-4">
              <div className="flex items-center gap-5 py-1">
                <div className="flex-shrink-0"><CompassSVG direction={lastRoll.direction} /></div>
                <div className="flex-1 space-y-3">
                  <div>
                    <div className="text-white/30 text-[10px] uppercase tracking-widest mb-0.5">方向</div>
                    <div className="text-white font-bold text-3xl">{lastRoll.direction}</div>
                  </div>
                  <div>
                    <div className="text-white/30 text-[10px] uppercase tracking-widest mb-0.5">距离</div>
                    <div className="text-white font-semibold text-xl">
                      {lastRoll.distance_km.toLocaleString()}
                      <span className="text-white/40 text-sm ml-1">km</span>
                    </div>
                  </div>
                </div>
              </div>

              <div className="border-t border-white/8" />

              <div className="text-center py-3">
                <div className="text-white/30 text-[10px] tracking-[0.35em] uppercase mb-2">目标城市</div>
                <div className="city-in text-white font-black text-5xl"
                  style={{ textShadow: '0 0 40px rgba(99,102,241,0.65), 0 0 80px rgba(99,102,241,0.25)' }}>
                  {lastRoll.target_city.name}
                </div>
                <div className="text-white/35 text-sm mt-1.5">{lastRoll.target_city.province}</div>
              </div>

              <button onClick={handleLand}
                className="w-full py-4 rounded-2xl text-white font-bold text-base transition-all duration-200 flex items-center justify-center gap-2"
                style={{ background: 'linear-gradient(135deg,#6366f1,#818cf8)', boxShadow: '0 0 40px rgba(99,102,241,0.45)' }}>
                <span>✈️</span> 降落 {lastRoll.target_city.name}
              </button>
            </div>
          )}

          {rollPhase === 'idle' && (
            <button onClick={handleRoll}
              className="w-full py-4 rounded-2xl text-white font-bold text-base transition-all duration-200"
              style={{ background: 'linear-gradient(135deg,#4f46e5,#6366f1)', boxShadow: '0 0 30px rgba(99,102,241,0.35)' }}>
              🎲 抛掷骰子
            </button>
          )}

          {isRolling && (
            <div className="text-center text-white/35 text-sm py-1 tracking-wide">
              正在计算命运轨迹…
            </div>
          )}
        </div>
      </div>
    </>
  );
};
