import React, { useCallback, useEffect, useRef, useState } from 'react';
import AMapLoader from '@amap/amap-jsapi-loader';
import { api } from '../api';
import type { CityDetail, CityListItem, Landmark, LandmarkMapItem } from '../api/types';
import { useMapStore } from '../store/useMapStore';
import { useCityStore } from '../store/useCityStore';
import { useGameStore } from '../store/useGameStore';
import { useUserStore } from '../store/useUserStore';
import { useViewStore } from '../store/useViewStore';

const DEFAULT_USER_POSITION = { lat: 39.9042, lng: 116.4074 };

const FOX_SVG = `
  <img src="/icon-transparent.png" alt="我" style="position:relative;width:38px;height:38px;object-fit:contain;filter:drop-shadow(0 8px 12px rgba(43,58,54,0.28));" />
`;

type MapCity = Pick<CityListItem, 'id' | 'name' | 'province' | 'lat' | 'lng'> & {
  landmarks?: Array<Landmark | LandmarkMapItem>;
};

function escapeHTML(value: string): string {
  return value.replace(/[&<>"']/g, (char) => {
    switch (char) {
      case '&': return '&amp;';
      case '<': return '&lt;';
      case '>': return '&gt;';
      case '"': return '&quot;';
      case "'": return '&#39;';
      default: return char;
    }
  });
}

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
  const safeName = escapeHTML(name);
  if (isActive) {
    return `
      <div style="position:relative;display:flex;flex-direction:column;align-items:center;">
        <div style="position:absolute;top:-4px;left:50%;transform:translateX(-50%);width:36px;height:36px;border-radius:50%;background:rgba(43,58,54,0.18);animation:mapLandPulse 1.2s ease-out infinite;"></div>
        <div class="px-3 py-1 backdrop-blur rounded-full text-xs font-bold shadow-md cursor-pointer flex items-center gap-1 transition-colors" style="background:#2B3A36;color:#fff;border:1.5px solid #2B3A36;z-index:1;">
          <span style="width:6px;height:6px;border-radius:50%;background:#C29F60;display:inline-block;flex-shrink:0;"></span>
          ${safeName}
        </div>
      </div>
    `;
  }
  return `
    <div class="px-3 py-1 bg-[#F5F3EB]/90 backdrop-blur border border-[#E5E0D5] text-[#2B3A36] rounded-full text-xs font-bold shadow-md cursor-pointer hover:bg-white transition-colors flex items-center gap-1">
      <span class="w-2 h-2 rounded-full bg-[#C29F60]"></span>
      ${safeName}
    </div>
  `;
}

function landmarkMarkerContent(name: string, emphasized: boolean, showLabel: boolean): string {
  const safeName = escapeHTML(name);
  const label = showLabel
    ? `<span style="white-space:nowrap;font-size:11px;font-weight:800;line-height:1;">${safeName}</span>`
    : '';
  return `
    <div style="position:relative;display:flex;align-items:center;gap:6px;padding:6px ${showLabel ? '9px' : '6px'};border-radius:999px;background:${emphasized ? '#2B3A36' : 'rgba(245,243,235,0.94)'};color:${emphasized ? '#fff' : '#2B3A36'};border:1px solid ${emphasized ? '#2B3A36' : 'rgba(229,224,213,0.95)'};box-shadow:0 10px 22px rgba(43,58,54,0.18);cursor:pointer;">
      <span style="width:10px;height:10px;border-radius:3px;background:#C29F60;box-shadow:0 0 0 3px ${emphasized ? 'rgba(194,159,96,0.22)' : 'rgba(194,159,96,0.16)'};display:inline-block;flex-shrink:0;"></span>
      ${label}
    </div>
  `;
}

function hasCoordinate(item: Landmark | LandmarkMapItem): item is (Landmark | LandmarkMapItem) & { lat: number; lng: number } {
  return typeof item.lat === 'number' && Number.isFinite(item.lat) &&
    typeof item.lng === 'number' && Number.isFinite(item.lng);
}

