import { apiClient } from './client';
import type {
  UserResponse, CityListResponse, CityDetail,
  FreeVisitResponse, GameInitResponse, GameRollResponse,
  ChatResponse, GenerateImageResponse, CheckinResponse,
  AchievementWallResponse,
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
};
