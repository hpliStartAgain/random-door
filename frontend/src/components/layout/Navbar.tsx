import React from 'react';
import { useViewStore } from '../../store/useViewStore';

export const Navbar: React.FC = () => {
  const { setView } = useViewStore();

  return (
    <header className="absolute top-0 left-0 w-full h-[60px] bg-background/90 backdrop-blur-md border-b border-border z-30 flex items-center justify-between px-6 shadow-sm">
      <div className="flex items-center gap-2 cursor-pointer" onClick={() => setView('HOME')}>
        <div className="w-6 h-6 rounded-full bg-primary flex items-center justify-center text-primary-foreground font-bold text-xs">
          任
        </div>
        <span className="font-bold text-lg text-foreground tracking-wide">
          任意门 <span className="font-normal text-muted-foreground ml-2 text-sm">推开门，遇见大美中国</span>
        </span>
      </div>
      
      <div className="flex items-center gap-4">
        <button className="text-sm font-medium text-foreground hover:text-primary transition-colors">
          成就墙
        </button>
        <div className="w-8 h-8 rounded-full bg-secondary border border-border flex items-center justify-center text-muted-foreground cursor-pointer hover:bg-card">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
        </div>
      </div>
    </header>
  );
};
