import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useGameStore } from '../store/useGameStore';

export const GameModePage: React.FC = () => {
  const navigate = useNavigate();
  const { roll, rolling, targetCity } = useGameStore();

  const handleRoll = async () => {
    await roll(1, 1, 0, 0);
  };

  return (
    <div className="min-h-screen bg-background relative flex flex-col items-center justify-center p-4">
      <header className="absolute top-4 left-4 z-10">
        <div className="glass-panel inline-block px-6 py-2 rounded-full font-bold">
          <span className="text-primary cursor-pointer" onClick={() => navigate('/mode')}>&lt; 返回</span>
        </div>
      </header>

      <div className="glass-panel p-10 rounded-3xl max-w-sm w-full text-center space-y-8 z-10">
        <h2 className="text-2xl font-bold">命运的掷骰</h2>
        <div className="h-24 flex items-center justify-center">
          {rolling ? (
            <div className="animate-spin text-4xl text-destructive">🎲</div>
          ) : targetCity ? (
            <div className="text-lg">
              目标：<span className="font-bold text-primary">{targetCity.name}</span>
            </div>
          ) : (
            <div className="text-muted-foreground">点击掷骰子，寻找下一站</div>
          )}
        </div>
        
        <button
          onClick={handleRoll}
          disabled={rolling}
          className="w-full py-3 bg-destructive text-destructive-foreground font-semibold rounded-xl hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {rolling ? '掷骰中...' : '掷骰子'}
        </button>

        {targetCity && !rolling && (
          <button
            onClick={() => navigate(`/city/${targetCity.id}`)}
            className="w-full py-3 bg-primary text-primary-foreground font-semibold rounded-xl hover:opacity-90 transition-opacity"
          >
            前往 {targetCity.name}
          </button>
        )}
      </div>
    </div>
  );
};
