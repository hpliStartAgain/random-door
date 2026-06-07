import React, { useEffect, useState } from 'react';
import { Check, Pencil, X } from 'lucide-react';
import { api } from '../api';
import { useUserStore } from '../store/useUserStore';
import { useViewStore } from '../store/useViewStore';
import type { UserAssetsResponse, UserProfileResponse } from '../api/types';
import { ProfileVisitedList } from './ProfileVisitedList';
import { ProfilePosterGrid } from './ProfilePosterGrid';

export const ProfilePanel: React.FC = () => {
  const { profileOpen, setProfileOpen } = useViewStore();
  const { userId } = useUserStore();
  const [data, setData] = useState<UserAssetsResponse | null>(null);
  const [profile, setProfile] = useState<UserProfileResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editing, setEditing] = useState(false);
  const [activeTab, setActiveTab] = useState<'cities' | 'posters' | 'achievements'>('cities');
  const [nickname, setNickname] = useState('');
  const [age, setAge] = useState('');
  const [homeRegion, setHomeRegion] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (!profileOpen || !userId) return;
    setLoading(true);
    setError('');
    Promise.all([api.getUserAssets(userId), api.getUserProfile(userId)])
      .then(([assets, nextProfile]) => {
        setData(assets);
        setProfile(nextProfile);
        setNickname(nextProfile.nickname || '');
        setAge(nextProfile.age ? String(nextProfile.age) : '');
        setHomeRegion(nextProfile.home_region || '');
      })
      .catch((e) => {
        console.error(e);
        setError('资料加载失败');
      })
      .finally(() => setLoading(false));
  }, [profileOpen, userId]);

  const handleSaveProfile = async () => {
    if (!userId || saving) return;
    setSaving(true);
    setError('');
    try {
      const payload: { nickname?: string; age?: number; home_region?: string } = {};
      if (nickname.trim()) payload.nickname = nickname.trim();
      if (age.trim()) payload.age = Number(age);
      if (homeRegion.trim()) payload.home_region = homeRegion.trim();
      const next = await api.updateUserProfile(userId, payload);
      setProfile(next);
      setEditing(false);
    } catch (e: any) {
      setError(e?.message || '资料保存失败');
    } finally {
      setSaving(false);
    }
  };

  return (
    <>
      {/* 蒙层 */}
      <div
        className={`fixed inset-0 bg-black/20 backdrop-blur-sm z-40 transition-opacity duration-300 ${profileOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'}`}
        onClick={() => setProfileOpen(false)}
      />

      {/* 面板 */}
      <div
        className={`fixed top-0 right-0 w-full sm:w-[380px] max-w-[100vw] h-full bg-background/97 backdrop-blur-2xl shadow-2xl z-50 border-l border-border/50 flex flex-col transition-transform duration-500 ease-[cubic-bezier(0.16,1,0.3,1)] ${profileOpen ? 'translate-x-0' : 'translate-x-full'}`}
      >
        {/* 头部 */}
        <div className="px-6 pt-6 pb-4 border-b border-border/30 shrink-0">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary/30 to-accent/40 flex items-center justify-center border border-border/40 overflow-hidden">
                <img src="/icon-transparent.png" alt="我" className="w-8 h-8 object-contain" />
              </div>
              <div>
                <h3 className="font-bold text-base text-foreground">{profile?.nickname || '未设置昵称'}</h3>
                <p className="text-xs text-muted-foreground mt-0.5">
                  {data ? `走过 ${data.visited_cities.length} 座城市` : '加载中...'}
                </p>
              </div>
            </div>
            <button
              onClick={() => setProfileOpen(false)}
              className="w-8 h-8 rounded-full hover:bg-secondary flex items-center justify-center text-muted-foreground transition-colors shrink-0"
            >
              <X className="h-4 w-4" />
            </button>
          </div>

          <div className="mt-4 rounded-xl border border-border/60 bg-background/70 p-3 space-y-2">
            <div className="flex items-center justify-between">
              <div className="text-xs font-semibold text-muted-foreground">个人资料</div>
              <button
                onClick={() => editing ? handleSaveProfile() : setEditing(true)}
                disabled={saving || loading}
                className="h-7 px-2.5 rounded-lg border border-border text-xs font-semibold hover:bg-secondary disabled:opacity-50 flex items-center gap-1"
              >
                {editing ? <Check className="h-3.5 w-3.5" /> : <Pencil className="h-3.5 w-3.5" />}
                {editing ? (saving ? '保存中' : '保存') : '编辑'}
              </button>
            </div>
            {editing ? (
              <div className="grid grid-cols-1 gap-2">
                <input
                  value={nickname}
                  onChange={(e) => setNickname(e.target.value)}
                  placeholder="昵称"
                  className="px-3 py-2 rounded-lg border border-border bg-card text-xs outline-none focus:border-primary/50"
                />
                <div className="grid grid-cols-2 gap-2">
                  <input
                    value={age}
                    onChange={(e) => setAge(e.target.value.replace(/[^\d]/g, ''))}
                    placeholder="年龄"
                    inputMode="numeric"
                    className="px-3 py-2 rounded-lg border border-border bg-card text-xs outline-none focus:border-primary/50"
                  />
                  <input
                    value={homeRegion}
                    onChange={(e) => setHomeRegion(e.target.value)}
                    placeholder="家乡/地区"
                    className="px-3 py-2 rounded-lg border border-border bg-card text-xs outline-none focus:border-primary/50"
                  />
                </div>
              </div>
            ) : (
              <div className="grid grid-cols-2 gap-2 text-xs">
                <div className="rounded-lg bg-card px-3 py-2">
                  <span className="text-muted-foreground">年龄</span>
                  <div className="font-semibold mt-0.5">{profile?.age || '未设置'}</div>
                </div>
                <div className="rounded-lg bg-card px-3 py-2">
                  <span className="text-muted-foreground">地区</span>
                  <div className="font-semibold mt-0.5">{profile?.home_region || '未设置'}</div>
                </div>
              </div>
            )}
            {error && <div className="text-xs text-red-500">{error}</div>}
          </div>

          {/* Tabs */}
          <div className="flex gap-1 mt-4">
            {(['cities', 'posters', 'achievements'] as const).map(tab => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`flex-1 py-1.5 rounded-lg text-xs font-medium transition-colors ${activeTab === tab ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground hover:bg-secondary'}`}
              >
                {tab === 'cities' ? `城市 ${data?.visited_cities.length ?? ''}` : tab === 'posters' ? `海报 ${data?.posters.length ?? ''}` : '成就'}
              </button>
            ))}
          </div>
        </div>

        {/* 内容区 */}
        <div className="flex-1 overflow-y-auto p-5">
          {loading && <div className="text-sm text-muted-foreground text-center py-8">加载中…</div>}
          {!loading && data && (
            <>
              {activeTab === 'cities' && (
                <ProfileVisitedList cities={data.visited_cities} compact={false} />
              )}
              {activeTab === 'posters' && (
                <ProfilePosterGrid posters={data.posters} compact={false} />
              )}
              {activeTab === 'achievements' && (
                <div className="space-y-2">
                  {data.achievement_progress.length ? data.achievement_progress.map(item => (
                    <div key={item.code} className="p-3 rounded-xl border border-border bg-background/70">
                      <div className="flex justify-between text-xs mb-2">
                        <span className="font-semibold text-foreground">{item.code}</span>
                        <span className="text-muted-foreground">{item.current} / {item.target}</span>
                      </div>
                      <div className="h-1.5 bg-border rounded-full overflow-hidden">
                        <div className="h-full bg-primary rounded-full transition-all" style={{ width: `${Math.min(100, (item.current / item.target) * 100)}%` }} />
                      </div>
                    </div>
                  )) : (
                    <div className="text-sm text-muted-foreground text-center py-6">暂无进行中的成就</div>
                  )}
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </>
  );
};
