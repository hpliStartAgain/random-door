import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

export const CheckinPage: React.FC = () => {
  const { id: _id } = useParams();
  const navigate = useNavigate();
  const [generating, setGenerating] = useState(false);

  const handleGenerate = () => {
    setGenerating(true);
    setTimeout(() => {
      setGenerating(false);
      alert('打卡成功！获得成就：古都初见');
    }, 2000);
  };

  return (
    <div className="min-h-screen bg-background p-4 flex flex-col items-center justify-center">
      <header className="absolute top-4 left-4">
        <div className="glass-panel inline-block px-4 py-2 rounded-full font-bold">
          <span className="text-primary cursor-pointer" onClick={() => navigate(-1)}>&lt; 返回</span>
        </div>
      </header>

      <div className="glass-panel p-8 rounded-3xl max-w-md w-full space-y-6 text-center">
        <h2 className="text-2xl font-bold">赛博打卡</h2>
        <p className="text-muted-foreground text-sm">上传你的自拍，生成带有城市地标背景的打卡照！</p>
        
        <div className="aspect-[3/4] w-full border-2 border-dashed border-border rounded-xl flex items-center justify-center text-muted-foreground cursor-pointer hover:bg-white/50 transition-colors">
          点击上传照片 (≤5MB)
        </div>

        <button 
          onClick={handleGenerate}
          disabled={generating}
          className="w-full py-3 bg-destructive text-destructive-foreground font-semibold rounded-xl disabled:opacity-50"
        >
          {generating ? '生成中...' : '生成并打卡'}
        </button>
      </div>
    </div>
  );
};
