import React from 'react';

interface Props {
  imageUrl: string;
  cityName: string;
  landmarkName: string;
}

export const CheckinPoster: React.FC<Props> = ({ imageUrl, cityName, landmarkName }) => {
  const handleDownload = () => {
    const img = new Image();
    img.crossOrigin = 'anonymous';
    img.onload = () => {
      const BRAND_H = Math.round(img.height * 0.08);
      const canvas = document.createElement('canvas');
      canvas.width = img.width;
      canvas.height = img.height + BRAND_H;
      const ctx = canvas.getContext('2d');
      if (!ctx) return;
      ctx.drawImage(img, 0, 0);
      ctx.fillStyle = '#2B3A36';
      ctx.fillRect(0, img.height, img.width, BRAND_H);
      const fontSize = Math.max(12, Math.round(img.width / 28));
      ctx.font = `bold ${fontSize}px sans-serif`;
      ctx.fillStyle = 'rgba(255,255,255,0.85)';
      ctx.textBaseline = 'middle';
      ctx.fillText(`任意门 · ${cityName} · ${landmarkName}`, img.width * 0.04, img.height + BRAND_H / 2);
      const a = document.createElement('a');
      a.download = `random-door-${cityName}.png`;
      a.href = canvas.toDataURL('image/png');
      a.click();
    };
    img.src = imageUrl;
  };

  return (
    <div className="rounded-2xl overflow-hidden border border-border shadow-lg">
      <img src={imageUrl} alt="赛博打卡" className="w-full object-cover" />
      <div className="bg-primary px-4 py-3 flex items-center justify-between">
        <div>
          <span className="text-white/80 text-xs font-bold tracking-widest">任意门</span>
          <span className="text-white/40 text-xs ml-2">{cityName} · {landmarkName}</span>
        </div>
        <button
          onClick={handleDownload}
          className="text-xs px-3 py-1.5 bg-white/10 hover:bg-white/20 border border-white/20 text-white rounded-full transition-colors"
        >
          ↓ 保存
        </button>
      </div>
    </div>
  );
};
