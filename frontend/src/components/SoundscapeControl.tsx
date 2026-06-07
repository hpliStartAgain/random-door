import React, { useEffect, useRef, useState } from 'react';
import { Pause, Play, Volume2, VolumeX } from 'lucide-react';

interface Props {
  url?: string;
  label: string;
}

export const SoundscapeControl: React.FC<Props> = ({ url, label }) => {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [playing, setPlaying] = useState(false);
  const [volume, setVolume] = useState(0.45);
  const [error, setError] = useState(false);

  useEffect(() => {
    setPlaying(false);
    setError(false);
  }, [url]);

  if (!url) return null;

  const handleToggle = async () => {
    const audio = audioRef.current;
    if (!audio) return;
    setError(false);
    if (playing) {
      audio.pause();
      setPlaying(false);
      return;
    }
    try {
      audio.volume = volume;
      await audio.play();
      setPlaying(true);
    } catch {
      setError(true);
      setPlaying(false);
    }
  };

  return (
    <div className="absolute left-6 bottom-6 z-10 pointer-events-auto max-w-[calc(100vw-3rem)]">
      <audio
        ref={audioRef}
        src={url}
        loop
        onEnded={() => setPlaying(false)}
        onError={() => { setError(true); setPlaying(false); }}
      />
      <div className="rounded-2xl border border-white/15 bg-black/45 backdrop-blur-xl text-white shadow-2xl px-3 py-2 flex items-center gap-3">
        <button
          onClick={handleToggle}
          className="h-9 w-9 rounded-full bg-white text-[#22302C] flex items-center justify-center hover:bg-white/90 transition-colors"
          aria-label={playing ? '暂停声景' : '播放声景'}
        >
          {playing ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4 ml-0.5" />}
        </button>
        <div className="min-w-0">
          <div className="text-xs font-semibold truncate">{label}</div>
          <div className="text-[11px] text-white/65">{error ? '声景加载失败' : '地标声景'}</div>
        </div>
        <div className="hidden sm:flex items-center gap-1.5">
          {volume > 0 ? <Volume2 className="h-3.5 w-3.5 text-white/70" /> : <VolumeX className="h-3.5 w-3.5 text-white/70" />}
          <input
            type="range"
            min={0}
            max={1}
            step={0.05}
            value={volume}
            onChange={(e) => {
              const next = Number(e.target.value);
              setVolume(next);
              if (audioRef.current) audioRef.current.volume = next;
            }}
            className="w-20 accent-white"
            aria-label="声景音量"
          />
        </div>
      </div>
    </div>
  );
};
