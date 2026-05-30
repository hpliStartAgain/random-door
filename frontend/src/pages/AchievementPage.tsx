import React from 'react';
import { useNavigate } from 'react-router-dom';

export const AchievementPage: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-background p-4 max-w-4xl mx-auto">
      <header className="mb-8">
        <div className="glass-panel inline-block px-6 py-2 rounded-full font-bold">
          <span className="text-primary cursor-pointer" onClick={() => navigate(-1)}>&lt; 返回</span>
          <span className="mx-4">成就墙</span>
        </div>
      </header>

      <div className="glass-panel p-8 rounded-3xl grid grid-cols-2 md:grid-cols-4 gap-6">
        <div className="flex flex-col items-center space-y-2">
          <div className="w-20 h-20 bg-primary/20 rounded-full border-4 border-primary flex items-center justify-center">
            🏆
          </div>
          <span className="font-bold text-sm">初次打卡</span>
          <span className="text-xs text-muted-foreground">已解锁</span>
        </div>
        <div className="flex flex-col items-center space-y-2 opacity-50 grayscale">
          <div className="w-20 h-20 bg-muted rounded-full border-4 border-border flex items-center justify-center">
            ?
          </div>
          <span className="font-bold text-sm text-muted-foreground">古都巡礼</span>
          <span className="text-xs text-muted-foreground">0 / 3</span>
        </div>
      </div>
    </div>
  );
};
