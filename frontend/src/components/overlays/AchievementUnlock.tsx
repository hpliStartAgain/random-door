import React, { useEffect, useState } from 'react';

export interface Achievement {
  code: string;
  name: string;
  description?: string;
}

interface Props {
  achievements: Achievement[];
  onClose: () => void;
}

const PARTICLE_COUNT = 28;

export const AchievementUnlock: React.FC<Props> = ({ achievements, onClose }) => {
  const [visible, setVisible] = useState(false);
  const ach = achievements[0];

  useEffect(() => {
    requestAnimationFrame(() => setVisible(true));
    const t = setTimeout(onClose, 6000);
    return () => clearTimeout(t);
  }, [onClose]);

  if (!ach) return null;

  return (
    <>
      <style>{`
        @keyframes achParticle {
          0%   { transform: translateY(0) scale(1) rotate(0deg); opacity: 1; }
          100% { transform: translateY(-120px) scale(0.3) rotate(360deg); opacity: 0; }
        }
        @keyframes achBadge {
          0%   { transform: scale(0.3) rotate(-20deg); opacity: 0; }
          60%  { transform: scale(1.12) rotate(4deg); opacity: 1; }
          100% { transform: scale(1) rotate(0deg); opacity: 1; }
        }
        @keyframes achFadeIn {
          from { opacity: 0; transform: translateY(20px); }
          to   { opacity: 1; transform: translateY(0); }
        }
        .ach-badge { animation: achBadge 0.7s cubic-bezier(0.34,1.56,0.64,1) 0.2s both; }
        .ach-text  { animation: achFadeIn 0.5s ease-out 0.8s both; }
      `}</style>

      <div
        className={`fixed inset-0 z-[150] flex flex-col items-center justify-center transition-opacity duration-500 ${visible ? 'opacity-100' : 'opacity-0'}`}
        style={{ background: 'rgba(0,0,0,0.75)', backdropFilter: 'blur(6px)' }}
        onClick={onClose}
      >
        {/* Particles */}
        <div className="absolute inset-0 pointer-events-none overflow-hidden">
          {Array.from({ length: PARTICLE_COUNT }).map((_, i) => (
            <div
              key={i}
              className="absolute w-2 h-2 rounded-sm"
              style={{
                left: `${(i * 13 + 5) % 95}%`,
                top: `${40 + (i % 4) * 5}%`,
                background: ['#C29F60', '#2B3A36', '#8A7A45', '#F5F3EB', '#A8894E'][i % 5],
                animation: `achParticle ${1.2 + (i % 4) * 0.3}s ease-out ${i * 0.06}s both`,
              }}
            />
          ))}
        </div>

        {/* Card */}
        <div className="relative flex flex-col items-center text-center px-8 py-10 max-w-xs" onClick={e => e.stopPropagation()}>
          <div className="ach-badge w-24 h-24 rounded-full bg-gradient-to-br from-accent to-primary flex items-center justify-center mb-6 shadow-2xl"
            style={{ boxShadow: '0 0 60px rgba(194,159,96,0.48)' }}>
            <span className="text-4xl">🏆</span>
          </div>

          <div className="ach-text space-y-2">
            <div className="text-xs text-accent/90 tracking-[0.3em] uppercase font-semibold">成就解锁</div>
            <div className="text-white font-black text-3xl">{ach.name}</div>
            <div className="text-white/60 text-sm leading-relaxed">{ach.description}</div>
            {achievements.length > 1 && (
              <div className="text-white/40 text-xs mt-2">+{achievements.length - 1} 个其他成就</div>
            )}
          </div>

          <button
            onClick={onClose}
            className="mt-8 px-8 py-2.5 rounded-full bg-white/10 hover:bg-white/20 border border-white/20 text-white text-sm font-semibold transition-colors"
          >
            太棒了！
          </button>
        </div>
      </div>
    </>
  );
};
