import React, { useRef, useState } from 'react';
import { Camera, Sparkles, X } from 'lucide-react';
import { useViewStore } from '../store/useViewStore';
import { useUserStore } from '../store/useUserStore';
import { api } from '../api';
import type { CityDetail, GuessCaptionResponse } from '../api/types';
import { CheckinFlow } from './CheckinFlow';
import { AchievementUnlock } from './overlays/AchievementUnlock';
import type { Achievement } from './overlays/AchievementUnlock';
import { GuessChallengeModal } from './GuessChallengeModal';
import { SoundscapeControl } from './SoundscapeControl';

export const StreetViewCanvas: React.FC = () => {
  const { setCanvasMode, streetTarget } = useViewStore();
  const { userId } = useUserStore();
  const rootRef = useRef<HTMLDivElement>(null);
  const [shotUrl, setShotUrl] = useState<string | null>(null);
  const [caption, setCaption] = useState<GuessCaptionResponse | null>(null);
  const [loadingCaption, setLoadingCaption] = useState(false);
  const [loadingCheckin, setLoadingCheckin] = useState(false);
  const [checkinCity, setCheckinCity] = useState<CityDetail | null>(null);
  const [sceneFile, setSceneFile] = useState<File | null>(null);
  const [unlockedAchievements, setUnlockedAchievements] = useState<Achievement[]>([]);
  const [showChallengeModal, setShowChallengeModal] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const panoramaUrl = streetTarget?.cover_image_url || '/static/landmarks/beijing_cover.png';
  const targetName = streetTarget?.name || '未知异境';
  const captionCityId = streetTarget?.city_id || streetTarget?.id;

  const captureCurrentView = () => {
    let nextShot = '';
    const canvas = rootRef.current?.querySelector('canvas') as HTMLCanvasElement | null;
    if (canvas) {
      try {
        nextShot = canvas.toDataURL('image/png');
      } catch {
        nextShot = '';
      }
    }
    if (!nextShot) nextShot = panoramaUrl;
    setShotUrl(nextShot);
    return nextShot;
  };

  const captureSceneFile = async (): Promise<File | null> => {
    const canvas = rootRef.current?.querySelector('canvas') as HTMLCanvasElement | null;
    if (canvas) {
      const blob = await new Promise<Blob | null>((resolve) => {
        try {
          canvas.toBlob(resolve, 'image/jpeg', 0.86);
        } catch {
          resolve(null);
        }
      });
      if (blob) {
        return new File([blob], `street-view-${Date.now()}.jpg`, { type: 'image/jpeg', lastModified: Date.now() });
      }
    }

    try {
      return await api.fetchLocalImageFile(panoramaUrl, `street-view-${Date.now()}.png`);
    } catch {
      return null;
    }
  };

  const handleGenerateCaption = async () => {
    if (!captionCityId || loadingCaption) return;
    captureCurrentView();
    setLoadingCaption(true);
    setError(null);
    try {
      const res = await api.generateGuessCaption({
        user_id: userId,
        city_id: captionCityId,
        target_name: targetName,
        scene_hint: `全景截图视角：${targetName}`,
      });
      setCaption(res);
      setShowChallengeModal(true);
    } catch (e: any) {
      setError(e?.message || '文案生成失败');
    } finally {
      setLoadingCaption(false);
    }
  };

  const handleStartCheckin = async () => {
    if (!captionCityId || loadingCheckin) return;
    setLoadingCheckin(true);
    setError(null);
    try {
      const captured = await captureSceneFile();
      if (!captured) {
        setError('全景视角截图失败，将使用地标参考图');
      }
      const detail = await api.getCityDetail(captionCityId);
      setSceneFile(captured);
      setCheckinCity(detail);
    } catch (e: any) {
      setError(e?.message || '打卡入口加载失败');
    } finally {
      setLoadingCheckin(false);
    }
  };

  return (
    <div ref={rootRef} className="absolute inset-0 z-0 bg-black flex flex-col overflow-hidden">
      
      <div className="absolute inset-0 z-0 opacity-90">
        <img
          src={panoramaUrl}
          alt={targetName}
          className="w-full h-full object-cover"
        />
      </div>
      
      {/* 操作蒙层 (上下暗角) */}
      <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-black/30 pointer-events-none" />

      {/* 退出按钮 */}
      <div className="absolute top-6 left-6 z-10 pointer-events-auto">
        <button 
          onClick={() => setCanvasMode('map')}
          className="px-6 py-2 bg-black/40 hover:bg-black/60 backdrop-blur-lg text-white font-bold rounded-full border border-white/20 shadow-2xl flex items-center gap-2 transition-all hover:scale-105"
        >
          <span>&larr;</span> 退出风光浏览
        </button>
      </div>

      <div className="absolute top-20 right-4 sm:top-6 sm:right-6 z-10 pointer-events-auto w-[360px] max-w-[calc(100vw-2rem)]">
        <div className="rounded-2xl border border-white/15 bg-black/42 backdrop-blur-xl text-white shadow-2xl overflow-hidden">
          <div className="flex items-center justify-between gap-3 px-4 py-3 border-b border-white/10">
            <div className="flex items-center gap-2 min-w-0">
              <Sparkles className="h-4 w-4 text-[#F1C76B] shrink-0" />
              <span className="text-sm font-bold truncate">猜一猜</span>
            </div>
            {(shotUrl || caption) && (
              <button
                onClick={() => { setShotUrl(null); setCaption(null); setError(null); }}
                className="h-7 w-7 rounded-full hover:bg-white/10 flex items-center justify-center"
                aria-label="关闭猜一猜面板"
              >
                <X className="h-4 w-4" />
              </button>
            )}
          </div>

          <div className="p-4 space-y-3">
            <button
              onClick={handleGenerateCaption}
              disabled={loadingCaption || !captionCityId}
              className="w-full py-2.5 rounded-xl bg-white text-[#22302C] text-sm font-bold flex items-center justify-center gap-2 hover:bg-white/90 disabled:opacity-55 transition-colors"
            >
              <Camera className="h-4 w-4" />
              {loadingCaption ? '生成中…' : '截图生成文案'}
            </button>
            <button
              onClick={handleStartCheckin}
              disabled={loadingCheckin || !captionCityId}
              className="w-full py-2.5 rounded-xl border border-white/15 bg-white/10 text-white text-sm font-bold flex items-center justify-center gap-2 hover:bg-white/15 disabled:opacity-55 transition-colors"
            >
              <Sparkles className="h-4 w-4" />
              {loadingCheckin ? '准备中…' : '生成赛博打卡'}
            </button>

            {error && <div className="text-xs text-red-200 bg-red-500/12 border border-red-300/20 rounded-lg px-3 py-2">{error}</div>}
          </div>
        </div>
      </div>

      <div className="absolute bottom-12 left-1/2 -translate-x-1/2 z-10 text-white text-center pointer-events-none">
        <h2 className="text-5xl font-bold tracking-[0.2em] shadow-black drop-shadow-2xl">{targetName}</h2>
        <p className="mt-4 text-white/80 tracking-widest text-sm font-medium bg-black/20 backdrop-blur-md px-4 py-1 rounded-full border border-white/10">
          鼠标拖拽漫游 · 仿佛身临其境
        </p>
      </div>

      <SoundscapeControl url={streetTarget?.soundscape_url} label={targetName} />

      {checkinCity && (
        <CheckinFlow
          city={checkinCity}
          initialLandmarkId={streetTarget?.city_id ? streetTarget.id : undefined}
          sceneFile={sceneFile}
          onClose={() => { setCheckinCity(null); setSceneFile(null); }}
          onAchievementUnlocked={(achievements) => setUnlockedAchievements(achievements)}
        />
      )}

      {unlockedAchievements.length > 0 && (
        <AchievementUnlock
          achievements={unlockedAchievements}
          onClose={() => setUnlockedAchievements([])}
        />
      )}

      <GuessChallengeModal
        isOpen={showChallengeModal}
        onClose={() => setShowChallengeModal(false)}
        shotUrl={shotUrl}
        caption={caption}
        targetName={targetName}
        cityId={captionCityId}
        userId={userId}
      />
    </div>
  );
};
