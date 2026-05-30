import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useUserStore } from '../store/useUserStore';

export const HomePage: React.FC = () => {
  const navigate = useNavigate();
  const initUser = useUserStore((state) => state.initUser);

  useEffect(() => {
    initUser();
  }, [initUser]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <div className="glass-panel p-10 rounded-3xl max-w-lg text-center space-y-6">
        <h1 className="text-4xl font-extrabold text-primary tracking-wider">AI 城市漫游</h1>
        <p className="text-muted-foreground text-lg">
          穿梭神州大地，邂逅文化记忆，探索你的专属赛博旅程。
        </p>
        <button
          onClick={() => navigate('/mode')}
          className="mt-4 px-8 py-3 bg-primary text-primary-foreground font-semibold rounded-xl shadow hover:opacity-90 transition-all"
        >
          启程
        </button>
      </div>
    </div>
  );
};
