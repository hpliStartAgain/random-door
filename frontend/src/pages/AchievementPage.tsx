import React, { useEffect, useState } from 'react';
import { api } from '../api';
import { useUserStore } from '../store/useUserStore';
import { useViewStore } from '../store/useViewStore';
import type { AchievementWallResponse } from '../api/types';

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

      <div className="glass-panel p-8 rounded-3xl grid grid-cols-2 md:grid-cols-4 gap-6">
        {data?.unlocked.map((ach) => (
          <div key={ach.code} className="flex flex-col items-center space-y-2 text-center">
            <div className="w-20 h-20 bg-primary/20 rounded-full border-4 border-primary flex items-center justify-center text-4xl shadow-md">
              🏅
            </div>
            <span className="font-bold text-sm">{ach.name}</span>
            <span className="text-xs text-primary font-bold">已解锁</span>
            <span className="text-[10px] text-muted-foreground">{new Date(ach.unlocked_at).toLocaleDateString()}</span>
          </div>
        ))}
        {data?.progress.map((prog) => (
          <div key={prog.code} className="flex flex-col items-center space-y-2 text-center opacity-70">
            <div className="w-20 h-20 bg-muted rounded-full border-4 border-border flex items-center justify-center text-3xl shadow-inner relative overflow-hidden">
              <div className="absolute bottom-0 left-0 right-0 bg-primary/20" style={{ height: `${(prog.current / prog.target) * 100}%` }}></div>
              <span className="relative z-10">⏳</span>
            </div>
            <span className="font-bold text-sm">{prog.code}</span>
            <span className="text-xs text-muted-foreground">{prog.current} / {prog.target}</span>
          </div>
        ))}
        {data?.locked.map((ach) => (
          <div key={ach.code} className="flex flex-col items-center space-y-2 text-center opacity-40 grayscale">
            <div className="w-20 h-20 bg-muted rounded-full border-4 border-border flex items-center justify-center text-3xl">
              🔒
            </div>
            <span className="font-bold text-sm text-muted-foreground">{ach.name}</span>
            <span className="text-[10px] text-muted-foreground line-clamp-2">{ach.description}</span>
          </div>
        ))}
      </div>
    </div>
  );
};
