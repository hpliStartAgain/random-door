import React, { useEffect, useRef, useState } from 'react';
import AMapLoader from '@amap/amap-jsapi-loader';
import { useMapStore } from '../store/useMapStore';
import { useCityStore } from '../store/useCityStore';
import { useGameStore } from '../store/useGameStore';
import { useViewStore } from '../store/useViewStore';

export const MapCanvas: React.FC = () => {
  const { setMapContext, flyTo } = useMapStore();
  const { cities } = useCityStore();
  const { lastRoll, fromPoint } = useGameStore();
  const { rollPhase } = useViewStore();
  const mapContainer = useRef<HTMLDivElement>(null);
  const [map, setMap] = useState<any>(null);
  const [error, setError] = useState<string>('');
  const flightLayersRef = useRef<any[]>([]);

  useEffect(() => {
    const key = import.meta.env.VITE_AMAP_KEY;
    const securityCode = import.meta.env.VITE_AMAP_SECURITY_CODE;

    if (!key) {
      setError('未配置 VITE_AMAP_KEY。');
      return;
    }

    if (securityCode) {
      (window as any)._AMapSecurityConfig = { securityJsCode: securityCode };
    }

    let mapInstance: any;

    AMapLoader.load({
      key: key,
      version: '2.0',
      plugins: ['AMap.Scale', 'AMap.ToolBar', 'AMap.ControlBar'],
    })
      .then((AMap) => {
        if (!mapContainer.current) return;
        mapInstance = new AMap.Map(mapContainer.current, {
          viewMode: '3D',
          zoom: 4.5,
          center: [104.195397, 35.86166], 
          mapStyle: 'amap://styles/whitesmoke', 
          pitch: 35, 
        });
        setMap(mapInstance);
        setMapContext(mapInstance, AMap);
      })
      .catch((e) => {
        console.error(e);
        setError('地图加载失败，请检查网络或密钥配置是否正确。');
      });

    return () => {
      if (mapInstance) mapInstance.destroy();
    };
  }, []);

  useEffect(() => {
    const { mapInstance, AMap } = useMapStore.getState();
    if (!mapInstance || !AMap || cities.length === 0) return;

    mapInstance.clearMap();

    cities.forEach(city => {
      const markerContent = `
        <div class="px-3 py-1 bg-[#F5F3EB]/90 backdrop-blur border border-[#E5E0D5] text-[#2B3A36] rounded-full text-xs font-bold shadow-md cursor-pointer hover:bg-white transition-colors flex items-center gap-1">
          <span class="w-2 h-2 rounded-full bg-[#C29F60]"></span>
          ${city.name}
        </div>
      `;
      
      const marker = new AMap.Marker({
        position: [city.lng, city.lat],
        content: markerContent,
        offset: new AMap.Pixel(-20, -15)
      });
      
      marker.on('click', () => {
         flyTo(city.lng, city.lat, 13);
      });
      
      mapInstance.add(marker);
    });
  }, [map, cities, flyTo]);

  useEffect(() => {
    const { mapInstance, AMap } = useMapStore.getState();
    if (!mapInstance || !AMap) return;
    flightLayersRef.current.forEach((layer) => { try { mapInstance.remove(layer); } catch {} });
    flightLayersRef.current = [];
    if (rollPhase !== 'flying' || !lastRoll || !fromPoint) return;

    const polyline = new AMap.Polyline({
      path: [[fromPoint.lng, fromPoint.lat], [lastRoll.target_point.lng, lastRoll.target_point.lat]],
      strokeColor: '#818cf8', strokeWeight: 2, strokeOpacity: 0.85,
      strokeStyle: 'dashed', strokeDasharray: [12, 6], zIndex: 50,
    });
    mapInstance.add(polyline);
    flightLayersRef.current.push(polyline);

    const markerDiv = document.createElement('div');
    markerDiv.style.cssText = 'position:relative;width:40px;height:40px;display:flex;align-items:center;justify-content:center;';
    markerDiv.innerHTML = `
      <div style="position:absolute;width:40px;height:40px;border-radius:50%;background:rgba(129,140,248,0.2);animation:mapLandPulse 1s ease-out infinite;"></div>
      <div style="position:absolute;width:20px;height:20px;border-radius:50%;background:rgba(129,140,248,0.35);animation:mapLandPulse 1s ease-out 0.3s infinite;"></div>
      <div style="position:relative;width:10px;height:10px;border-radius:50%;background:#818cf8;border:2px solid white;"></div>
    `;
    const landingMarker = new AMap.Marker({
      position: [lastRoll.target_city.lng, lastRoll.target_city.lat],
      content: markerDiv,
      offset: new AMap.Pixel(-20, -20),
      zIndex: 200,
    });
    mapInstance.add(landingMarker);
    flightLayersRef.current.push(landingMarker);

    mapInstance.setPitch(15, true, 600);
    setTimeout(() => {
      mapInstance.setZoomAndCenter(4, [104.195397, 35.86166], false, 700);
      setTimeout(() => {
        mapInstance.setPitch(60, true, 1800);
        mapInstance.setZoomAndCenter(11, [lastRoll.target_city.lng, lastRoll.target_city.lat], false, 2200);
      }, 800);
    }, 300);
  }, [rollPhase, lastRoll, fromPoint]);

  return (
    <div className="absolute inset-0 z-0 bg-background overflow-hidden">
      <style>{`
        @keyframes mapLandPulse {
          0%   { transform: scale(1); opacity: 0.8; }
          100% { transform: scale(2.8); opacity: 0; }
        }
      `}</style>
      {error ? (
        <div className="w-full h-full flex flex-col items-center justify-center bg-card text-muted-foreground p-10 text-center z-10 relative">
          <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center mb-4 border border-border">🗺️</div>
          <h3 className="text-xl font-bold mb-2 text-foreground">高德地图模块就绪</h3>
          <p>{error}</p>
        </div>
      ) : (
        <div ref={mapContainer} className="w-full h-full z-10 relative" />
      )}
      
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-primary/[0.03] font-bold text-[12rem] pointer-events-none select-none tracking-widest z-0 whitespace-nowrap">
        任意门
      </div>
    </div>
  );
};
