import type { CityListItem } from '../api/types';

// 区域分组：每个区域包含一组 tag（来自 cities.json 实际 tags）
export const REGION_GROUPS: Record<string, { label: string; tags: string[] }> = {
  north_china: { label: '华北 · 古都', tags: ['north_china', 'ancient_capital'] },
  jiangnan:    { label: '江南 · 水乡', tags: ['jiangnan', 'canal_city', 'water_town'] },
  southwest:   { label: '西南 · 人文', tags: ['southwest', 'ethnic_culture', 'plateau'] },
  coastal:     { label: '沿海 · 港口', tags: ['coastal', 'port', 'maritime'] },
  northwest:   { label: '西北 · 丝路', tags: ['northwest', 'silk_road', 'desert'] },
  food_city:   { label: '美食之都', tags: ['food_city', 'spicy_food', 'dim_sum'] },
  historic:    { label: '历史名城', tags: ['historic', 'ancient_wall', 'heritage'] },
};

export interface RegionOption {
  key: string;
  label: string;
}

/** 从城市列表中提取实际存在对应 tag 的区域选项 */
export function getRegionOptions(cities: CityListItem[]): RegionOption[] {
  const allTags = new Set(cities.flatMap(c => c.tags ?? []));
  return Object.entries(REGION_GROUPS)
    .filter(([, grp]) => grp.tags.some(t => allTags.has(t)))
    .map(([key, grp]) => ({ key, label: grp.label }));
}

/** 从城市列表中提取所有独立 tag */
export function getAllTags(cities: CityListItem[]): string[] {
  const set = new Set(cities.flatMap(c => c.tags ?? []));
  return Array.from(set).sort();
}

/** 多维度筛选城市 */
export function filterCities(
  cities: CityListItem[],
  query: string,
  activeRegion: string | null,
  activeTag: string | null,
): CityListItem[] {
  let result = cities;

  if (query) {
    const q = query.toLowerCase();
    result = result.filter(c =>
      c.name.includes(query) ||
      c.province.includes(query) ||
      c.tags?.some(t => t.includes(q))
    );
  }

  if (activeRegion && REGION_GROUPS[activeRegion]) {
    const regionTags = new Set(REGION_GROUPS[activeRegion].tags);
    result = result.filter(c => c.tags?.some(t => regionTags.has(t)));
  }

  if (activeTag) {
    result = result.filter(c => c.tags?.includes(activeTag));
  }

  return result;
}
