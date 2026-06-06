import React from 'react';
import { Dices, MapPinned } from 'lucide-react';
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
          0%, 100% { transform: translateY(0px); opacity: 0.14; }
          50% { transform: translateY(-18px); opacity: 0.28; }
        }
        @keyframes welcomeFadeUp {
          from { opacity: 0; transform: translateY(28px); }
          to { opacity: 1; transform: translateY(0); }
        }
        .wf-1 { animation: welcomeFloat 7s ease-in-out infinite; }
        .wf-2 { animation: welcomeFloat 9s ease-in-out infinite 1.5s; }
        .wf-3 { animation: welcomeFloat 8s ease-in-out infinite 3s; }
        .wfu-0 { animation: welcomeFadeUp 0.7s ease-out 0.1s both; }
        .wfu-1 { animation: welcomeFadeUp 0.7s ease-out 0.25s both; }
        .wfu-2 { animation: welcomeFadeUp 0.7s ease-out 0.4s both; }
        .wfu-3 { animation: welcomeFadeUp 0.7s ease-out 0.6s both; }
        .wfu-4 { animation: welcomeFadeUp 0.7s ease-out 0.8s both; }
      `}</style>

      <div className="fixed inset-0 z-50 overflow-hidden flex flex-col items-center justify-center">
        {/* Background */}
        <div className="absolute inset-0 bg-background" />
        <div
          className="absolute inset-0 opacity-80"
          style={{
            backgroundImage:
              'linear-gradient(rgba(43,58,54,0.055) 1px, transparent 1px), linear-gradient(90deg, rgba(43,58,54,0.055) 1px, transparent 1px)',
            backgroundSize: '56px 56px',
          }}
        />
        <div className="absolute inset-x-0 top-0 h-64 bg-gradient-to-b from-accent/15 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 h-80 bg-gradient-to-t from-primary/10 to-transparent" />

        {/* Floating city names */}
        <div className="absolute inset-0 pointer-events-none select-none overflow-hidden">
          {CITIES.map((city, i) => (
            <span
              key={city}
              className={`absolute text-primary font-serif-display font-bold text-sm wf-${(i % 3) + 1}`}
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
        <div className="relative z-10 flex flex-col items-center text-center px-6 max-w-2xl w-full">
          {/* Logo */}
          <img
            src="/icon-transparent.png"
            alt="任意门"
            className="wfu-0 w-24 h-24 sm:w-28 sm:h-28 object-contain mb-5 drop-shadow-[0_16px_32px_rgba(43,58,54,0.18)]"
          />

          {/* Badge */}
          <div className="wfu-0 mb-6">
            <span className="inline-flex items-center gap-1.5 text-[10px] font-semibold tracking-[0.35em] text-accent uppercase border border-accent/30 bg-accent/10 rounded-full px-4 py-1.5">
              AI 互动探索产品
            </span>
          </div>

          {/* Title */}
          <h1
            className="wfu-1 font-serif-display text-7xl sm:text-8xl font-black text-foreground mb-3 leading-none"
            style={{ textShadow: '0 16px 48px rgba(43,58,54,0.14), 0 0 80px rgba(194,159,96,0.18)' }}
          >
            任意门
          </h1>

          {/* Subtitle */}
          <p className="wfu-2 text-base text-foreground/65 mb-10 leading-relaxed tracking-wide max-w-sm">
            随机降落中国城市，遇见历史地标<br />
            与古今人物跨时空对话，生成专属赛博大片
          </p>

          {/* CTA buttons */}
          <div className="wfu-3 flex flex-col sm:flex-row gap-3 w-full">
            <button
              onClick={() => handleEnter('FREE_EXPLORE')}
              className="flex-1 flex items-center justify-center gap-3 py-4 px-5 rounded-2xl bg-primary hover:bg-primary/90 border border-primary/30 text-primary-foreground transition-all duration-300 shadow-[0_14px_34px_rgba(43,58,54,0.18)]"
            >
              <MapPinned className="h-6 w-6 shrink-0" />
              <div className="text-left">
                <div className="font-bold text-sm">自由探索</div>
                <div className="text-xs text-primary-foreground/70">主动选择目的地</div>
              </div>
            </button>

            <button
              onClick={() => handleEnter('GAME_DICE')}
              className="flex-1 flex items-center justify-center gap-3 py-4 px-5 rounded-2xl border border-accent/40 text-accent-foreground transition-all duration-300 hover:brightness-95"
              style={{
                background: 'linear-gradient(135deg,hsl(var(--accent)),hsl(var(--primary)))',
                boxShadow: '0 16px 36px rgba(194,159,96,0.24)',
              }}
            >
              <Dices className="h-6 w-6 shrink-0" />
              <div className="text-left">
                <div className="font-bold text-sm">随机漫游</div>
                <div className="text-xs text-accent-foreground/75">让命运决定</div>
              </div>
            </button>
          </div>

          {/* Feature tags */}
          <div className="wfu-4 flex flex-wrap gap-2 mt-8 justify-center">
            {['35座城市', 'AI人物对话', '赛博打卡', '成就收集', '3D地图探索'].map((tag) => (
              <span
                key={tag}
                className="text-xs px-3 py-1 rounded-full bg-card/80 text-muted-foreground border border-border/70 shadow-sm"
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
