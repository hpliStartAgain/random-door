import React from 'react';
import { useToastStore } from '../store/useToastStore';

const ICONS: Record<string, string> = { success: '✓', error: '✕', info: 'ℹ' };
const COLORS: Record<string, string> = {
  success: 'bg-green-500/90 border-green-400/40',
  error: 'bg-red-500/90 border-red-400/40',
  info: 'bg-primary/90 border-primary/40',
};

export const Toast: React.FC = () => {
  const { toasts, dismiss } = useToastStore();

  return (
    <div className="fixed bottom-6 right-6 z-[200] flex flex-col gap-2 pointer-events-none">
      {toasts.map((t) => (
        <div
          key={t.id}
          className={`flex items-center gap-3 px-4 py-3 rounded-xl border text-white text-sm shadow-xl backdrop-blur-xl pointer-events-auto animate-in slide-in-from-right duration-300 ${COLORS[t.type]}`}
        >
          <span className="font-bold text-base leading-none">{ICONS[t.type]}</span>
          <span className="flex-1">{t.message}</span>
          <button onClick={() => dismiss(t.id)} className="text-white/60 hover:text-white ml-1 leading-none">✕</button>
        </div>
      ))}
    </div>
  );
};
