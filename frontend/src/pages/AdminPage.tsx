import React, { useState, useEffect, useRef } from 'react';
import { ClipboardPaste, Link2, Upload } from 'lucide-react';
import { api } from '../api';
import { useCityStore } from '../store/useCityStore';
import type { AdminCoverageResponse, CityDetail } from '../api/types';
import { useToastStore } from '../store/useToastStore';

const MAX_IMAGE_SIZE_BYTES = 5 * 1024 * 1024;
const ALLOWED_IMAGE_TYPES = new Set(['image/jpeg', 'image/png', 'image/webp']);

function imageExtension(mimeType: string): string {
  if (mimeType === 'image/jpeg') return 'jpg';
  if (mimeType === 'image/webp') return 'webp';
  return 'png';
}

function validateImageFile(file: File): string | null {
  if (!ALLOWED_IMAGE_TYPES.has(file.type)) {
    return '仅支持 JPG、PNG、WEBP 图片';
  }
  if (file.size > MAX_IMAGE_SIZE_BYTES) {
    return '图片不能超过 5MB';
  }
  return null;
}

interface UploadBtnProps {
  label: string;
  currentUrl?: string | null;
  onUpload: (file: File) => Promise<unknown>;
  onBindUrl?: (url: string) => Promise<unknown>;
}