export const MapCanvas: React.FC = () => {
  const { setMapContext, flyTo } = useMapStore();
  const { cities, filteredCities, searchQuery, cityCache, loadCity } = useCityStore();
  const { lastRoll, fromPoint } = useGameStore();
  const { rollPhase, activeCityId, setView, openDrawer } = useViewStore();
  const { userId, setCurrentCityId } = useUserStore();
  const mapContainer = useRef<HTMLDivElement>(null);
  const [map, setMap] = useState<any>(null);
  const [mapZoom, setMapZoom] = useState(4.5);
  const [error, setError] = useState<string>('');
  const [userPosition, setUserPosition] = useState(DEFAULT_USER_POSITION);
  const baseLayersRef = useRef<any[]>([]);
  const flightLayersRef = useRef<any[]>([]);

  const recordFreeVisit = useCallback((cityId: number) => {
    if (!userId || activeCityId === cityId) return;
    api.createFreeVisit(userId, cityId).catch((err) => {
      console.error('record free visit failed', err);
    });
  }, [activeCityId, userId]);

  const activateCity = useCallback((city: MapCity, zoom = 13) => {
    flyTo(city.lng, city.lat, zoom);
    setCurrentCityId(city.id);
    setView('CITY_DETAIL', city.id);
    recordFreeVisit(city.id);
    loadCity(city.id).catch((err) => {
      console.error('load city failed', err);
    });
  }, [flyTo, loadCity, recordFreeVisit, setCurrentCityId, setView]);

  const activateLandmark = useCallback(async (city: MapCity, landmark: Landmark | LandmarkMapItem) => {
    if (!hasCoordinate(landmark)) return;
    flyTo(landmark.lng, landmark.lat, 16);
    setCurrentCityId(city.id);
    setView('CITY_DETAIL', city.id);
    recordFreeVisit(city.id);

    let detail: CityDetail | undefined = cityCache[city.id];
    if (!detail) {
      try {
        detail = await loadCity(city.id);
      } catch (err) {
        console.error('load city failed', err);
      }
    }
    const fullLandmark = detail?.landmarks.find((item) => item.id === landmark.id) ?? landmark;
    openDrawer('gallery', { ...fullLandmark, target_type: 'landmark' });
  }, [cityCache, flyTo, loadCity, openDrawer, recordFreeVisit, setCurrentCityId, setView]);

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
    if (!map) return;
    const updateZoom = () => {
      try {
        setMapZoom(map.getZoom());
      } catch {
        setMapZoom(4.5);
      }
    };
    updateZoom();
    map.on?.('zoomend', updateZoom);
    return () => {
      try { map.off?.('zoomend', updateZoom); } catch {}
    };
  }, [map]);

  useEffect(() => {
    if (!activeCityId || cityCache[activeCityId]) return;
    loadCity(activeCityId).catch((err) => {
      console.error('load active city failed', err);
    });
  }, [activeCityId, cityCache, loadCity]);

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
         activateCity(city, 13);
      });
      
      mapInstance.add(marker);
      baseLayersRef.current.push(marker);
    });

    const activeCity = activeCityId
      ? (cityCache[activeCityId] ?? cities.find((city) => city.id === activeCityId))
      : undefined;
    const landmarkCities: MapCity[] = mapZoom >= 8
      ? citiesToShow
      : activeCity ? [activeCity] : [];

    landmarkCities.forEach((city) => {
      city.landmarks?.forEach((landmark) => {
        if (!hasCoordinate(landmark)) return;
        const isActiveCity = city.id === activeCityId;
        const marker = new AMap.Marker({
          position: [landmark.lng, landmark.lat],
          content: landmarkMarkerContent(landmark.name, isActiveCity, isActiveCity || mapZoom >= 10),
          offset: new AMap.Pixel(-14, -14),
          zIndex: isActiveCity ? 170 : 120,
        });

        marker.on('click', () => {
          void activateLandmark(city, landmark);
        });

        mapInstance.add(marker);
        baseLayersRef.current.push(marker);
      });
    });

    const userMarker = new AMap.Marker({
      position: [userPosition.lng, userPosition.lat],
      content: foxMarkerContent('当前位置', true),
      offset: new AMap.Pixel(-23, -23),
      zIndex: 180,
    });
    mapInstance.add(userMarker);
    baseLayersRef.current.push(userMarker);
  }, [map, searchQuery, filteredCities, userPosition, activeCityId, cityCache, cities, mapZoom, activateCity, activateLandmark]);

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
