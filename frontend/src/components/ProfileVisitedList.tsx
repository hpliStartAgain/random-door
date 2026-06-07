import React from 'react';
import type { UserAssetCity } from '../api/types';

interface Props {
  cities: UserAssetCity[];
  compact?: boolean; // true=紧凑模式(ProfilePanel), false=完整模式(AssetPage)
}

export const ProfileVisitedList: React.FC<Props> = ({ cities, compact = false }) => {
  if (!cities.length) {
    return (
      <div className="text-center py-6 text-sm text-muted-foreground">
        还没有走过任何城市<br />
        <span className="text-xs">先开启一次任意门吧</span>
      </div>
    );
  }

  const displayCities = compact ? cities.slice(0, 6) : cities;

  return (
    <div className={`space-y-1.5 ${compact ? 'max-h-[200px] overflow-y-auto' : 'max-h-[520px] overflow-y-auto'}`}>
      {displayCities.map(city => (
        <div key={`${city.id}-${city.visited_at}`} className="flex items-center justify-between p-2.5 rounded-lg border border-border/60 bg-background/70 hover:border-border transition-colors">
          <div>
            <div className="font-medium text-sm text-foreground">{city.name}</div>
            <div className="text-[11px] text-muted-foreground">{city.province}</div>
          </div>
          <div className="text-[11px] text-muted-foreground shrink-0">
            {new Date(city.visited_at).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })}
          </div>
        </div>
      ))}
      {compact && cities.length > 6 && (
        <div className="text-center text-xs text-muted-foreground pt-1">
          还有 {cities.length - 6} 座城市…
        </div>
      )}
    </div>
  );
};
