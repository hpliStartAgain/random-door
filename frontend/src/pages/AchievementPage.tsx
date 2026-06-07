import React, { useEffect, useState } from 'react';
import { api } from '../api';
import { useUserStore } from '../store/useUserStore';
import { useViewStore } from '../store/useViewStore';
import type { AchievementWallResponse } from '../api/types';
import { foxImages } from '../assets/foxImages';

export const AchievementPage: React.FC = () => {
  const { setView } = useViewStore();
  const { userId } = useUserStore();
  const [data, setData] = useState<AchievementWallResponse | null>(null);

  useEffect(() => {
    if (userId) {
      api.getAchievements(userId).then(setData).catch(console.error);
    }
  }, [userId]);

  return (
    <div className="min-h-screen bg-background p-4 max-w-4xl mx-auto pt-24">
      <header className="mb-8">
        <div className="glass-panel inline-block px-6 py-2 rounded-full font-bold shadow-sm border border-border">
          <span className="text-primary cursor-pointer hover:underline" onClick={() => setView('HOME')}>&lt; 返回</span>
          <span className="mx-4">成就墙</span>
        </div>
      </header>

      {/* 顶部主视觉 */}
      <div className="flex flex-col items-center mb-10 text-center">
        <img
          src={foxImages.passportStamp}
          alt="狐狸护照印章"
          className="w-40 h-40 object-contain fox-float drop-shadow-lg mb-4"
          onError={(e) => { (e.currentTarget as HTMLImageElement).style.display = 'none'; }}
        />
        <h2 className="font-serif-display text-2xl font-black text-primary mb-1">旅行成就墙</h2>
        <p className="text-sm text-muted-foreground">每一次打卡，都会在这里留下一枚印章</p>
      </div>

      <div className="glass-panel p-8 rounded-3xl grid grid-cols-2 md:grid-cols-4 gap-6">
        {(data?.unlocked ?? []).map((ach) => (
          <div key={ach.code} className="flex flex-col items-center space-y-2 text-center">
            <div className="w-20 h-20 bg-primary/20 rounded-full border-4 border-primary flex items-center justify-center text-4xl shadow-md">
              🏅
            </div>
            <span className="font-bold text-sm">{ach.name}</span>
            <span className="text-xs text-primary font-bold">已解锁</span>
            <span className="text-[10px] text-muted-foreground">{new Date(ach.unlocked_at).toLocaleDateString()}</span>
          </div>
        ))}
        {(data?.progress ?? []).map((prog) => {
          const clamped = Math.min(prog.current, prog.target);
          const pct = Math.min(clamped / prog.target, 1) * 100;
          const done = clamped >= prog.target;
          return (
            <div key={prog.code} className="flex flex-col items-center space-y-2 text-center opacity-70">
              <div className="w-20 h-20 bg-muted rounded-full border-4 border-border flex items-center justify-center text-3xl shadow-inner relative overflow-hidden">
                <div className="absolute bottom-0 left-0 right-0 bg-primary/20" style={{ height: `${pct}%` }}></div>
                <span className="relative z-10">⏳</span>
              </div>
              <span className="font-bold text-sm">{prog.code}</span>
              <span className="text-xs text-muted-foreground">{done ? '已完成' : `${clamped} / ${prog.target}`}</span>
            </div>
          );
        })}
        {(data?.locked ?? []).map((ach) => (
          <div key={ach.code} className="flex flex-col items-center space-y-2 text-center opacity-40 grayscale">
            <div className="w-20 h-20 bg-muted rounded-full border-4 border-border flex items-center justify-center text-3xl">
              🔒
            </div>
            <span className="font-bold text-sm text-muted-foreground">{ach.name}</span>
            <span className="text-[10px] text-muted-foreground line-clamp-2">{ach.description}</span>
          </div>
        ))}
        {data && (data.unlocked ?? []).length === 0 && (data.progress ?? []).length === 0 && (data.locked ?? []).length === 0 && (
          <div className="col-span-2 md:col-span-4 text-center text-muted-foreground py-10 text-sm">
            还没有任何成就，快去探索城市吧！
          </div>
        )}
      </div>
    </div>
  );
};
