import React, { useState } from 'react';
import { useViewStore } from '../store/useViewStore';

export const RightDrawer: React.FC = () => {
  const { drawer, closeDrawer } = useViewStore();
  const [chatInput, setChatInput] = useState('');
  
  const isOpen = drawer.isOpen;
  
  return (
    <>
      {/* 遮罩 */}
      <div 
        className={`fixed inset-0 bg-black/20 backdrop-blur-sm z-40 transition-opacity duration-500 ${isOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'}`}
        onClick={closeDrawer}
      />
      
      {/* 抽屉主体 */}
      <div className={`fixed top-0 right-0 w-[420px] h-full bg-background/95 backdrop-blur-2xl shadow-2xl z-50 border-l border-border/50 flex flex-col transition-transform duration-500 ease-[cubic-bezier(0.16,1,0.3,1)] ${isOpen ? 'translate-x-0' : 'translate-x-full'}`}>
        
        {/* 头部 */}
        <div className="h-[60px] flex items-center justify-between px-6 border-b border-border/30 bg-background/50">
          <h3 className="font-bold text-lg text-primary tracking-tight">
            {drawer.type === 'chat' && drawer.data ? `与 ${drawer.data.name} 跨时空对话` : '风物鉴赏'}
          </h3>
          <button onClick={closeDrawer} className="w-8 h-8 rounded-full hover:bg-black/5 flex items-center justify-center text-muted-foreground transition-colors font-sans">
            ✕
          </button>
        </div>

        {/* 内容区 */}
        <div className="flex-1 overflow-y-auto p-6 relative">
           {drawer.type === 'chat' && drawer.data && (
             <div className="space-y-6">
                {/* 角色介绍卡 */}
                <div className="bg-primary/5 p-4 rounded-2xl border border-primary/10 flex gap-4 items-start shadow-sm">
                   <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center text-2xl shrink-0 border border-primary/20 shadow-inner">
                     🗿
                   </div>
                   <div>
                     <h4 className="font-bold text-primary mb-1 text-lg">{drawer.data.name} <span className="text-xs font-normal text-muted-foreground ml-1">· {drawer.data.dynasty}</span></h4>
                     <p className="text-sm text-foreground/80 leading-relaxed">{drawer.data.desc}</p>
                   </div>
                </div>

                {/* 聊天气泡 Demo */}
                <div className="flex flex-col gap-4 mt-6">
                  <div className="self-end bg-primary/90 text-primary-foreground px-4 py-2.5 rounded-2xl rounded-tr-sm max-w-[80%] text-sm shadow-md">
                    久仰大名，敢问阁下当年在长安是怎样一番光景？
                  </div>
                  <div className="self-start bg-card border border-border px-4 py-2.5 rounded-2xl rounded-tl-sm max-w-[80%] text-sm shadow-md flex gap-3 items-start">
                    <span className="text-xl pt-1">🗿</span>
                    <div>
                      <span className="text-xs text-primary font-bold mb-1 block">{drawer.data.name}</span>
                      “哈哈哈哈，长安市上酒家眠，天子呼来不上船！当年的风月，你这后生可懂？”
                    </div>
                  </div>
                </div>
             </div>
           )}

           {drawer.type === 'gallery' && drawer.data && (
             <div className="space-y-4">
                <h4 className="text-3xl font-bold text-foreground mb-4">{drawer.data.name}</h4>
                <div className="w-full h-72 bg-gradient-to-br from-secondary to-muted rounded-2xl flex items-center justify-center text-6xl shadow-inner border border-border/50">
                  🍜
                </div>
                <p className="text-muted-foreground leading-relaxed text-sm">
                  {drawer.data.desc}
                </p>
                <button className="w-full py-3 mt-4 bg-primary text-primary-foreground font-bold rounded-xl hover:opacity-90 transition-opacity shadow-sm">
                  保存图集至本地
                </button>
             </div>
           )}
        </div>

        {/* 底栏 (聊天输入) */}
        {drawer.type === 'chat' && (
          <div className="p-4 border-t border-border/30 bg-background/80 backdrop-blur pb-8">
            <div className="relative flex items-center">
              <input 
                type="text" 
                value={chatInput}
                onChange={(e) => setChatInput(e.target.value)}
                placeholder="在此输入对话，或使用语音..."
                className="w-full pl-5 pr-14 py-3.5 bg-card border border-border rounded-full text-sm outline-none focus:border-primary/50 transition-colors shadow-inner"
              />
              <button className="absolute right-2 w-10 h-10 rounded-full bg-primary flex items-center justify-center text-white hover:scale-105 transition-transform shadow-md">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3Z"/><path d="M19 10v2a7 7 0 0 1-14 0v-2"/><line x1="12" x2="12" y1="19" y2="22"/></svg>
              </button>
            </div>
          </div>
        )}
      </div>
    </>
  );
};