const UploadBtn: React.FC<UploadBtnProps> = ({ label, currentUrl, onUpload, onBindUrl }) => {
  const { push: pushToast } = useToastStore();
  const [uploading, setUploading] = useState(false);
  const [showUrlInput, setShowUrlInput] = useState(false);
  const [urlValue, setUrlValue] = useState('');
  const ref = useRef<HTMLInputElement>(null);

  const uploadFile = async (file: File) => {
    const validationError = validateImageFile(file);
    if (validationError) {
      pushToast(validationError, 'error');
      return;
    }
    setUploading(true);
    try {
      await onUpload(file);
    } catch (err) {
      pushToast((err as { message?: string })?.message || '图片处理失败', 'error');
    } finally {
      setUploading(false);
    }
  };

  const handleChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) await uploadFile(file);
    e.target.value = '';
  };

  const handlePaste = async (event: React.ClipboardEvent<HTMLDivElement>) => {
    const imageItem = Array.from(event.clipboardData.items).find(
      (item) => item.kind === 'file' && item.type.startsWith('image/'),
    );
    if (!imageItem) return;

    const file = imageItem.getAsFile();
    if (!file) return;

    event.preventDefault();
    const pastedFile = new File(
      [file],
      file.name || `clipboard-${Date.now()}.${imageExtension(file.type)}`,
      { type: file.type, lastModified: Date.now() },
    );
    await uploadFile(pastedFile);
  };

  const handleReadClipboard = async () => {
    const clipboard = navigator.clipboard as Clipboard & { read?: () => Promise<ClipboardItem[]> };
    if (!clipboard?.read) {
      pushToast('当前浏览器不支持直接读取剪贴板图片', 'error');
      return;
    }

    try {
      const items = await clipboard.read();
      for (const item of items) {
        const imageType = item.types.find((type) => type.startsWith('image/'));
        if (!imageType) continue;
        const blob = await item.getType(imageType);
        await uploadFile(new File(
          [blob],
          `clipboard-${Date.now()}.${imageExtension(imageType)}`,
          { type: imageType, lastModified: Date.now() },
        ));
        return;
      }
      pushToast('剪贴板中没有可上传的图片', 'error');
    } catch (err) {
      pushToast((err as { message?: string })?.message || '读取剪贴板失败', 'error');
    }
  };

  const handleBindUrl = async () => {
    if (!urlValue.trim() || !onBindUrl) return;
    setUploading(true);
    try {
      await onBindUrl(urlValue.trim());
      setShowUrlInput(false);
      setUrlValue('');
    } catch (err) {
      pushToast((err as { message?: string })?.message || '图片导入失败', 'error');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="flex flex-col gap-1 mt-1" onPaste={handlePaste}>
      <div className="flex items-center gap-2 flex-wrap">
        {currentUrl && (
          <img src={currentUrl} alt="" className="w-10 h-10 object-cover rounded-lg border border-border shrink-0" />
        )}
        <button
          onClick={() => ref.current?.click()}
          disabled={uploading}
          className="text-xs px-3 py-1.5 rounded-lg bg-primary/10 hover:bg-primary/20 border border-primary/20 text-primary font-medium transition-colors disabled:opacity-50 inline-flex items-center gap-1.5"
        >
          <Upload className="h-3.5 w-3.5" />
          {uploading ? '处理中…' : label}
        </button>
        <button
          onClick={handleReadClipboard}
          disabled={uploading}
          className="text-xs px-3 py-1.5 rounded-lg bg-secondary hover:bg-border border border-border text-muted-foreground font-medium transition-colors disabled:opacity-50 inline-flex items-center gap-1.5"
        >
          <ClipboardPaste className="h-3.5 w-3.5" />
          粘贴图片
        </button>
        {onBindUrl && (
          <button
            onClick={() => setShowUrlInput(v => !v)}
            disabled={uploading}
            className="text-xs px-3 py-1.5 rounded-lg bg-secondary hover:bg-border border border-border text-muted-foreground font-medium transition-colors disabled:opacity-50 inline-flex items-center gap-1.5"
          >
            <Link2 className="h-3.5 w-3.5" />
            粘贴URL
          </button>
        )}
        <input ref={ref} type="file" accept="image/jpeg,image/png,image/webp" className="hidden" onChange={handleChange} />
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
  const { cities, loadCities, reloadCity } = useCityStore();
  const { push: pushToast } = useToastStore();
  const [token, setToken] = useState('');
  const [authed, setAuthed] = useState(false);
  const [authLoading, setAuthLoading] = useState(false);
  const [cityDetails, setCityDetails] = useState<Record<number, CityDetail>>({});
  const [cityDrafts, setCityDrafts] = useState<Record<number, { intro: string; tags: string; dialect_sample: string; dialect_explanation: string }>>({});
  const [coverage, setCoverage] = useState<AdminCoverageResponse | null>(null);
  const [expandedCity, setExpandedCity] = useState<number | null>(null);
  const [newDrafts, setNewDrafts] = useState<Record<number, {
    landmarkName: string;
    landmarkDescription: string;
    foodName: string;
    foodDescription: string;
    characterName: string;
  }>>({});

  useEffect(() => { loadCities(); }, [loadCities]);

  const handleAuth = async () => {
    if (!token.trim() || authLoading) return;
    setAuthLoading(true);
    try {
      const cov = await api.adminCoverage(token);
      setCoverage(cov);
      setAuthed(true);
    } catch {
      pushToast('Token 无效，请检查 ADMIN_TOKEN', 'error');
    } finally {
      setAuthLoading(false);
    }
  };

  const refreshCoverage = async () => {
    if (!token.trim()) return;
    try {
      setCoverage(await api.adminCoverage(token));
    } catch { pushToast('覆盖率加载失败，请检查 ADMIN_TOKEN', 'error'); }
  };

  const loadCityDetail = async (cityId: number) => {
    if (cityDetails[cityId]) return;
    try {
      const detail = await api.getCityDetail(cityId);
      setCityDetails(prev => ({ ...prev, [cityId]: detail }));
      setCityDrafts(prev => ({
        ...prev,
        [cityId]: {
          intro: detail.intro || '',
          tags: detail.tags?.join(',') || '',
          dialect_sample: detail.dialect_sample || '',
          dialect_explanation: detail.dialect_explanation || '',
        },
      }));
    } catch { pushToast('加载城市详情失败', 'error'); }
  };

  const handleExpand = async (cityId: number) => {
    if (expandedCity === cityId) { setExpandedCity(null); return; }
    setExpandedCity(cityId);
    await loadCityDetail(cityId);
  };

  const refreshCityDetail = async (cityId: number) => {
    const detail = await reloadCity(cityId);
    setCityDetails(prev => ({ ...prev, [cityId]: detail }));
    setCityDrafts(prev => ({
      ...prev,
      [cityId]: {
        intro: detail.intro || '',
        tags: detail.tags?.join(',') || '',
        dialect_sample: detail.dialect_sample || '',
        dialect_explanation: detail.dialect_explanation || '',
      },
    }));
  };

  const emptyNewDraft = {
    landmarkName: '',
    landmarkDescription: '',
    foodName: '',
    foodDescription: '',
    characterName: '',
  };
  const newDraftFor = (cityId: number) => newDrafts[cityId] || emptyNewDraft;
  const updateNewDraft = (cityId: number, field: keyof typeof emptyNewDraft, value: string) => {
    setNewDrafts(prev => ({
      ...prev,
      [cityId]: { ...emptyNewDraft, ...prev[cityId], [field]: value },
    }));
  };

  const saveCityDraft = async (cityId: number) => {
    const draft = cityDrafts[cityId];
    if (!draft) return;
    await api.adminUpdateCity(cityId, {
      intro: draft.intro,
      dialect_sample: draft.dialect_sample,
      dialect_explanation: draft.dialect_explanation,
      tags: draft.tags.split(',').map(t => t.trim()).filter(Boolean),
    }, token);
    pushToast('城市内容已更新', 'success');
    await Promise.all([loadCities(), refreshCityDetail(cityId), refreshCoverage()]);
  };

  const refreshAdminCity = async (cityId: number) => {
    await Promise.all([refreshCityDetail(cityId), refreshCoverage()]);
  };

  const coverageFor = (cityId: number) => coverage?.items.find(item => item.city_id === cityId);

  const createLandmark = async (cityId: number) => {
    const draft = newDraftFor(cityId);
    if (!draft.landmarkName.trim()) return;
    await api.adminCreateLandmark(cityId, {
      name: draft.landmarkName.trim(),
      description: draft.landmarkDescription.trim(),
    }, token);
    updateNewDraft(cityId, 'landmarkName', '');
    updateNewDraft(cityId, 'landmarkDescription', '');
    pushToast('地标已新增', 'success');
    await refreshAdminCity(cityId);
  };

  const createFood = async (cityId: number) => {
    const draft = newDraftFor(cityId);
    if (!draft.foodName.trim()) return;
    await api.adminCreateFood(cityId, {
      name: draft.foodName.trim(),
      description: draft.foodDescription.trim(),
    }, token);
    updateNewDraft(cityId, 'foodName', '');
    updateNewDraft(cityId, 'foodDescription', '');
    pushToast('美食已新增', 'success');
    await refreshAdminCity(cityId);
  };

  const createCharacter = async (cityId: number) => {
    const draft = newDraftFor(cityId);
    if (!draft.characterName.trim()) return;
    await api.adminCreateCharacter(cityId, { name: draft.characterName.trim() }, token);
    updateNewDraft(cityId, 'characterName', '');
    pushToast('人物已新增', 'success');
    await refreshAdminCity(cityId);
  };

  const editPOI = async (
    cityId: number,
    kind: 'landmark' | 'food',
    item: { id: number; name: string; description?: string },
  ) => {
    const name = window.prompt('名称', item.name);
    if (name === null) return;
    const description = window.prompt('描述', item.description || '');
    if (description === null) return;
    if (kind === 'landmark') {
      await api.adminUpdateLandmark(item.id, { name: name.trim(), description }, token);
      pushToast('地标已更新', 'success');
    } else {
      await api.adminUpdateFood(item.id, { name: name.trim(), description }, token);
      pushToast('美食已更新', 'success');
    }
    await refreshAdminCity(cityId);
  };

  const editCharacter = async (cityId: number, item: { id: number; name: string; character_type: string; dialect_style?: string }) => {
    const name = window.prompt('人物名称', item.name);
    if (name === null) return;
    const dialectStyle = window.prompt('表达风格', item.dialect_style || '');
    if (dialectStyle === null) return;
    await api.adminUpdateCharacter(item.id, { name: name.trim(), dialect_style: dialectStyle }, token);
    pushToast('人物已更新', 'success');
    await refreshAdminCity(cityId);
  };

  const deleteItem = async (cityId: number, kind: 'landmark' | 'food' | 'character', id: number, name: string) => {
    if (!window.confirm(`删除「${name}」？`)) return;
    if (kind === 'landmark') {
      await api.adminDeleteLandmark(id, token);
      pushToast('地标已删除', 'success');
    } else if (kind === 'food') {
      await api.adminDeleteFood(id, token);
      pushToast('美食已删除', 'success');
    } else {
      await api.adminDeleteCharacter(id, token);
      pushToast('人物已删除', 'success');
    }
    await refreshAdminCity(cityId);
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
            <button onClick={handleAuth} disabled={authLoading} className="flex-1 py-2.5 rounded-xl bg-primary text-primary-foreground text-sm font-bold transition-colors disabled:opacity-60">{authLoading ? '验证中…' : '进入'}</button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 z-[100] bg-background flex flex-col">
      <div className="h-14 flex items-center justify-between px-6 border-b border-border/40 shrink-0">
        <h1 className="font-bold text-lg">内容 CMS</h1>
        <button onClick={onClose} className="text-sm text-muted-foreground hover:text-foreground transition-colors">✕ 关闭</button>
      </div>

      <div className="flex-1 overflow-y-auto px-6 py-5 space-y-3 max-w-3xl mx-auto w-full">
        {coverage && (
          <div className="border border-border rounded-2xl p-4 bg-secondary/20">
            <div className="flex items-center justify-between gap-3">
              <div>
                <div className="text-sm font-bold">内容覆盖率</div>
                <div className="text-xs text-muted-foreground mt-0.5">
                  {coverage.complete_cities} / {coverage.total_cities} 城完整
                </div>
              </div>
              <button
                onClick={refreshCoverage}
                className="text-xs px-3 py-1.5 rounded-lg border border-border hover:bg-secondary transition-colors"
              >
                刷新
              </button>
            </div>
            <div className="mt-3 h-2 bg-border rounded-full overflow-hidden">
              <div
                className="h-full bg-primary rounded-full"
                style={{ width: `${coverage.total_cities ? (coverage.complete_cities / coverage.total_cities) * 100 : 0}%` }}
              />
            </div>
          </div>
        )}

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
                <div className="text-xs text-muted-foreground flex items-center gap-2 flex-wrap">
                  <span>{city.province}</span>
                  {coverageFor(city.id)?.missing_fields.length ? (
                    <span className="text-red-500">{coverageFor(city.id)?.missing_fields.length} 项待补</span>
                  ) : (
                    <span className="text-primary">完整</span>
                  )}
                </div>
              </div>
              <span className="text-muted-foreground text-sm">{expandedCity === city.id ? '▲' : '▼'}</span>
            </button>

            {expandedCity === city.id && (
              <div className="px-4 pb-5 pt-1 border-t border-border/40 space-y-5 bg-secondary/20">
                {coverageFor(city.id)?.missing_fields.length ? (
                  <div className="text-xs text-red-500 bg-red-500/5 border border-red-500/20 rounded-xl p-3">
                    待补：{coverageFor(city.id)?.missing_fields.join('、')}
                  </div>
                ) : null}

                {cityDrafts[city.id] && (
                  <div>
                    <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">城市内容</div>
                    <div className="space-y-2">
                      <textarea
                        value={cityDrafts[city.id].intro}
                        onChange={e => setCityDrafts(prev => ({ ...prev, [city.id]: { ...prev[city.id], intro: e.target.value } }))}
                        className="w-full min-h-20 text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                        placeholder="城市简介"
                      />
                      <input
                        value={cityDrafts[city.id].tags}
                        onChange={e => setCityDrafts(prev => ({ ...prev, [city.id]: { ...prev[city.id], tags: e.target.value } }))}
                        className="w-full text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                        placeholder="标签，逗号分隔"
                      />
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                        <input
                          value={cityDrafts[city.id].dialect_sample}
                          onChange={e => setCityDrafts(prev => ({ ...prev, [city.id]: { ...prev[city.id], dialect_sample: e.target.value } }))}
                          className="w-full text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                          placeholder="方言样例"
                        />
                        <input
                          value={cityDrafts[city.id].dialect_explanation}
                          onChange={e => setCityDrafts(prev => ({ ...prev, [city.id]: { ...prev[city.id], dialect_explanation: e.target.value } }))}
                          className="w-full text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                          placeholder="方言解释"
                        />
                      </div>
                      <button
                        onClick={() => saveCityDraft(city.id)}
                        className="text-xs px-3 py-2 rounded-lg bg-primary text-primary-foreground font-bold"
                      >
                        保存城市内容
                      </button>
                    </div>
                  </div>
                )}

                <div>
                  <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">封面图</div>
                  <UploadBtn
                    label="上传封面"
                    currentUrl={city.cover_image_url}
                    onUpload={async (file) => {
                      const res = await api.adminUploadCityCover(city.id, file, token);
                      pushToast(`${city.name} 封面已更新`, 'success');
                      await Promise.all([loadCities(), refreshCityDetail(city.id), refreshCoverage()]);
                      return res;
                    }}
                    onBindUrl={async (url) => {
                      const res = await api.adminBindCityCoverURL(city.id, url, token);
                      pushToast(`${city.name} 封面已导入本地`, 'success');
                      await Promise.all([loadCities(), refreshCityDetail(city.id), refreshCoverage()]);
                      return res;
                    }}
                  />
                </div>

                {cityDetails[city.id] && (
                  <div>
                    <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">地标</div>
                    <div className="grid grid-cols-1 md:grid-cols-[1fr_1fr_auto] gap-2 mb-3">
                      <input
                        value={newDraftFor(city.id).landmarkName}
                        onChange={e => updateNewDraft(city.id, 'landmarkName', e.target.value)}
                        className="text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                        placeholder="地标名称"
                      />
                      <input
                        value={newDraftFor(city.id).landmarkDescription}
                        onChange={e => updateNewDraft(city.id, 'landmarkDescription', e.target.value)}
                        className="text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                        placeholder="描述"
                      />
                      <button
                        onClick={() => createLandmark(city.id)}
                        disabled={!newDraftFor(city.id).landmarkName.trim()}
                        className="text-xs px-3 py-2 rounded-lg bg-primary text-primary-foreground font-bold disabled:opacity-50"
                      >
                        新增
                      </button>
                    </div>
                    <div className="space-y-2">
                      {cityDetails[city.id].landmarks.map(lm => (
                        <div key={lm.id} className="flex items-center gap-3 flex-wrap">
                          <span className="text-sm flex-1">{lm.name}</span>
                          <button
                            onClick={() => editPOI(city.id, 'landmark', lm)}
                            className="text-xs px-2 py-1 rounded-lg border border-border hover:bg-secondary transition-colors"
                          >
                            编辑
                          </button>
                          <button
                            onClick={() => deleteItem(city.id, 'landmark', lm.id, lm.name)}
                            className="text-xs px-2 py-1 rounded-lg border border-red-500/30 text-red-500 hover:bg-red-500/5 transition-colors"
                          >
                            删除
                          </button>
                          <UploadBtn
                            label="上传"
                            currentUrl={lm.image_url}
                            onUpload={async (file) => {
                              const res = await api.adminUploadLandmarkImage(lm.id, file, token);
                              pushToast(`${lm.name} 图片已更新`, 'success');
                              await Promise.all([refreshCityDetail(city.id), refreshCoverage()]);
                              return res;
                            }}
                            onBindUrl={async (url) => {
                              const res = await api.adminBindLandmarkImageURL(lm.id, url, token);
                              pushToast(`${lm.name} 图片已导入本地`, 'success');
                              await Promise.all([refreshCityDetail(city.id), refreshCoverage()]);
                              return res;
                            }}
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {cityDetails[city.id] && (
                  <div>
                    <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">人物</div>
                    <div className="grid grid-cols-1 md:grid-cols-[1fr_auto] gap-2 mb-3">
                      <input
                        value={newDraftFor(city.id).characterName}
                        onChange={e => updateNewDraft(city.id, 'characterName', e.target.value)}
                        className="text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                        placeholder="人物名称"
                      />
                      <button
                        onClick={() => createCharacter(city.id)}
                        disabled={!newDraftFor(city.id).characterName.trim()}
                        className="text-xs px-3 py-2 rounded-lg bg-primary text-primary-foreground font-bold disabled:opacity-50"
                      >
                        新增
                      </button>
                    </div>
                    <div className="space-y-2">
                      {cityDetails[city.id].characters.map(ch => (
                        <div key={ch.id} className="flex items-center gap-3 flex-wrap">
                          <span className="text-sm flex-1">{ch.name}</span>
                          <button
                            onClick={() => editCharacter(city.id, ch)}
                            className="text-xs px-2 py-1 rounded-lg border border-border hover:bg-secondary transition-colors"
                          >
                            编辑
                          </button>
                          <button
                            onClick={() => deleteItem(city.id, 'character', ch.id, ch.name)}
                            className="text-xs px-2 py-1 rounded-lg border border-red-500/30 text-red-500 hover:bg-red-500/5 transition-colors"
                          >
                            删除
                          </button>
                          <UploadBtn
                            label="上传"
                            currentUrl={ch.avatar_url}
                            onUpload={async (file) => {
                              const res = await api.adminUploadCharacterAvatar(ch.id, file, token);
                              pushToast(`${ch.name} 头像已更新`, 'success');
                              await Promise.all([refreshCityDetail(city.id), refreshCoverage()]);
                              return res;
                            }}
                            onBindUrl={async (url) => {
                              const res = await api.adminBindCharacterAvatarURL(ch.id, url, token);
                              pushToast(`${ch.name} 头像已导入本地`, 'success');
                              await Promise.all([refreshCityDetail(city.id), refreshCoverage()]);
                              return res;
                            }}
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {cityDetails[city.id] && (
                  <div>
                    <div className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-2">美食</div>
                    <div className="grid grid-cols-1 md:grid-cols-[1fr_1fr_auto] gap-2 mb-3">
                      <input
                        value={newDraftFor(city.id).foodName}
                        onChange={e => updateNewDraft(city.id, 'foodName', e.target.value)}
                        className="text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                        placeholder="美食名称"
                      />
                      <input
                        value={newDraftFor(city.id).foodDescription}
                        onChange={e => updateNewDraft(city.id, 'foodDescription', e.target.value)}
                        className="text-sm px-3 py-2 rounded-xl border border-border bg-background outline-none focus:border-primary/50"
                        placeholder="描述"
                      />
                      <button
                        onClick={() => createFood(city.id)}
                        disabled={!newDraftFor(city.id).foodName.trim()}
                        className="text-xs px-3 py-2 rounded-lg bg-primary text-primary-foreground font-bold disabled:opacity-50"
                      >
                        新增
                      </button>
                    </div>
                    <div className="space-y-2">
                      {cityDetails[city.id].foods.map(f => (
                        <div key={f.id} className="flex items-center gap-3 flex-wrap">
                          <span className="text-sm flex-1">{f.name}</span>
                          <button
                            onClick={() => editPOI(city.id, 'food', f)}
                            className="text-xs px-2 py-1 rounded-lg border border-border hover:bg-secondary transition-colors"
                          >
                            编辑
                          </button>
                          <button
                            onClick={() => deleteItem(city.id, 'food', f.id, f.name)}
                            className="text-xs px-2 py-1 rounded-lg border border-red-500/30 text-red-500 hover:bg-red-500/5 transition-colors"
                          >
                            删除
                          </button>
                          <UploadBtn
                            label="上传"
                            currentUrl={f.image_url}
                            onUpload={async (file) => {
                              const res = await api.adminUploadFoodImage(f.id, file, token);
                              pushToast(`${f.name} 图片已更新`, 'success');
                              await Promise.all([refreshCityDetail(city.id), refreshCoverage()]);
                              return res;
                            }}
                            onBindUrl={async (url) => {
                              const res = await api.adminBindFoodImageURL(f.id, url, token);
                              pushToast(`${f.name} 图片已导入本地`, 'success');
                              await Promise.all([refreshCityDetail(city.id), refreshCoverage()]);
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
