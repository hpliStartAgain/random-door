import React, { useState, useEffect, useRef } from 'react';
import { api } from '../api';
import { useCityStore } from '../store/useCityStore';
import type { CityDetail } from '../api/types';
import { useToastStore } from '../store/useToastStore';

interface UploadBtnProps {
  label: string;
  currentUrl?: string | null;
  onUpload: (file: File) => Promise<unknown>;
  onBindUrl?: (url: string) => Promise<unknown>;
}

const UploadBtn: React.FC<UploadBtnProps> = ({ label, currentUrl, onUpload, onBindUrl }) => {
  const [uploading, setUploading] = useState(false);
  const [showUrlInput, setShowUrlInput] = useState(false);
  const [urlValue, setUrlValue] = useState('');
  const ref = useRef<HTMLInputElement>(null);

  const handleChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setUploading(true);
    try { await onUpload(file); } finally { setUploading(false); e.target.value = ''; }
  };

  const handleBindUrl = async () => {
    if (!urlValue.trim() || !onBindUrl) return;
    setUploading(true);
    try { await onBindUrl(urlValue.trim()); setShowUrlInput(false); setUrlValue(''); }
    finally { setUploading(false); }
  };

  return (
    <div className="flex flex-col gap-1 mt-1">
      <div className="flex items-center gap-2 flex-wrap">
        {currentUrl && (
          <img src={currentUrl} alt="" className="w-10 h-10 object-cover rounded-lg border border-border shrink-0" />
        )}
        <button
          onClick={() => ref.current?.click()}
          disabled={uploading}
          className="text-xs px-3 py-1.5 rounded-lg bg-primary/10 hover:bg-primary/20 border border-primary/20 text-primary font-medium transition-colors disabled:opacity-50"
        >
          {uploading ? '处理中…' : label}
        </button>
        {onBindUrl && (
          <button
            onClick={() => setShowUrlInput(v => !v)}
            disabled={uploading}
            className="text-xs px-3 py-1.5 rounded-lg bg-secondary hover:bg-border border border-border text-muted-foreground font-medium transition-colors disabled:opacity-50"
          >
            🔗 粘贴URL
          </button>
        )}
        <input ref={ref} type="file" accept="image/*" className="hidden" onChange={handleChange} />
      </div>
      {showUrlInput && onBindUrl && (
        <div className="flex gap-2 mt-1">
          <input
            type="url"
            value={urlValue}
            onChange={e => setUrlValue(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleBindUrl()}
            placeholder="https://example.com/image.jpg"
            className="flex-1 text-xs px-3 py-1.5 rounded-lg border border-border bg-background outline-none focus:border-primary/50 min-w-0"
          />
          <button
            onClick={handleBindUrl}
            disabled={!urlValue.trim() || uploading}
            className="text-xs px-3 py-1.5 rounded-lg bg-primary text-primary-foreground font-medium disabled:opacity-50"
          >
            确认
          </button>
        </div>
      )}
    </div>
  );
};

interface Props {
  onClose: () => void;
}

