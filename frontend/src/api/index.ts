import { apiClient } from './client';
import type {
  UserResponse, CityListResponse, CityDetail,
  Landmark, Food, Character,
  FreeVisitResponse, GameInitResponse, GameRollResponse,
  ChatResponse, GenerateImageResponse, CheckinResponse,
  AchievementWallResponse, AdminUploadResponse, ImageTaskResponse,
  AdminCoverageResponse, AdminUpdateResponse, UserAssetsResponse,
  CommentListResponse, CommentItem, CommentTargetType,
  GuessCaptionResponse, GuessChallengeResponse, GuessAnswerResponse,
  UserProfileResponse, AuthResponse, AdminTagListResponse,
  AdminAchievementListResponse, AdminAchievement,
} from './types';

export const api = {
  createAnonymousUser: (anonymousId: string) =>
    apiClient.post<unknown, UserResponse>('/users/anonymous', { anonymous_id: anonymousId }),

  register: (payload: { user_id?: number | null; username: string; password: string; nickname?: string }) =>
    apiClient.post<unknown, AuthResponse>('/auth/register', payload),

  login: (payload: { username: string; password: string }) =>
    apiClient.post<unknown, AuthResponse>('/auth/login', payload),

  getCities: () =>
    apiClient.get<unknown, CityListResponse>('/cities'),

  getCityDetail: (cityId: number) =>
    apiClient.get<unknown, CityDetail>(`/cities/${cityId}`),

  createFreeVisit: (userId: number, cityId: number) =>
    apiClient.post<unknown, FreeVisitResponse>('/visits/free', { user_id: userId, city_id: cityId, source: 'map_click' }),

  gameInit: (userId: number, lat?: number, lng?: number) =>
    apiClient.post<unknown, GameInitResponse>('/game/init', { user_id: userId, lat, lng }),

  gameRoll: (userId: number, fromCityId: number, lat: number, lng: number) =>
    apiClient.post<unknown, GameRollResponse>('/game/roll', { user_id: userId, from_city_id: fromCityId, lat, lng }),

  chat: (userId: number, cityId: number, characterId: number, message: string) =>
    apiClient.post<unknown, ChatResponse>('/chat', { user_id: userId, city_id: cityId, character_id: characterId, message }),

  getComments: (targetType: CommentTargetType, targetId: number, limit = 50) =>
    apiClient.get<unknown, CommentListResponse>('/comments', {
      params: { target_type: targetType, target_id: targetId, limit },
    }),

  createComment: (payload: {
    target_type: CommentTargetType;
    target_id: number;
    user_id?: number | null;
    nickname?: string;
    content: string;
  }) => apiClient.post<unknown, CommentItem>('/comments', payload),

  generateGuessCaption: (payload: {
    user_id?: number | null;
    city_id: number;
    target_name?: string;
    scene_hint?: string;
  }) => apiClient.post<unknown, GuessCaptionResponse>('/guess/caption', payload),

  generateImage: (formData: FormData) =>
    apiClient.post<unknown, GenerateImageResponse>('/checkin/generate-image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  getImageTask: (taskId: number, userId: number) =>
    apiClient.get<unknown, ImageTaskResponse>(`/checkin/image-tasks/${taskId}`, {
      params: { user_id: userId },
    }),

  retryImageTask: (taskId: number, userId: number) =>
    apiClient.post<unknown, GenerateImageResponse>(`/checkin/image-tasks/${taskId}/retry`, { user_id: userId }),

  createCheckin: (userId: number, cityId: number, landmarkId?: number, visitId?: number, generatedImageUrl?: string) =>
    apiClient.post<unknown, CheckinResponse>('/checkin', {
      user_id: userId, city_id: cityId,
      landmark_id: landmarkId, visit_id: visitId,
      generated_image_url: generatedImageUrl,
    }),

  getAchievements: (userId: number) =>
    apiClient.get<unknown, AchievementWallResponse>(`/users/${userId}/achievements`),

  getUserAssets: (userId: number) =>
    apiClient.get<unknown, UserAssetsResponse>(`/users/${userId}/assets`),

  getUserProfile: (userId: number) =>
    apiClient.get<unknown, UserProfileResponse>(`/users/${userId}/profile`),

  updateUserProfile: (userId: number, payload: {
    nickname?: string;
    age?: number;
    home_region?: string;
  }) => apiClient.patch<unknown, UserProfileResponse>(`/users/${userId}/profile`, payload),

  createGuessChallenge: (payload: {
    user_id?: number | null;
    city_id: number;
    target_name?: string;
    image_url?: string;
    image_data_url?: string;
    caption?: string;
  }) => apiClient.post<unknown, GuessChallengeResponse>('/guess/challenges', payload),

  getGuessChallenge: (code: string) =>
    apiClient.get<unknown, GuessChallengeResponse>(`/guess/challenges/${encodeURIComponent(code)}`),

  answerGuessChallenge: (code: string, answerText: string) =>
    apiClient.post<unknown, GuessAnswerResponse>(`/guess/challenges/${encodeURIComponent(code)}/answers`, {
      answer_text: answerText,
    }),

  adminCoverage: (token: string) =>
    apiClient.get<unknown, AdminCoverageResponse>('/admin/catalog/coverage', {
      headers: { 'X-Admin-Token': token },
    }),

  adminTags: (token: string) =>
    apiClient.get<unknown, AdminTagListResponse>('/admin/tags', {
      headers: { 'X-Admin-Token': token },
    }),

  adminRenameTag: (tag: string, nextTag: string, token: string) =>
    apiClient.patch<unknown, AdminUpdateResponse>(`/admin/tags/${encodeURIComponent(tag)}`, { tag: nextTag }, {
      headers: { 'X-Admin-Token': token },
    }),

  adminDeleteTag: (tag: string, token: string) =>
    apiClient.delete<unknown, AdminUpdateResponse>(`/admin/tags/${encodeURIComponent(tag)}`, {
      headers: { 'X-Admin-Token': token },
    }),

  adminAchievements: (token: string) =>
    apiClient.get<unknown, AdminAchievementListResponse>('/admin/achievements', {
      headers: { 'X-Admin-Token': token },
    }),

  adminCreateAchievement: (payload: Record<string, unknown>, token: string) =>
    apiClient.post<unknown, AdminAchievement>('/admin/achievements', payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminUpdateAchievement: (achievementId: number, payload: Record<string, unknown>, token: string) =>
    apiClient.patch<unknown, AdminUpdateResponse>(`/admin/achievements/${achievementId}`, payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminDeleteAchievement: (achievementId: number, token: string) =>
    apiClient.delete<unknown, AdminUpdateResponse>(`/admin/achievements/${achievementId}`, {
      headers: { 'X-Admin-Token': token },
    }),

  adminUpdateCity: (cityId: number, payload: Record<string, unknown>, token: string) =>
    apiClient.patch<unknown, AdminUpdateResponse>(`/admin/cities/${cityId}`, payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminUpdateLandmark: (landmarkId: number, payload: Record<string, unknown>, token: string) =>
    apiClient.patch<unknown, AdminUpdateResponse>(`/admin/landmarks/${landmarkId}`, payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminUpdateFood: (foodId: number, payload: Record<string, unknown>, token: string) =>
    apiClient.patch<unknown, AdminUpdateResponse>(`/admin/foods/${foodId}`, payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminUpdateCharacter: (characterId: number, payload: Record<string, unknown>, token: string) =>
    apiClient.patch<unknown, AdminUpdateResponse>(`/admin/characters/${characterId}`, payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminCreateLandmark: (cityId: number, payload: Record<string, unknown>, token: string) =>
    apiClient.post<unknown, Landmark>(`/admin/cities/${cityId}/landmarks`, payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminCreateFood: (cityId: number, payload: Record<string, unknown>, token: string) =>
    apiClient.post<unknown, Food>(`/admin/cities/${cityId}/foods`, payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminCreateCharacter: (cityId: number, payload: Record<string, unknown>, token: string) =>
    apiClient.post<unknown, Character>(`/admin/cities/${cityId}/characters`, payload, {
      headers: { 'X-Admin-Token': token },
    }),

  adminDeleteLandmark: (landmarkId: number, token: string) =>
    apiClient.delete<unknown, AdminUpdateResponse>(`/admin/landmarks/${landmarkId}`, {
      headers: { 'X-Admin-Token': token },
    }),

  adminDeleteFood: (foodId: number, token: string) =>
    apiClient.delete<unknown, AdminUpdateResponse>(`/admin/foods/${foodId}`, {
      headers: { 'X-Admin-Token': token },
    }),

  adminDeleteCharacter: (characterId: number, token: string) =>
    apiClient.delete<unknown, AdminUpdateResponse>(`/admin/characters/${characterId}`, {
      headers: { 'X-Admin-Token': token },
    }),

  adminUploadCityCover: (cityId: number, file: File, token: string) => {
    const fd = new FormData(); fd.append('file', file);
    return apiClient.post<unknown, AdminUploadResponse>(`/admin/cities/${cityId}/cover-image`, fd, {
      headers: { 'Content-Type': 'multipart/form-data', 'X-Admin-Token': token },
    });
  },

  adminUploadLandmarkImage: (landmarkId: number, file: File, token: string) => {
    const fd = new FormData(); fd.append('file', file);
    return apiClient.post<unknown, AdminUploadResponse>(`/admin/landmarks/${landmarkId}/image`, fd, {
      headers: { 'Content-Type': 'multipart/form-data', 'X-Admin-Token': token },
    });
  },

  adminUploadCharacterAvatar: (characterId: number, file: File, token: string) => {
    const fd = new FormData(); fd.append('file', file);
    return apiClient.post<unknown, AdminUploadResponse>(`/admin/characters/${characterId}/avatar`, fd, {
      headers: { 'Content-Type': 'multipart/form-data', 'X-Admin-Token': token },
    });
  },

  adminUploadFoodImage: (foodId: number, file: File, token: string) => {
    const fd = new FormData(); fd.append('file', file);
    return apiClient.post<unknown, AdminUploadResponse>(`/admin/foods/${foodId}/image`, fd, {
      headers: { 'Content-Type': 'multipart/form-data', 'X-Admin-Token': token },
    });
  },

  adminUploadAchievementBadge: (achievementId: number, file: File, token: string) => {
    const fd = new FormData(); fd.append('file', file);
    return apiClient.post<unknown, AdminUploadResponse>(`/admin/achievements/${achievementId}/badge`, fd, {
      headers: { 'Content-Type': 'multipart/form-data', 'X-Admin-Token': token },
    });
  },

  adminBindCityCoverURL: (cityId: number, url: string, token: string) =>
    apiClient.patch<unknown, AdminUploadResponse>(`/admin/cities/${cityId}/cover-image`, { url }, {
      headers: { 'X-Admin-Token': token },
    }),

  adminBindLandmarkImageURL: (landmarkId: number, url: string, token: string) =>
    apiClient.patch<unknown, AdminUploadResponse>(`/admin/landmarks/${landmarkId}/image`, { url }, {
      headers: { 'X-Admin-Token': token },
    }),

  adminBindCharacterAvatarURL: (characterId: number, url: string, token: string) =>
    apiClient.patch<unknown, AdminUploadResponse>(`/admin/characters/${characterId}/avatar`, { url }, {
      headers: { 'X-Admin-Token': token },
    }),

  adminBindFoodImageURL: (foodId: number, url: string, token: string) =>
    apiClient.patch<unknown, AdminUploadResponse>(`/admin/foods/${foodId}/image`, { url }, {
      headers: { 'X-Admin-Token': token },
    }),

  adminBindAchievementBadgeURL: (achievementId: number, url: string, token: string) =>
    apiClient.patch<unknown, AdminUploadResponse>(`/admin/achievements/${achievementId}/badge`, { url }, {
      headers: { 'X-Admin-Token': token },
    }),

  fetchLocalImageFile: async (url: string, filename: string): Promise<File | null> => {
    if (!url.startsWith('/static/') && !url.startsWith('/uploads/')) return null;
    const res = await fetch(url);
    if (!res.ok) return null;
    const blob = await res.blob();
    return new File([blob], filename, { type: blob.type || 'image/png', lastModified: Date.now() });
  },
};
