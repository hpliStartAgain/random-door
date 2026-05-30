import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';

export const ChatPage: React.FC = () => {
  const { id, cid } = useParams();
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-background flex flex-col p-4 max-w-2xl mx-auto">
      <header className="mb-4 shrink-0">
        <div className="glass-panel inline-block px-4 py-2 rounded-full font-bold">
          <span className="text-primary cursor-pointer" onClick={() => navigate(-1)}>&lt; 退出对话</span>
        </div>
      </header>
      
      <div className="glass-panel flex-1 rounded-2xl p-4 flex flex-col">
        <div className="flex-1 overflow-y-auto space-y-4 mb-4">
          <div className="text-center text-xs text-muted-foreground mb-4">与城市人物对话中...</div>
          <div className="flex items-start gap-2">
            <div className="w-8 h-8 rounded-full bg-primary/20 shrink-0" />
            <div className="glass-panel p-3 rounded-2xl rounded-tl-none max-w-[80%] text-sm">
              欢迎来到这里！你想了解些什么？
            </div>
          </div>
        </div>
        
        <div className="flex gap-2 shrink-0">
          <input 
            type="text" 
            placeholder="输入对话..." 
            className="flex-1 px-4 py-2 rounded-full border border-border bg-white outline-none focus:border-primary"
          />
          <button className="px-6 py-2 bg-primary text-primary-foreground rounded-full font-bold">
            发送
          </button>
        </div>
      </div>
    </div>
  );
};
