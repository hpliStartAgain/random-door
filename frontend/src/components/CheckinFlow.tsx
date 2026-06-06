import React, { useState, useRef, useEffect } from 'react';
import type { CityDetail, Landmark } from '../api/types';
import { api } from '../api';
import { useUserStore } from '../store/useUserStore';
import { CheckinPoster } from './CheckinPoster';
import type { Achievement } from './overlays/AchievementUnlock';

interface Props {
  city: CityDetail;
  visitId?: number;
  onClose: () => void;
  onAchievementUnlocked: (achievements: Achievement[]) => void;
}

type Step = 'landmark' | 'upload' | 'confirm';

const STEPS: Step[] = ['landmark', 'upload', 'confirm'];

export const CheckinFlow: React.FC<Props> = ({ city, visitId, onClose, onAchievementUnlocked }) => {
  const { userId } = useUserStore();
  const [step, setStep] = useState<Step>('landmark');
  const [selectedLandmark, setSelectedLandmark] = useState<Landmark | null>(null);
  const [selfieFile, setSelfieFile] = useState<File | null>(null);
  const [selfiePreview, setSelfiePreview] = useState<string | null>(null);
  const [generating, setGenerating] = useState(false);
  const [progress, setProgress] = useState(0);
  const [generatedUrl, setGeneratedUrl] = useState<string | null>(null);
  const [confirming, setConfirming] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const progressTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const landmarks: Landmark[] = city.landmarks?.length > 0
    ? city.landmarks
    : [{ id: 1, name: city.name, image_url: city.cover_image_url, description: city.intro }];

  useEffect(() => {
    if (!generating) {
      if (progressTimerRef.current) clearInterval(progressTimerRef.current);
      return;
    }
    setProgress(0);
    progressTimerRef.current = setInterval(() => {
      setProgress((p) => (p >= 88 ? 88 : p + Math.random() * 9));
    }, 350);
    return () => { if (progressTimerRef.current) clearInterval(progressTimerRef.current); };
  }, [generating]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setSelfieFile(file);
    setSelfiePreview(URL.createObjectURL(file));
    setError(null);
  };

  const handleGenerate = async () => {
    if (!selfieFile || !userId || !selectedLandmark) return;
    setGenerating(true);
    setError(null);
    try {
      const formData = new FormData();
      formData.append('selfie_file', selfieFile);
      formData.append('user_id', userId.toString());
      formData.append('city_id', city.id.toString());
      formData.append('landmark_id', selectedLandmark.id.toString());
      const res = await api.generateImage(formData);
      setProgress(100);
      setTimeout(() => {
        setGeneratedUrl(res.generated_image_url);
        setGenerating(false);
        setStep('confirm');
      }, 400);
    } catch {
      setGenerating(false);
      setError('AI 生成失败，请重试');
    }
  };

  const handleConfirm = async () => {
    if (!userId || !selectedLandmark || !generatedUrl) return;
    setConfirming(true);
    setError(null);
    try {
      const res = await api.createCheckin(userId, city.id, selectedLandmark.id, visitId, generatedUrl);
      if (res.unlocked_achievements?.length > 0) {
        onAchievementUnlocked(res.unlocked_achievements);
      }
      onClose();
    } catch {
      setError('打卡失败，请重试');
    } finally {
      setConfirming(false);
    }
  };

  const stepIdx = STEPS.indexOf(step);

  return (
    <div className="absolute inset-0 bg-background z-40 flex flex-col animate-in slide-in-from-right duration-300">
      <div className="flex items-center gap-3 px-5 py-4 border-b border-border/40 shrink-0">
        <button onClick={onClose} className="text-sm text-muted-foreground hover:text-foreground transition-colors">
          ← 返回
        </button>
        <h3 className="font-bold text-base flex-1 text-center">赛博打卡</h3>
        <div className="flex gap-1.5">
          {STEPS.map((s, i) => (
            <div key={s} className={`h-1.5 rounded-full transition-all duration-300 ${i === stepIdx ? 'w-6 bg-primary' : i < stepIdx ? 'w-4 bg-primary/40' : 'w-4 bg-border'}`} />
          ))}
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-5 py-5 space-y-4">
        {step === 'landmark' && (
          <>
            <p className="text-sm text-muted-foreground">选择地标，AI 将为你合成专属赛博大片</p>
            <div className="space-y-2">
              {landmarks.map((lm) => (
                <button
                  key={lm.id}
                  onClick={() => { setSelectedLandmark(lm); setStep('upload'); }}
                  className="w-full flex items-center gap-3 p-3 rounded-xl border border-border hover:border-primary/50 hover:bg-primary/5 transition-all text-left group"
                >
                  {lm.image_url ? (
                    <img src={lm.image_url} alt={lm.name} className="w-14 h-14 object-cover rounded-lg shrink-0" />
                  ) : (
                    <div className="w-14 h-14 bg-accent/10 rounded-lg flex items-center justify-center text-xl shrink-0">🏯</div>
                  )}
                  <div className="flex-1 min-w-0">
                    <div className="font-semibold text-sm">{lm.name}</div>
                    {lm.description && <div className="text-xs text-muted-foreground line-clamp-2 mt-0.5">{lm.description}</div>}
                  </div>
                  <span className="text-muted-foreground group-hover:text-primary text-sm">→</span>
                </button>
              ))}
            </div>
          </>
        )}

        {step === 'upload' && selectedLandmark && (
          <>
            <div className="text-sm text-muted-foreground flex items-center gap-1.5">
              <span>📍</span>
              <span>地标：<span className="font-semibold text-foreground">{selectedLandmark.name}</span></span>
            </div>
            <div
              onClick={() => !generating && fileInputRef.current?.click()}
              className={`relative w-full aspect-square rounded-2xl border-2 border-dashed flex flex-col items-center justify-center cursor-pointer transition-all overflow-hidden ${selfiePreview ? 'border-primary/40' : 'border-border hover:border-primary/40 bg-secondary/30'}`}
            >
              {selfiePreview ? (
                <img src={selfiePreview} alt="自拍预览" className="absolute inset-0 w-full h-full object-cover" />
              ) : (
                <>
                  <span className="text-4xl mb-2">📸</span>
                  <span className="text-sm text-muted-foreground">点击选择或拍摄照片</span>
                </>
              )}
            </div>
            <input ref={fileInputRef} type="file" accept="image/*" capture="user" className="hidden" onChange={handleFileChange} />

            {generating && (
              <div className="space-y-1.5">
                <div className="flex justify-between text-xs text-muted-foreground">
                  <span className="animate-pulse">AI 正在穿梭时空合成大片…</span>
                  <span>{Math.round(progress)}%</span>
                </div>
                <div className="w-full bg-border rounded-full h-1.5 overflow-hidden">
                  <div className="bg-primary h-1.5 rounded-full transition-all duration-300" style={{ width: `${progress}%` }} />
                </div>
              </div>
            )}

            {error && <p className="text-sm text-red-500 text-center">{error}</p>}

            <button
              onClick={handleGenerate}
              disabled={!selfieFile || generating}
              className="w-full py-3.5 rounded-2xl font-bold text-sm text-white transition-all disabled:opacity-40 disabled:cursor-not-allowed"
              style={selfieFile && !generating ? { background: 'linear-gradient(135deg,#6366f1,#818cf8)', boxShadow: '0 0 20px rgba(99,102,241,0.3)' } : {}}
            >
              {generating ? '生成中…' : '✨ 生成赛博大片'}
            </button>
          </>
        )}

        {step === 'confirm' && generatedUrl && selectedLandmark && (
          <>
            <p className="text-sm font-semibold text-foreground">你的专属赛博大片已生成！</p>
            <CheckinPoster imageUrl={generatedUrl} cityName={city.name} landmarkName={selectedLandmark.name} />
            {error && <p className="text-sm text-red-500 text-center">{error}</p>}
            <div className="flex gap-3 pb-4">
              <button
                onClick={() => { setStep('upload'); setGeneratedUrl(null); }}
                className="flex-1 py-3 rounded-xl border border-border text-sm font-semibold hover:bg-secondary transition-colors"
              >
                重新生成
              </button>
              <button
                onClick={handleConfirm}
                disabled={confirming}
                className="flex-1 py-3 rounded-xl text-white font-bold text-sm disabled:opacity-50 transition-all"
                style={{ background: 'linear-gradient(135deg,#6366f1,#818cf8)', boxShadow: '0 0 20px rgba(99,102,241,0.3)' }}
              >
                {confirming ? '打卡中…' : '✓ 确认打卡'}
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
};
