import React, { useEffect, useRef, useState } from 'react';
import AMapLoader from '@amap/amap-jsapi-loader';
import { useMapStore } from '../store/useMapStore';
import { useCityStore } from '../store/useCityStore';
import { useGameStore } from '../store/useGameStore';
import { useViewStore } from '../store/useViewStore';

const DEFAULT_USER_POSITION = { lat: 39.9042, lng: 116.4074 };

const FOX_SVG = `
  <img src="/icon-transparent.png" alt="我" style="position:relative;width:38px;height:38px;object-fit:contain;filter:drop-shadow(0 8px 12px rgba(43,58,54,0.28));" />
`;

function foxMarkerContent(label?: string, pulse = false): HTMLDivElement {
  const markerDiv = document.createElement('div');
  markerDiv.style.cssText = 'position:relative;width:46px;height:46px;display:flex;align-items:center;justify-content:center;';
  markerDiv.innerHTML = `
    ${pulse ? '<div style="position:absolute;width:40px;height:40px;border-radius:50%;background:rgba(194,159,96,0.22);animation:mapLandPulse 1s ease-out infinite;"></div>' : ''}
    ${pulse ? '<div style="position:absolute;width:20px;height:20px;border-radius:50%;background:rgba(194,159,96,0.38);animation:mapLandPulse 1s ease-out 0.3s infinite;"></div>' : ''}
    ${FOX_SVG}
    ${label ? `<div style="position:absolute;top:35px;left:50%;transform:translateX(-50%);white-space:nowrap;border:1px solid rgba(229,224,213,0.9);border-radius:999px;background:rgba(245,243,235,0.92);color:#2B3A36;font-size:10px;font-weight:700;padding:2px 7px;box-shadow:0 8px 18px rgba(43,58,54,0.12);">${label}</div>` : ''}
  `;
  return markerDiv;
}

function cityMarkerContent(name: string, isActive: boolean): string {
  if (isActive) {
    return `
      <div style="position:relative;display:flex;flex-direction:column;align-items:center;">
        <div style="position:absolute;top:-4px;left:50%;transform:translateX(-50%);width:36px;height:36px;border-radius:50%;background:rgba(43,58,54,0.18);animation:mapLandPulse 1.2s ease-out infinite;"></div>
        <div class="px-3 py-1 backdrop-blur rounded-full text-xs font-bold shadow-md cursor-pointer flex items-center gap-1 transition-colors" style="background:#2B3A36;color:#fff;border:1.5px solid #2B3A36;z-index:1;">
          <span style="width:6px;height:6px;border-radius:50%;background:#C29F60;display:inline-block;flex-shrink:0;"></span>
          ${name}
        </div>
      </div>
    `;
  }
  return `
    <div class="px-3 py-1 bg-[#F5F3EB]/90 backdrop-blur border border-[#E5E0D5] text-[#2B3A36] rounded-full text-xs font-bold shadow-md cursor-pointer hover:bg-white transition-colors flex items-center gap-1">
      <span class="w-2 h-2 rounded-full bg-[#C29F60]"></span>
      ${name}
    </div>
  `;
}

export const MapCanvas: React.FC = () => {
  const { setMapContext, flyTo } = useMapStore();
  const { filteredCities, searchQuery } = useCityStore();
  const { lastRoll, fromPoint } = useGameStore();
  const { rollPhase, activeCityId } = useViewStore();
  const mapContainer = useRef<HTMLDivElement>(null);
  const [map, setMap] = useState<any>(null);
  const [error, setError] = useState<string>('');
  const [userPosition, setUserPosition] = useState(DEFAULT_USER_POSITION);
  const baseLayersRef = useRef<any[]>([]);
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
    if (!navigator.geolocation) return;
    let cancelled = false;
    navigator.geolocation.getCurrentPosition(
      (position) => {
        if (cancelled) return;
        setUserPosition({
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        });
      },
      () => undefined,
      { enableHighAccuracy: false, timeout: 8000, maximumAge: 5 * 60 * 1000 },
    );
    return () => { cancelled = true; };
  }, []);

  useEffect(() => {
    const { mapInstance, AMap } = useMapStore.getState();
    if (!mapInstance || !AMap) return;

    baseLayersRef.current.forEach((layer) => { try { mapInstance.remove(layer); } catch {} });
    baseLayersRef.current = [];

    const citiesToShow = filteredCities();
    citiesToShow.forEach(city => {
      const isActive = city.id === activeCityId;
      const markerContent = cityMarkerContent(city.name, isActive);
      
      const marker = new AMap.Marker({
        position: [city.lng, city.lat],
        content: markerContent,
        offset: new AMap.Pixel(-20, -15),
        zIndex: isActive ? 160 : 100,
      });
      
      marker.on('click', () => {
         flyTo(city.lng, city.lat, 13);
      });
      
      mapInstance.add(marker);
      baseLayersRef.current.push(marker);
    });

    const userMarker = new AMap.Marker({
      position: [userPosition.lng, userPosition.lat],
      content: foxMarkerContent('当前位置', true),
      offset: new AMap.Pixel(-23, -23),
      zIndex: 180,
    });
    mapInstance.add(userMarker);
    baseLayersRef.current.push(userMarker);
  }, [map, searchQuery, filteredCities, flyTo, userPosition, activeCityId]);

  useEffect(() => {
    const { mapInstance, AMap } = useMapStore.getState();
    if (!mapInstance || !AMap) return;
    flightLayersRef.current.forEach((layer) => { try { mapInstance.remove(layer); } catch {} });
    flightLayersRef.current = [];
    if (rollPhase !== 'flying' || !lastRoll || !fromPoint) return;

    const polyline = new AMap.Polyline({
      path: [[fromPoint.lng, fromPoint.lat], [lastRoll.target_point.lng, lastRoll.target_point.lat]],
      strokeColor: '#C29F60', strokeWeight: 2, strokeOpacity: 0.9,
      strokeStyle: 'dashed', strokeDasharray: [12, 6], zIndex: 50,
    });
    mapInstance.add(polyline);
    flightLayersRef.current.push(polyline);

    const landingMarker = new AMap.Marker({
      position: [lastRoll.target_city.lng, lastRoll.target_city.lat],
      content: foxMarkerContent(undefined, true),
      offset: new AMap.Pixel(-23, -23),
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
