import React from 'react';
import type { UserPosterAsset } from '../api/types';

interface Props {
  posters: UserPosterAsset[];
  compact?: boolean;
}

export const ProfilePosterGrid: React.FC<Props> = ({ posters, compact = false }) => {
  if (!posters.length) {
    return (
      <div className="text-sm text-muted-foreground py-4 text-center">还没有生成海报</div>
    );
  }

  const displayPosters = compact ? posters.slice(0, 4) : posters;

  return (
    <div className={`grid gap-2 ${compact ? 'grid-cols-2' : 'grid-cols-1 sm:grid-cols-2'}`}>
      {displayPosters.map(poster => (
        <article key={poster.checkin_id} className="rounded-xl overflow-hidden border border-border bg-background/70">
          <img
            src={poster.generated_image_url}
            alt={`${poster.city_name}打卡海报`}
            className="w-full aspect-video object-cover"
          />
          <div className="p-2">
            <div className="font-medium text-xs text-foreground">{poster.city_name}</div>
            <div className="text-[10px] text-muted-foreground">
              {poster.landmark_name || '城市打卡'} · {new Date(poster.created_at).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })}
            </div>
          </div>
        </article>
      ))}
      {compact && posters.length > 4 && (
        <div className="col-span-2 text-center text-xs text-muted-foreground">
          还有 {posters.length - 4} 张海报…
        </div>
      )}
    </div>
  );
};
