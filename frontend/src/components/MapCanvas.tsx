import React, { useEffect, useRef, useState } from 'react';
import AMapLoader from '@amap/amap-jsapi-loader';
import { useMapStore } from '../store/useMapStore';
import { useCityStore } from '../store/useCityStore';

export const MapCanvas: React.FC = () => {
  const { setMapContext, flyTo } = useMapStore();
  const { cities } = useCityStore();
  const mapContainer = useRef<HTMLDivElement>(null);
  const [map, setMap] = useState<any>(null);
  const [error, setError] = useState<string>('');

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

  return (
    <div className="absolute inset-0 z-0 bg-background overflow-hidden">
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
