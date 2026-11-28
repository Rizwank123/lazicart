import axios, { AxiosError, create, InternalAxiosRequestConfig, isAxiosError } from "axios";

import { env } from "@/config/env";
import { clearTokens, getTokens, saveTokens } from "@/services/token-storage";
import type { ApiResponse, AuthTokens } from "@/types/api";

type RetriableConfig = InternalAxiosRequestConfig & { _retry?: boolean };

type TokenListener = (tokens: AuthTokens | null) => void;

const tokenListeners = new Set<TokenListener>();

/**
 * Subscribe to token changes performed by the refresh interceptor.
 *
 * The interceptor owns token rotation, so the auth store has to learn about
 * new tokens. Without this the store keeps a stale jti, and logout would
 * revoke the already-rotated refresh token instead of the live one.
 */
export function subscribeToTokenChanges(listener: TokenListener) {
  tokenListeners.add(listener);
  return () => {
    tokenListeners.delete(listener);
  };
}

function notifyTokenChange(tokens: AuthTokens | null) {
  for (const listener of tokenListeners) {
    try {
      listener(tokens);
    } catch {
      // A misbehaving listener must not break the request pipeline.
    }
  }
}

export const api = create({
  baseURL: env.apiBaseUrl,
  timeout: 20_000,
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use(async (config) => {
  const tokens = await readTokens();
  if (tokens?.access_token) {
    config.headers.Authorization = `Bearer ${tokens.access_token}`;
  }
  return config;
});

/**
 * The in-flight refresh operation, shared by every caller.
 *
 * A refresh token is single use: the backend rotates it on every refresh.
 * Without this guard two requests that fail with 401 at the same time would
 * both try to refresh, one of them would fail because its refresh token had
 * already been consumed, and that failure would clear the session that the
 * successful refresh had just created.
 */
let refreshPromise: Promise<AuthTokens | null> | null = null;

async function refreshAccessToken(): Promise<AuthTokens | null> {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      try {
        const tokens = await getTokens();
        if (!tokens?.refresh_token) return null;
        const refreshed = await axios.post<ApiResponse<AuthTokens>>(
          `${env.apiBaseUrl}/api/auth/refresh-token`,
          undefined,
          { params: { refresh_token: tokens.refresh_token }, timeout: 20_000 },
        );
        const rotated = refreshed.data.data;
        if (!rotated) return null;
        await saveTokens(rotated);
        notifyTokenChange(rotated);
        return rotated;
      } catch {
        // The session cannot be recovered.
        return null;
      } finally {
        refreshPromise = null;
      }
    })();
  }
  return refreshPromise;
}

/**
 * Read the stored tokens without throwing.
 *
 * A storage failure must not replace the original request error with a
 * storage error, so the failure is reported as "no session".
 */
async function readTokens(): Promise<AuthTokens | null> {
  try {
    return await getTokens();
  } catch {
    return null;
  }
}

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiResponse<unknown>>) => {
    const original = error.config as RetriableConfig | undefined;
    const isAuthRoute = original?.url?.startsWith("/api/auth/");
    if (error.response?.status === 401 && original && !original._retry && !isAuthRoute) {
      original._retry = true;

      // The token that the failed request used. If it is no longer the stored
      // token then another request already refreshed the session, so the
      // request can simply be replayed with the new token.
      const usedAuthorization = original.headers?.Authorization;
      const current = await readTokens();
      const currentAuthorization = current?.access_token ? `Bearer ${current.access_token}` : undefined;

      if (usedAuthorization && currentAuthorization && usedAuthorization !== currentAuthorization) {
        original.headers = original.headers ?? {};
        original.headers.Authorization = currentAuthorization;
        return api(original);
      }

      const refreshed = await refreshAccessToken();
      if (refreshed) {
        original.headers = original.headers ?? {};
        original.headers.Authorization = `Bearer ${refreshed.access_token}`;
        return api(original);
      }

      // The session is unusable. Clear it everywhere so the app returns to
      // the sign-in screen instead of retrying with a dead token.
      try {
        await clearTokens();
      } catch {
        // The store is cleared anyway, so the app returns to sign-in.
      }
      notifyTokenChange(null);
    }
    return Promise.reject(error);
  },
);

export function getApiError(error: unknown, fallback = "Something went wrong") {
  if (isAxiosError<ApiResponse<unknown>>(error)) {
    return error.response?.data?.message || error.message || fallback;
  }
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return fallback;
}

export async function unwrap<T>(promise: Promise<{ data: ApiResponse<T> }>): Promise<T> {
  const response = await promise;
  const data = response.data.data;
  if (data === undefined || data === null) {
    throw new Error(response.data.message || "The server returned an empty response");
  }
  return data;
}

/**
 * Resolve a response that carries no payload.
 *
 * The backend omits the `data` key entirely for these operations, so
 * `unwrap` would resolve to `undefined` rather than the declared `null`.
 */
export async function unwrapVoid(promise: Promise<{ data: ApiResponse<unknown> }>): Promise<void> {
  await promise;
}
