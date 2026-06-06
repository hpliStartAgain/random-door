import React, { useState, useEffect, useRef } from 'react';
import { useViewStore } from '../store/useViewStore';
import { useUserStore } from '../store/useUserStore';
import { api } from '../api';
import { CommentThread } from './CommentThread';
import type { CommentTargetType } from '../api/types';

type Msg = { role: 'user' | 'assistant'; content: string };

const SUGGESTED_QUESTIONS = [
  '请做个自我介绍',
  '这里最值得一看的是什么？',
  '这座城市有什么历史典故？',
  '你最难忘的经历是什么？',
];

export const RightDrawer: React.FC = () => {
  const { drawer, closeDrawer, activeCityId } = useViewStore();
  const { userId } = useUserStore();
  const [chatInput, setChatInput] = useState('');
  const [messages, setMessages] = useState<Msg[]>([]);
  const [loading, setLoading] = useState(false);
  const [typingText, setTypingText] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const typeTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const isOpen = drawer.isOpen;

  useEffect(() => {
    if (isOpen && drawer.type === 'chat' && drawer.data) {
      setChatInput('');
      const char = drawer.data;
      const greeting = char.character_type === 'history'
        ? `汝好！吾乃${char.name}，久候矣。有何要叙？`
        : `你好！我是${char.name}。${char.dialect_style ? `（${char.dialect_style}）` : ''} 有什么想聊的？`;
      setMessages([{ role: 'assistant', content: greeting }]);
    }
  }, [isOpen, drawer.data?.id]);

  useEffect(() => {
    if (messages.length === 0) { setTypingText(''); return; }
    const last = messages[messages.length - 1];
    if (last.role !== 'assistant') return;
    let i = 0;
    setIsTyping(true);
    setTypingText('');
    if (typeTimerRef.current) clearInterval(typeTimerRef.current);
    typeTimerRef.current = setInterval(() => {
      i++;
      setTypingText(last.content.slice(0, i));
      if (i >= last.content.length) {
        clearInterval(typeTimerRef.current!);
        setIsTyping(false);
      }
    }, 22);
    return () => { if (typeTimerRef.current) clearInterval(typeTimerRef.current); };
  }, [messages.length]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, typingText]);

  const handleSend = async (text?: string) => {
    const msg = (text ?? chatInput).trim();
    if (!msg || !userId || !activeCityId || !drawer.data || loading) return;
    setMessages(prev => [...prev, { role: 'user', content: msg }]);
    setChatInput('');
    setLoading(true);
    try {
      const res = await api.chat(userId, activeCityId, drawer.data.id, msg);
      setMessages(prev => [...prev, { role: 'assistant', content: res.reply }]);
    } catch (e) {
      console.error(e);
      setMessages(prev => [...prev, { role: 'assistant', content: '抱歉，暂时无法回答，请稍后再试。' }]);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSend(); }
  };

  const avatarEl = drawer.data?.avatar_url
    ? <img src={drawer.data.avatar_url} alt={drawer.data.name} className="w-full h-full object-cover rounded-full" />
    : <span className="text-2xl">🗿</span>;

  const commentTarget: { targetType: CommentTargetType; targetId: number } | null = drawer.data?.id
    ? drawer.type === 'chat'
      ? { targetType: 'character', targetId: drawer.data.id }
      : { targetType: (drawer.data.target_type || 'food') as CommentTargetType, targetId: drawer.data.id }
    : null;

  return (
    <>
      <div
        className={`fixed inset-0 bg-black/20 backdrop-blur-sm z-40 transition-opacity duration-500 ${isOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'}`}
        onClick={closeDrawer}
      />

      <div className={`fixed top-0 right-0 w-full sm:w-[420px] max-w-[100vw] h-full bg-background/95 backdrop-blur-2xl shadow-2xl z-50 border-l border-border/50 flex flex-col transition-transform duration-500 ease-[cubic-bezier(0.16,1,0.3,1)] ${isOpen ? 'translate-x-0' : 'translate-x-full'}`}>

        <div className="h-[60px] flex items-center justify-between px-6 border-b border-border/30 bg-background/50 shrink-0">
          <h3 className="font-bold text-lg text-primary tracking-tight truncate">
            {drawer.type === 'chat' && drawer.data ? `与 ${drawer.data.name} 跨时空对话` : '风物鉴赏'}
          </h3>
          <button onClick={closeDrawer} className="w-8 h-8 rounded-full hover:bg-black/5 flex items-center justify-center text-muted-foreground transition-colors shrink-0">✕</button>
        </div>

        <div className="flex-1 overflow-y-auto p-6">
          {drawer.type === 'chat' && drawer.data && (
            <div className="space-y-5">
              <div className="bg-primary/5 p-4 rounded-2xl border border-primary/10 flex gap-4 items-start shadow-sm">
                <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center shrink-0 border border-primary/20 shadow-inner overflow-hidden">
                  {avatarEl}
                </div>
                <div>
                  <h4 className="font-bold text-primary mb-1 text-base">
                    {drawer.data.name}
                    <span className="text-xs font-normal text-muted-foreground ml-1.5">
                      · {drawer.data.character_type === 'history' ? '历史人物' : drawer.data.character_type === 'culture' ? '文化符号' : drawer.data.character_type}
                    </span>
                  </h4>
                  {drawer.data.dialect_style && (
                    <p className="text-xs text-muted-foreground">{drawer.data.dialect_style}</p>
                  )}
                </div>
              </div>

              <div className="flex flex-col gap-4">
                {messages.map((msg, i) => {
                  const isLastAssistant = i === messages.length - 1 && msg.role === 'assistant';
                  const content = isLastAssistant ? typingText : msg.content;
                  return msg.role === 'user' ? (
                    <div key={i} className="self-end flex gap-2.5 items-start max-w-[82%] flex-row-reverse">
                      <div className="w-7 h-7 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center shrink-0 overflow-hidden mt-0.5">
                        <img src="/icon-transparent.png" alt="我" className="w-full h-full object-contain" />
                      </div>
                      <div className="bg-primary/90 text-primary-foreground px-4 py-2.5 rounded-2xl rounded-tr-sm text-sm shadow-md whitespace-pre-wrap">
                        {msg.content}
                      </div>
                    </div>
                  ) : (
                    <div key={i} className="self-start bg-card border border-border px-4 py-3 rounded-2xl rounded-tl-sm max-w-[82%] text-sm shadow-md flex gap-3 items-start">
                      <div className="w-7 h-7 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center shrink-0 overflow-hidden mt-0.5">
                        {drawer.data?.avatar_url
                          ? <img src={drawer.data.avatar_url} alt="" className="w-full h-full object-cover" />
                          : <span className="text-sm">🗿</span>}
                      </div>
                      <div>
                        <span className="text-xs text-primary font-bold mb-1 block">{drawer.data?.name}</span>
                        <span className="whitespace-pre-wrap leading-relaxed">{content}
                          {isLastAssistant && isTyping && <span className="inline-block w-0.5 h-3.5 bg-primary ml-0.5 animate-pulse" />}
                        </span>
                      </div>
                    </div>
                  );
                })}
                {loading && (
                  <div className="self-start bg-card border border-border px-4 py-3 rounded-2xl rounded-tl-sm text-sm shadow-md flex gap-2 items-center text-muted-foreground">
                    <span className="flex gap-1">
                      {[0,1,2].map(n => <span key={n} className="w-1.5 h-1.5 rounded-full bg-primary/60 animate-bounce" style={{ animationDelay: `${n*0.15}s` }} />)}
                    </span>
                    <span>思考中…</span>
                  </div>
                )}
                <div ref={messagesEndRef} />
              </div>

              {!loading && messages.length <= 1 && (
                <div className="flex flex-wrap gap-2 pt-1">
                  {SUGGESTED_QUESTIONS.map(q => (
                    <button key={q} onClick={() => handleSend(q)}
                      className="text-xs px-3 py-1.5 rounded-full bg-secondary hover:bg-primary/10 border border-border hover:border-primary/30 text-foreground/70 hover:text-primary transition-colors">
                      {q}
                    </button>
                  ))}
                </div>
              )}

              {commentTarget && (
                <CommentThread targetType={commentTarget.targetType} targetId={commentTarget.targetId} />
              )}
            </div>
          )}

          {drawer.type === 'gallery' && drawer.data && (
            <div className="space-y-4">
              <h4 className="text-3xl font-bold text-foreground mb-4">{drawer.data.name}</h4>
              {drawer.data.image_url ? (
                <img src={drawer.data.image_url} alt={drawer.data.name} className="w-full h-72 object-cover rounded-2xl shadow-md" />
              ) : (
                <div className="w-full h-72 bg-gradient-to-br from-secondary to-muted rounded-2xl flex items-center justify-center text-6xl shadow-inner border border-border/50">🍜</div>
              )}
              <p className="text-muted-foreground leading-relaxed text-sm">{drawer.data.description}</p>
              {commentTarget && (
                <CommentThread targetType={commentTarget.targetType} targetId={commentTarget.targetId} />
              )}
            </div>
          )}
        </div>

        {drawer.type === 'chat' && (
          <div className="p-4 border-t border-border/30 bg-background/80 backdrop-blur pb-safe shrink-0">
            <div className="relative flex items-center">
              <input
                type="text"
                value={chatInput}
                onChange={(e) => setChatInput(e.target.value)}
                onKeyDown={handleKeyDown}
                disabled={loading}
                placeholder="问他一个问题…"
                className="w-full pl-5 pr-14 py-3.5 bg-card border border-border rounded-full text-sm outline-none focus:border-primary/50 transition-colors shadow-inner disabled:opacity-50"
              />
              <button
                onClick={() => handleSend()}
                disabled={loading || !chatInput.trim()}
                className="absolute right-2 w-10 h-10 rounded-full bg-primary flex items-center justify-center text-white hover:scale-105 transition-transform shadow-md disabled:opacity-50 disabled:hover:scale-100"
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="22" y1="2" x2="11" y2="13" /><polygon points="22 2 15 22 11 13 2 9 22 2" />
                </svg>
              </button>
            </div>
          </div>
        )}
      </div>
    </>
  );
};
