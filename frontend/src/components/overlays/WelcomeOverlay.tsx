import React from 'react';
import { useViewStore } from '../../store/useViewStore';

const CITIES = ['北京', '上海', '西安', '成都', '杭州', '苏州', '南京', '武汉', '重庆', '广州', '厦门', '洛阳'];

export const WelcomeOverlay: React.FC = () => {
  const { enter, setView } = useViewStore();

  const handleEnter = (mode: 'FREE_EXPLORE' | 'GAME_DICE') => {
    enter();
    setView(mode);
  };

  return (
    <>
      <style>{`
        @keyframes welcomeFloat {
          0%, 100% { transform: translateY(0px); opacity: 0.05; }
          50% { transform: translateY(-18px); opacity: 0.12; }
        }
        @keyframes welcomeFadeUp {
          from { opacity: 0; transform: translateY(28px); }
          to { opacity: 1; transform: translateY(0); }
        }
        @keyframes welcomePulse {
          0%, 100% { opacity: 0.15; }
          50% { opacity: 0.35; }
        }
        .wf-1 { animation: welcomeFloat 7s ease-in-out infinite; }
        .wf-2 { animation: welcomeFloat 9s ease-in-out infinite 1.5s; }
        .wf-3 { animation: welcomeFloat 8s ease-in-out infinite 3s; }
        .wfu-0 { animation: welcomeFadeUp 0.7s ease-out 0.1s both; }
        .wfu-1 { animation: welcomeFadeUp 0.7s ease-out 0.25s both; }
        .wfu-2 { animation: welcomeFadeUp 0.7s ease-out 0.4s both; }
        .wfu-3 { animation: welcomeFadeUp 0.7s ease-out 0.6s both; }
        .wfu-4 { animation: welcomeFadeUp 0.7s ease-out 0.8s both; }
        .welcome-glow { animation: welcomePulse 3s ease-in-out infinite; }
      `}</style>

      <div className="fixed inset-0 z-50 overflow-hidden flex flex-col items-center justify-center">
        {/* Background */}
        <div className="absolute inset-0 bg-[#070c18]" />
        <div
          className="absolute inset-0"
          style={{
            backgroundImage:
              'linear-gradient(rgba(255,255,255,0.025) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.025) 1px, transparent 1px)',
            backgroundSize: '56px 56px',
          }}
        />
        <div
          className="absolute inset-0 welcome-glow"
          style={{
            background:
              'radial-gradient(ellipse 80% 60% at 50% 50%, rgba(99,102,241,0.12) 0%, transparent 70%)',
          }}
        />

        {/* Floating city names */}
        <div className="absolute inset-0 pointer-events-none select-none overflow-hidden">
          {CITIES.map((city, i) => (
            <span
              key={city}
              className={`absolute text-white font-bold text-sm wf-${(i % 3) + 1}`}
              style={{
                left: `${(i * 19 + 4) % 88}%`,
                top: `${(i * 17 + 6) % 82}%`,
                animationDelay: `${i * 0.4}s`,
                fontSize: `${12 + (i % 3) * 3}px`,
              }}
            >
              {city}
            </span>
          ))}
        </div>

        {/* Main content */}
        <div className="relative z-10 flex flex-col items-center text-center px-6 max-w-lg w-full">
          {/* Badge */}
          <div className="wfu-0 mb-6">
            <span className="inline-flex items-center gap-1.5 text-[10px] font-semibold tracking-[0.35em] text-white/35 uppercase border border-white/10 rounded-full px-4 py-1.5">
              ✦ AI 互动探索产品 ✦
            </span>
          </div>

          {/* Title */}
          <h1
            className="wfu-1 text-7xl sm:text-8xl font-black text-white mb-3 leading-none"
            style={{ textShadow: '0 0 80px rgba(99,102,241,0.45), 0 0 160px rgba(99,102,241,0.2)' }}
          >
            任意门
          </h1>

          {/* Subtitle */}
          <p className="wfu-2 text-base text-white/45 mb-10 leading-relaxed tracking-wide max-w-sm">
            随机降落中国城市，遇见历史地标<br />
            与古今人物跨时空对话，生成专属赛博大片
          </p>

          {/* CTA buttons */}
          <div className="wfu-3 flex flex-col sm:flex-row gap-3 w-full">
            <button
              onClick={() => handleEnter('FREE_EXPLORE')}
              className="flex-1 flex items-center justify-center gap-3 py-4 px-5 rounded-2xl bg-white/8 hover:bg-white/14 border border-white/15 hover:border-white/30 text-white transition-all duration-300 backdrop-blur-sm"
            >
              <span className="text-2xl">🗺️</span>
              <div className="text-left">
                <div className="font-bold text-sm">自由探索</div>
                <div className="text-xs text-white/40">主动选择目的地</div>
              </div>
            </button>

            <button
              onClick={() => handleEnter('GAME_DICE')}
              className="flex-1 flex items-center justify-center gap-3 py-4 px-5 rounded-2xl border border-indigo-500/40 hover:border-indigo-400/60 text-white transition-all duration-300"
              style={{
                background: 'linear-gradient(135deg, rgba(99,102,241,0.7), rgba(129,140,248,0.5))',
                boxShadow: '0 0 40px rgba(99,102,241,0.35)',
              }}
            >
              <span className="text-2xl">🎲</span>
              <div className="text-left">
                <div className="font-bold text-sm">随机漫游</div>
                <div className="text-xs text-indigo-200/60">让命运决定</div>
              </div>
            </button>
          </div>

          {/* Feature tags */}
          <div className="wfu-4 flex flex-wrap gap-2 mt-8 justify-center">
            {['12座城市', 'AI人物对话', '赛博打卡', '成就收集', '3D地图探索'].map((tag) => (
              <span
                key={tag}
                className="text-xs px-3 py-1 rounded-full text-white/25 border border-white/8"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>
      </div>
    </>
  );
};
