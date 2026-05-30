import React from 'react';
import { Pannellum } from 'pannellum-react';
import { useViewStore } from '../store/useViewStore';

export const StreetViewCanvas: React.FC = () => {
  const { setCanvasMode, streetTarget } = useViewStore();

  const panoramaUrl = streetTarget?.panoramaUrl || 'https://pannellum.org/images/cerro-toco-0.jpg';

  return (
    <div className="absolute inset-0 z-0 bg-black flex flex-col overflow-hidden">
      
      <div className="absolute inset-0 z-0 opacity-90">
        <Pannellum
          width="100%"
          height="100%"
          image={panoramaUrl}
          pitch={10}
          yaw={180}
          hfov={110}
          autoLoad
          showZoomCtrl={false}
          showFullscreenCtrl={false}
          mouseZoom={false}
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
          <span>&larr;</span> 退出 3D 全景
        </button>
      </div>

      <div className="absolute bottom-12 left-1/2 -translate-x-1/2 z-10 text-white text-center pointer-events-none">
        <h2 className="text-5xl font-bold tracking-[0.2em] shadow-black drop-shadow-2xl">{streetTarget?.name || '未知异境'}</h2>
        <p className="mt-4 text-white/80 tracking-widest text-sm font-medium bg-black/20 backdrop-blur-md px-4 py-1 rounded-full border border-white/10">
          鼠标拖拽漫游 · 仿佛身临其境
        </p>
      </div>
    </div>
  );
};
