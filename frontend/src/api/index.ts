import { apiClient } from './client';
import type {
  UserResponse, CityListResponse, CityDetail,
  FreeVisitResponse, GameInitResponse, GameRollResponse,
  ChatResponse, GenerateImageResponse, CheckinResponse,
  AchievementWallResponse, AdminUploadResponse,
} from './types';

export const api = {
  createAnonymousUser: (anonymousId: string) =>
    apiClient.post<unknown, UserResponse>('/users/anonymous', { anonymous_id: anonymousId }),

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

  generateImage: (formData: FormData) =>
    apiClient.post<unknown, GenerateImageResponse>('/checkin/generate-image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  createCheckin: (userId: number, cityId: number, landmarkId?: number, visitId?: number, generatedImageUrl?: string) =>
    apiClient.post<unknown, CheckinResponse>('/checkin', {
      user_id: userId, city_id: cityId,
      landmark_id: landmarkId, visit_id: visitId,
      generated_image_url: generatedImageUrl,
    }),

  getAchievements: (userId: number) =>
    apiClient.get<unknown, AchievementWallResponse>(`/users/${userId}/achievements`),

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
};
