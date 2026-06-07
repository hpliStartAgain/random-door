import React, { useState } from 'react';
import { Copy, Download, Send, X, Check } from 'lucide-react';
import type { GuessCaptionResponse } from '../api/types';
import { downloadImage, copyImageToClipboard } from '../lib/shareImage';
import { api } from '../api';

interface Props {
  isOpen: boolean;
  onClose: () => void;
  shotUrl: string | null;
  caption: GuessCaptionResponse | null;
  targetName: string;
  cityId?: number | null;
  userId?: number | null;
}

export const GuessChallengeModal: React.FC<Props> = ({
  isOpen, onClose, shotUrl, caption, targetName, cityId, userId,
}) => {
  const [captionMode, setCaptionMode] = useState<'weibo' | 'moments'>('weibo');
  const [copyTextDone, setCopyTextDone] = useState(false);
  const [copyImgDone, setCopyImgDone] = useState(false);
  const [challengeDone, setChallengeDone] = useState(false);
  const [creatingChallenge, setCreatingChallenge] = useState(false);
  const [error, setError] = useState('');
  const [supportsClipboardImg] = useState(() => !!(window.ClipboardItem));

  if (!isOpen) return null;

  const activeCaption = caption
    ? (captionMode === 'weibo' ? caption.weibo : caption.moments)
    : '';

  const handleCopyText = async () => {
    if (!activeCaption) return;
    try {
      await navigator.clipboard.writeText(activeCaption);
      setCopyTextDone(true);
      setTimeout(() => setCopyTextDone(false), 1400);
    } catch {}
  };

  const handleSave = () => {
    if (!shotUrl) return;
    downloadImage(shotUrl, `猜猜我在哪-${targetName}-${Date.now()}.png`);
  };

  const handleCopyImg = async () => {
    if (!shotUrl) return;
    const ok = await copyImageToClipboard(shotUrl);
    if (ok) {
      setCopyImgDone(true);
      setTimeout(() => setCopyImgDone(false), 1400);
    }
  };

  const handleChallengeFriends = async () => {
    if (!cityId || !shotUrl || creatingChallenge) return;
    setCreatingChallenge(true);
    setError('');
    try {
      const challenge = await api.createGuessChallenge({
        user_id: userId,
        city_id: cityId,
        target_name: targetName,
        image_data_url: shotUrl.startsWith('data:') ? shotUrl : undefined,
        image_url: shotUrl.startsWith('data:') ? undefined : shotUrl,
        caption: activeCaption,
      });
      const sharePath = challenge.share_url || `/?guess=${challenge.code}`;
      const shareURL = new URL(sharePath, window.location.origin).toString();
      await navigator.clipboard.writeText(shareURL);
      setChallengeDone(true);
      setTimeout(() => setChallengeDone(false), 1800);
    } catch (e: any) {
      setError(e?.message || '挑战链接生成失败');
    } finally {
      setCreatingChallenge(false);
    }
  };

  return (
    <>
      {/* 蒙层 */}
      <div
        className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50"
        onClick={onClose}
      />

      {/* 弹层主体 */}
      <div className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-[360px] max-w-[calc(100vw-2rem)] bg-background rounded-2xl shadow-2xl border border-border overflow-hidden">
        {/* 头部 */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-border/50">
          <div>
            <h3 className="font-bold text-base text-foreground">猜猜我在哪</h3>
            <p className="text-xs text-muted-foreground mt-0.5">{targetName}</p>
          </div>
          <button
            onClick={onClose}
            className="w-8 h-8 rounded-full hover:bg-secondary flex items-center justify-center text-muted-foreground transition-colors"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="p-5 space-y-4">
          {/* 截图预览 */}
          {shotUrl && (
            <div className="rounded-xl overflow-hidden border border-border aspect-video bg-card">
              <img src={shotUrl} alt="截图预览" className="w-full h-full object-cover" />
            </div>
          )}

          {/* 文案区 */}
          {caption && (
            <div className="space-y-2">
              <div className="grid grid-cols-2 gap-1 rounded-lg bg-secondary p-1">
                <button
                  onClick={() => setCaptionMode('weibo')}
                  className={`py-1.5 rounded-md text-xs font-semibold transition-colors ${captionMode === 'weibo' ? 'bg-white text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'}`}
                >
                  微博
                </button>
                <button
                  onClick={() => setCaptionMode('moments')}
                  className={`py-1.5 rounded-md text-xs font-semibold transition-colors ${captionMode === 'moments' ? 'bg-white text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'}`}
                >
                  朋友圈
                </button>
              </div>
              <div className="rounded-xl bg-card border border-border p-3 text-sm leading-relaxed text-foreground/80 min-h-[80px]">
                {activeCaption}
              </div>
              {caption.hashtags?.length > 0 && (
                <div className="flex flex-wrap gap-1">
                  {caption.hashtags.map(tag => (
                    <span key={tag} className="text-[10px] px-2 py-0.5 bg-primary/10 text-primary rounded-full">{tag}</span>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* 操作按钮 */}
          <div className="space-y-2">
            {shotUrl && cityId && (
              <button
                onClick={handleChallengeFriends}
                disabled={creatingChallenge}
                className="w-full flex items-center justify-center gap-2 py-3 rounded-xl bg-primary text-primary-foreground text-sm font-bold hover:brightness-95 disabled:opacity-55 transition-all"
              >
                {challengeDone ? <Check className="h-4 w-4" /> : <Send className="h-4 w-4" />}
                {challengeDone ? '挑战链接已复制' : creatingChallenge ? '生成中...' : '挑战朋友'}
              </button>
            )}
            <div className="grid grid-cols-2 gap-2">
            {activeCaption && (
              <button
                onClick={handleCopyText}
                className="flex items-center justify-center gap-1.5 py-2.5 rounded-xl border border-border bg-secondary text-sm font-medium text-foreground hover:bg-border/50 transition-colors"
              >
                {copyTextDone ? <Check className="h-4 w-4 text-green-500" /> : <Copy className="h-4 w-4" />}
                {copyTextDone ? '已复制' : '复制文案'}
              </button>
            )}
            {shotUrl && (
              <button
                onClick={handleSave}
                className="flex items-center justify-center gap-1.5 py-2.5 rounded-xl border border-border bg-secondary text-sm font-medium text-foreground hover:bg-border/50 transition-colors"
              >
                <Download className="h-4 w-4" />
                保存图片
              </button>
            )}
            {shotUrl && supportsClipboardImg && (
              <button
                onClick={handleCopyImg}
                className="col-span-2 flex items-center justify-center gap-1.5 py-2.5 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:brightness-95 transition-all"
              >
                {copyImgDone ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                {copyImgDone ? '图片已复制' : '复制图片'}
              </button>
            )}
            </div>
            {error && <div className="text-xs text-red-500 text-center">{error}</div>}
          </div>
        </div>
      </div>
    </>
  );
};