export const AdminPage: React.FC<Props> = ({ onClose }) => {
  const { cities, loadCities } = useCityStore();
  const { push: pushToast } = useToastStore();
  const [token, setToken] = useState('');
  const [authed, setAuthed] = useState(false);
  const [cityDetails, setCityDetails] = useState<Record<number, CityDetail>>({});
  const [expandedCity, setExpandedCity] = useState<number | null>(null);

  useEffect(() => { loadCities(); }, [loadCities]);

  const handleAuth = () => {
    if (!token.trim()) return;
    setAuthed(true);
  };

  const loadCityDetail = async (cityId: number) => {
    if (cityDetails[cityId]) return;
    try {
      const detail = await api.getCityDetail(cityId);
      setCityDetails(prev => ({ ...prev, [cityId]: detail }));
    } catch { pushToast('加载城市详情失败', 'error'); }
  };

  const handleExpand = async (cityId: number) => {
    if (expandedCity === cityId) { setExpandedCity(null); return; }
    setExpandedCity(cityId);
    await loadCityDetail(cityId);
  };

  if (!authed) {
    return (
      <div className="fixed inset-0 z-[100] flex items-center justify-center bg-background/95 backdrop-blur">
        <div className="w-full max-w-sm mx-4 bg-card border border-border rounded-3xl p-8 shadow-2xl">
          <div className="text-center mb-6">
            <div className="text-3xl mb-2">🔐</div>
            <h2 className="text-xl font-bold text-foreground">后台管理</h2>
            <p className="text-sm text-muted-foreground mt-1">输入 ADMIN_TOKEN 以继续</p>
          </div>
          <input
            type="password"
            value={token}
            onChange={e => setToken(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleAuth()}
            placeholder="Admin Token"
            className="w-full px-4 py-3 rounded-xl bg-secondary border border-border text-sm outline-none focus:border-primary/50 mb-4"
            autoFocus
          />
          <div className="flex gap-3">
            <button onClick={onClose} className="flex-1 py-2.5 rounded-xl border border-border text-sm font-semibold hover:bg-secondary transition-colors">取消</button>
            <button onClick={handleAuth} className="flex-1 py-2.5 rounded-xl bg-primary text-primary-foreground text-sm font-bold transition-colors">进入</button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 z-[100] bg-background flex flex-col">
      <div className="h-14 flex items-center justify-between px-6 border-b border-border/40 shrink-0">
        <h1 className="font-bold text-lg">🛠 媒体资产管理</h1>
        <button onClick={onClose} className="text-sm text-muted-foreground hover:text-foreground transition-colors">✕ 关闭</button>
      </div>

      <div className="flex-1 overflow-y-auto px-6 py-5 space-y-3 max-w-2xl mx-auto w-full">
        {cities.map(city => (
          <div key={city.id} className="border border-border rounded-2xl overflow-hidden">
            <button
              onClick={() => handleExpand(city.id)}
              className="w-full flex items-center gap-3 p-4 hover:bg-secondary/50 transition-colors text-left"
            >
              {city.cover_image_url
                ? <img src={city.cover_image_url} alt={city.name} className="w-12 h-12 object-cover rounded-xl shrink-0" />
                : <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center text-xl shrink-0">🏙</div>
              }
              <div className="flex-1 min-w-0">
                <div className="font-semibold text-sm">{city.name}</div>
                <div className="text-xs text-muted-foreground">{city.province}</div>
              </div>
              <span className="text-muted-foreground text-sm">{expandedCity === city.id ? '▲' : '▼'}</span>
            </button>

            {expandedCity === city.id && (
              <div className="px-4 pb-5 pt-1 border-t border-border/40 space-y-5 bg-secondary/20">
                <div>
                  <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">封面图</div>
                  <UploadBtn
                    label="上传封面"
                    currentUrl={city.cover_image_url}
                    onUpload={async (file) => {
                      const res = await api.adminUploadCityCover(city.id, file, token);
                      pushToast(`${city.name} 封面已更新`, 'success');
                      await loadCities();
                      return res;
                    }}
                    onBindUrl={async (url) => {
                      const res = await api.adminBindCityCoverURL(city.id, url, token);
                      pushToast(`${city.name} 封面 URL 已绑定`, 'success');
                      await loadCities();
                      return res;
                    }}
                  />
                </div>

                {cityDetails[city.id]?.landmarks?.length > 0 && (
                  <div>
                    <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">地标图片</div>
                    <div className="space-y-2">
                      {cityDetails[city.id].landmarks.map(lm => (
                        <div key={lm.id} className="flex items-center gap-3">
                          <span className="text-sm flex-1">{lm.name}</span>
                          <UploadBtn
                            label="上传"
                            currentUrl={lm.image_url}
                            onUpload={async (file) => {
                              const res = await api.adminUploadLandmarkImage(lm.id, file, token);
                              pushToast(`${lm.name} 图片已更新`, 'success');
                              api.getCityDetail(city.id).then(d => setCityDetails(prev => ({ ...prev, [city.id]: d })));
                              return res;
                            }}
                            onBindUrl={async (url) => {
                              const res = await api.adminBindLandmarkImageURL(lm.id, url, token);
                              pushToast(`${lm.name} URL 已绑定`, 'success');
                              api.getCityDetail(city.id).then(d => setCityDetails(prev => ({ ...prev, [city.id]: d })));
                              return res;
                            }}
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {cityDetails[city.id]?.characters?.length > 0 && (
                  <div>
                    <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">人物头像</div>
                    <div className="space-y-2">
                      {cityDetails[city.id].characters.map(ch => (
                        <div key={ch.id} className="flex items-center gap-3">
                          <span className="text-sm flex-1">{ch.name}</span>
                          <UploadBtn
                            label="上传"
                            currentUrl={ch.avatar_url}
                            onUpload={async (file) => {
                              const res = await api.adminUploadCharacterAvatar(ch.id, file, token);
                              pushToast(`${ch.name} 头像已更新`, 'success');
                              api.getCityDetail(city.id).then(d => setCityDetails(prev => ({ ...prev, [city.id]: d })));
                              return res;
                            }}
                            onBindUrl={async (url) => {
                              const res = await api.adminBindCharacterAvatarURL(ch.id, url, token);
                              pushToast(`${ch.name} 头像 URL 已绑定`, 'success');
                              api.getCityDetail(city.id).then(d => setCityDetails(prev => ({ ...prev, [city.id]: d })));
                              return res;
                            }}
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {cityDetails[city.id]?.foods?.length > 0 && (
                  <div>
                    <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">美食图片</div>
                    <div className="space-y-2">
                      {cityDetails[city.id].foods.map(f => (
                        <div key={f.id} className="flex items-center gap-3">
                          <span className="text-sm flex-1">{f.name}</span>
                          <UploadBtn
                            label="上传"
                            currentUrl={f.image_url}
                            onUpload={async (file) => {
                              const res = await api.adminUploadFoodImage(f.id, file, token);
                              pushToast(`${f.name} 图片已更新`, 'success');
                              api.getCityDetail(city.id).then(d => setCityDetails(prev => ({ ...prev, [city.id]: d })));
                              return res;
                            }}
                            onBindUrl={async (url) => {
                              const res = await api.adminBindFoodImageURL(f.id, url, token);
                              pushToast(`${f.name} URL 已绑定`, 'success');
                              api.getCityDetail(city.id).then(d => setCityDetails(prev => ({ ...prev, [city.id]: d })));
                              return res;
                            }}
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};
