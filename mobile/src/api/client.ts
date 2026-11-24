import axios, { AxiosError, create, InternalAxiosRequestConfig, isAxiosError } from "axios";

import { env } from "@/config/env";
import { clearTokens, getTokens, saveTokens } from "@/services/token-storage";
import type { ApiResponse, AuthTokens } from "@/types/api";

type RetriableConfig = InternalAxiosRequestConfig & { _retry?: boolean };

export const api = create({
  baseURL: env.apiBaseUrl,
  timeout: 20_000,
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use(async (config) => {
  const tokens = await getTokens();
  if (tokens?.access_token) {
    config.headers.Authorization = `Bearer ${tokens.access_token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiResponse<unknown>>) => {
    const original = error.config as RetriableConfig | undefined;
    const isAuthRoute = original?.url?.startsWith("/api/auth/");
    if (error.response?.status === 401 && original && !original._retry && !isAuthRoute) {
      original._retry = true;
      const tokens = await getTokens();
      if (tokens?.refresh_token) {
        try {
          const refreshed = await axios.post<ApiResponse<AuthTokens>>(
            `${env.apiBaseUrl}/api/auth/refresh-token`,
            undefined,
            { params: { refresh_token: tokens.refresh_token } },
          );
          if (refreshed.data.data) {
            await saveTokens(refreshed.data.data);
            original.headers.Authorization = `Bearer ${refreshed.data.data.access_token}`;
            return api(original);
          }
        } catch {
          await clearTokens();
        }
      }
    }
    return Promise.reject(error);
  },
);

export function getApiError(error: unknown, fallback = "Something went wrong") {
  if (isAxiosError<ApiResponse<unknown>>(error)) {
    return error.response?.data?.message || error.message || fallback;
  }
  return fallback;
}

export async function unwrap<T>(promise: Promise<{ data: ApiResponse<T> }>) {
  const response = await promise;
  return response.data.data as T;
}
