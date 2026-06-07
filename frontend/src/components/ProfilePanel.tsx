import React, { useEffect, useState } from 'react';
import { X } from 'lucide-react';
import { api } from '../api';
import { useUserStore } from '../store/useUserStore';
import { useViewStore } from '../store/useViewStore';
import type { UserAssetsResponse } from '../api/types';
import { ProfileVisitedList } from './ProfileVisitedList';
import { ProfilePosterGrid } from './ProfilePosterGrid';

export const ProfilePanel: React.FC = () => {
  const { profileOpen, setProfileOpen } = useViewStore();
  const { userId } = useUserStore();
  const [data, setData] = useState<UserAssetsResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState<'cities' | 'posters' | 'achievements'>('cities');

  useEffect(() => {
    if (!profileOpen || !userId) return;
    setLoading(true);
    api.getUserAssets(userId)
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [profileOpen, userId]);

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
                <h3 className="font-bold text-base text-foreground">我的足迹</h3>
                <p className="text-xs text-muted-foreground mt-0.5">
                  {data ? `走过 ${data.visited_cities.length} 座城市` : '加载中…'}
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
