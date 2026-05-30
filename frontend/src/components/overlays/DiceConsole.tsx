import React, { useState } from 'react';
import { useViewStore } from '../../store/useViewStore';
import { useGameStore } from '../../store/useGameStore';

export const DiceConsole: React.FC = () => {
  const { setView } = useViewStore();
  const { roll, rolling, targetCity } = useGameStore();
  const [showResult, setShowResult] = useState(false);

  const handleRoll = async () => {
    setShowResult(false);
    await roll(1, 39.9, 116.4);
    setShowResult(true);
  };

  return (
    <div className="absolute bottom-8 left-1/2 -translate-x-1/2 z-20 w-full max-w-sm px-4">
      <div className="glass-panel p-6 rounded-3xl flex flex-col items-center space-y-4 shadow-2xl">
        <div className="w-full flex justify-between items-center">
          <span className="font-bold text-foreground">命运掷骰</span>
          <button onClick={() => setView('HOME')} className="text-xs text-muted-foreground hover:text-primary">取消</button>
        </div>

        <div className="h-16 flex items-center justify-center">
          {rolling ? (
            <div className="animate-spin text-4xl">🎲</div>
          ) : showResult && targetCity ? (
            <div className="text-center">
              去往 <span className="text-primary font-bold">{targetCity.name}</span>
            </div>
          ) : (
            <div className="text-sm text-muted-foreground">准备就绪，点击开始</div>
          )}
        </div>

        {showResult && targetCity && !rolling ? (
          <button 
            onClick={() => setView('CITY_DETAIL', targetCity.id)}
            className="w-full py-2 bg-primary text-primary-foreground font-bold rounded-xl"
          >
            降落
          </button>
        ) : (
          <button 
            onClick={handleRoll}
            disabled={rolling}
            className="w-full py-2 bg-destructive text-destructive-foreground font-bold rounded-xl disabled:opacity-50"
          >
            抛掷
          </button>
        )}
      </div>
    </div>
  );
};
