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

export type CommentTargetType = 'landmark' | 'food' | 'character';

export interface CommentItem {
  id: number;
  target_type: CommentTargetType;
  target_id: number;
  user_id?: number | null;
  nickname: string;
  content: string;
  created_at: string;
}

export interface CommentListResponse {
  comments: CommentItem[];
}

export interface GuessCaptionResponse {
  weibo: string;
  moments: string;
  hashtags: string[];
}

export interface GenerateImageResponse {
  status: string;
  task_id: number;
}

export type ImageTaskStatusValue = 'queued' | 'running' | 'succeeded' | 'failed' | 'retryable';

export interface ImageTaskResponse {
  task_id: number;
  status: ImageTaskStatusValue;
  result_url?: string | null;
  error?: string | null;
  attempts: number;
  created_at: string;
  updated_at: string;
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

export interface AdminUploadResponse {
  cover_image_url?: string;
  image_url?: string;
  avatar_url?: string;
}

export interface AdminCoverageItem {
  city_id: number;
  city_name: string;
  has_cover_image: boolean;
  tag_count: number;
  landmark_count: number;
  food_count: number;
  character_count: number;
  missing_fields: string[];
}

export interface AdminCoverageResponse {
  total_cities: number;
  complete_cities: number;
  items: AdminCoverageItem[];
}

export interface AdminUpdateResponse {
  status: string;
}

export interface UserAssetCity {
  id: number;
  name: string;
  province: string;
  visited_at: string;
}

export interface UserPosterAsset {
  checkin_id: number;
  city_id: number;
  city_name: string;
  landmark_name?: string;
  generated_image_url: string;
  created_at: string;
}

export interface UserAssetsResponse {
  visited_cities: UserAssetCity[];
  posters: UserPosterAsset[];
  achievement_progress: ProgressItem[];
}
