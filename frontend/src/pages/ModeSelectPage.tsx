import React from 'react';
import { useNavigate } from 'react-router-dom';

export const ModeSelectPage: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-background p-4 space-y-8">
      <h2 className="text-3xl font-bold text-foreground">选择你的探索方式</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 w-full max-w-4xl">
        <div 
          onClick={() => navigate('/explore')}
          className="glass-panel p-8 rounded-2xl cursor-pointer hover:shadow-md transition-shadow flex flex-col items-center space-y-4"
        >
          <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center text-primary text-2xl font-bold">
            自
          </div>
          <h3 className="text-xl font-semibold">自由探索</h3>
          <p className="text-muted-foreground text-center">拨开云雾，自由点亮你想去的城市，查看地标与风物。</p>
        </div>
        <div 
          onClick={() => navigate('/game')}
          className="glass-panel p-8 rounded-2xl cursor-pointer hover:shadow-md transition-shadow flex flex-col items-center space-y-4"
        >
          <div className="w-16 h-16 bg-destructive/10 rounded-full flex items-center justify-center text-destructive text-2xl font-bold">
            游
          </div>
          <h3 className="text-xl font-semibold">游戏互动</h3>
          <p className="text-muted-foreground text-center">掷出命运的骰子，让随机的风带你去未知的远方。</p>
        </div>
      </div>
    </div>
  );
};
