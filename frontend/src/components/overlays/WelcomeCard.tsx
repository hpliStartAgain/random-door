import React from 'react';
import { Dices, MapPinned } from 'lucide-react';
import { useViewStore } from '../../store/useViewStore';

export const WelcomeCard: React.FC = () => {
  const { setView } = useViewStore();

  return (
    <div className="absolute inset-0 z-10 flex flex-col items-center justify-center pointer-events-none p-4">
      <div className="bg-background/95 backdrop-blur-xl border border-border p-10 rounded-3xl max-w-lg w-full text-center space-y-8 pointer-events-auto shadow-xl transition-all duration-500 ease-out">
        <div>
          <h1 className="font-serif-display text-4xl font-extrabold text-foreground tracking-wider mb-2">
            AI 城市漫游
          </h1>
          <p className="text-muted-foreground">拨开云雾，邂逅中国城市的文化与记忆</p>
        </div>
        
        <div className="flex flex-col gap-4">
          <button 
            onClick={() => setView('FREE_EXPLORE')}
            className="w-full py-3 bg-primary text-primary-foreground font-bold rounded-xl hover:opacity-90 transition-opacity flex items-center justify-center gap-2 shadow"
          >
            <MapPinned className="h-4 w-4" /> 自由探索
          </button>
          
          <button 
            onClick={() => setView('GAME_DICE')}
            className="w-full py-3 bg-accent text-accent-foreground font-bold rounded-xl hover:opacity-90 transition-opacity flex items-center justify-center gap-2 shadow"
          >
            <Dices className="h-4 w-4" /> 随机漫游
          </button>
        </div>
      </div>
    </div>
  );
};
