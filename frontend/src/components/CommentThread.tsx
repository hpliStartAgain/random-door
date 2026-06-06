import React, { useEffect, useMemo, useState } from 'react';
import { MessageCircle, Send } from 'lucide-react';
import { api } from '../api';
import type { CommentItem, CommentTargetType } from '../api/types';
import { useUserStore } from '../store/useUserStore';

interface Props {
  targetType: CommentTargetType;
  targetId: number;
}

const DEFAULT_NICKNAME = '游客';

export const CommentThread: React.FC<Props> = ({ targetType, targetId }) => {
  const { userId } = useUserStore();
  const [comments, setComments] = useState<CommentItem[]>([]);
  const [nickname, setNickname] = useState(() => localStorage.getItem('comment_nickname') || '');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);
    api.getComments(targetType, targetId, 50)
      .then((res) => {
        if (!cancelled) setComments(res.comments);
      })
      .catch(() => {
        if (!cancelled) setError('评论加载失败');
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [targetType, targetId]);

  const danmaku = useMemo(() => comments.slice(-8), [comments]);

  const handleSubmit = async () => {
    const text = content.trim();
    const name = nickname.trim() || DEFAULT_NICKNAME;
    if (!text || submitting) return;
    setSubmitting(true);
    setError(null);
    try {
      const created = await api.createComment({
        target_type: targetType,
        target_id: targetId,
        user_id: userId,
        nickname: name,
        content: text,
      });
      localStorage.setItem('comment_nickname', name);
      setComments(prev => [...prev, created]);
      setContent('');
      setNickname(name);
    } catch (e: any) {
      setError(e?.message || '评论提交失败');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="space-y-3 pt-2">
      <style>{`
        @keyframes commentBullet {
          from { transform: translateX(105%); }
          to { transform: translateX(-120%); }
        }
        @media (prefers-reduced-motion: reduce) {
          .comment-bullet {
            animation: none !important;
            transform: none !important;
          }
        }
      `}</style>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-sm font-bold text-foreground">
          <MessageCircle className="h-4 w-4 text-primary" />
          <span>评论弹幕</span>
        </div>
        <span className="text-xs text-muted-foreground">{comments.length}</span>
      </div>

      <div className="relative h-24 overflow-hidden rounded-xl bg-[#1f2a27] border border-white/10">
        {danmaku.length > 0 ? danmaku.map((item, index) => (
          <div
            key={`${item.id}-${index}`}
            className="comment-bullet absolute max-w-[86%] overflow-hidden text-ellipsis whitespace-nowrap rounded-full bg-white/88 px-3 py-1 text-xs font-medium text-[#22302c] shadow-sm"
            style={{
              top: `${8 + (index % 4) * 22}%`,
              animation: `commentBullet ${12 + (index % 5) * 1.2}s linear ${index * -1.7}s infinite`,
              willChange: 'transform',
            }}
          >
            {item.nickname}: {item.content}
          </div>
        )) : (
          <div className="absolute inset-0 flex items-center justify-center text-xs text-white/55">
            {loading ? '加载中…' : '暂无评论'}
          </div>
        )}
      </div>

      <div className="space-y-2">
        <div className="flex gap-2">
          <input
            value={nickname}
            onChange={(e) => setNickname(e.target.value.slice(0, 32))}
            placeholder={DEFAULT_NICKNAME}
            className="w-24 shrink-0 rounded-lg border border-border bg-background px-3 py-2 text-xs outline-none focus:border-primary/50"
          />
          <div className="relative flex-1">
            <input
              value={content}
              onChange={(e) => setContent(e.target.value.slice(0, 200))}
              onKeyDown={(e) => { if (e.key === 'Enter') handleSubmit(); }}
              placeholder="写一句评论"
              className="w-full rounded-lg border border-border bg-background py-2 pl-3 pr-10 text-xs outline-none focus:border-primary/50"
            />
            <button
              onClick={handleSubmit}
              disabled={!content.trim() || submitting}
              className="absolute right-1 top-1 h-7 w-7 rounded-md bg-primary text-primary-foreground flex items-center justify-center disabled:opacity-45"
              aria-label="发送评论"
            >
              <Send className="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
        {error && <div className="text-xs text-red-500">{error}</div>}
      </div>

      {comments.length > 0 && (
        <div className="max-h-44 overflow-y-auto space-y-2 pr-1">
          {comments.slice().reverse().map((item) => (
            <div key={item.id} className="rounded-lg border border-border bg-card px-3 py-2">
              <div className="flex items-center justify-between gap-2">
                <span className="text-xs font-semibold text-primary truncate">{item.nickname}</span>
                <span className="text-[10px] text-muted-foreground shrink-0">
                  {new Date(item.created_at).toLocaleString()}
                </span>
              </div>
              <p className="mt-1 text-xs leading-relaxed text-foreground/80 break-words">{item.content}</p>
            </div>
          ))}
        </div>
      )}
    </section>
  );
};
