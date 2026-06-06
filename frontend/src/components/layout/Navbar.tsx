import React from 'react';
import { useViewStore } from '../../store/useViewStore';

interface NavbarProps {
  onAdmin?: () => void;
}

export const Navbar: React.FC<NavbarProps> = ({ onAdmin }) => {
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
        <button
          onClick={() => setView('ASSETS')}
          className="text-sm font-medium text-foreground hover:text-primary transition-colors"
        >
          我的资产
        </button>
        <button 
          onClick={() => setView('ACHIEVEMENT')}
          className="text-sm font-medium text-foreground hover:text-primary transition-colors"
        >
          成就墙
        </button>
        {onAdmin && (
          <button
            onClick={onAdmin}
            title="后台管理"
            className="w-8 h-8 rounded-full bg-secondary border border-border flex items-center justify-center text-muted-foreground hover:bg-card hover:text-primary transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/></svg>
          </button>
        )}
      </div>
    </header>
  );
};
