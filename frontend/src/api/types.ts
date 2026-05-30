// API 类型定义，对齐 api-contract.md

export interface ApiError {
  code: string;
  message: string;
}

export interface UserResponse {
  user_id: number;
  anonymous_id: string;
  current_city_id: number | null;
}

export interface CityListItem {
  id: number;
  name: string;
  province: string;
  lat: number;
  lng: number;
  cover_image_url?: string;
  tags: string[];
}

export interface CityListResponse {
  cities: CityListItem[];
}

export interface Landmark {
  id: number;
  name: string;
  image_url?: string;
  description?: string;
}

export interface Food {
  id: number;
  name: string;
  image_url?: string;
  description?: string;
}

export interface Character {
  id: number;
  name: string;
  character_type: string;
  avatar_url?: string;
  dialect_style?: string;
}

export interface CityDetail {
  id: number;
  name: string;
  province: string;
  lat: number;
  lng: number;
  intro?: string;
  cover_image_url?: string;
  dialect_sample?: string;
  dialect_explanation?: string;
  tags: string[];
  landmarks: Landmark[];
  foods: Food[];
  characters: Character[];
}

export interface FreeVisitResponse {
  visit_id: number;
  city_id: number;
  visit_mode: string;
}

export interface NearestCity {
  id: number;
  name: string;
  province: string;
  lat: number;
  lng: number;
}

export interface GameInitResponse {
  nearest_city: NearestCity;
}

export interface TargetPoint {
  lat: number;
  lng: number;
}

export interface GameRollResponse {
  visit_id: number;
  dice_roll_id: number;
  direction: string;
  distance_km: number;
  target_point: TargetPoint;
  target_city: NearestCity;
}

export interface ChatResponse {
  reply: string;
}

export interface GenerateImageResponse {
  status: string;
  generated_image_url: string;
}

export interface AchievementBrief {
  code: string;
  name: string;
  description?: string;
}

export interface CheckinResponse {
  checkin_id: number;
  unlocked_achievements: AchievementBrief[];
}

export interface UnlockedAchievement {
  code: string;
  name: string;
  description?: string;
  badge_url?: string;
  unlocked_at: string;
}

export interface LockedAchievement {
  code: string;
  name: string;
  description?: string;
  badge_url?: string;
}

export interface ProgressItem {
  code: string;
  current: number;
  target: number;
}

export interface AchievementWallResponse {
  unlocked: UnlockedAchievement[];
  locked: LockedAchievement[];
  progress: ProgressItem[];
}
