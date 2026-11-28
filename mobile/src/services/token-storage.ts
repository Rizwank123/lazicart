import * as SecureStore from "expo-secure-store";
import { Platform } from "react-native";

import type { AuthTokens } from "@/types/api";

const ACCESS_TOKEN = "lahzicart.access_token";
const REFRESH_TOKEN = "lahzicart.refresh_token";
const JTI = "lahzicart.jti";

/**
 * SecureStore is not available on web, so the tokens are kept in
 * localStorage there. The in-memory fallback is used when localStorage is
 * missing as well, for example during server side rendering.
 */
const memory = new Map<string, string>();

const storage = Platform.OS === "web"
  ? {
      getItem: async (key: string) => {
        if (typeof localStorage !== "undefined") return localStorage.getItem(key);
        return memory.get(key) ?? null;
      },
      setItem: async (key: string, value: string) => {
        if (typeof localStorage !== "undefined") {
          localStorage.setItem(key, value);
          return;
        }
        memory.set(key, value);
      },
      removeItem: async (key: string) => {
        if (typeof localStorage !== "undefined") {
          localStorage.removeItem(key);
          return;
        }
        memory.delete(key);
      },
    }
  : {
      getItem: (key: string) => SecureStore.getItemAsync(key),
      setItem: (key: string, value: string) => SecureStore.setItemAsync(key, value),
      removeItem: (key: string) => SecureStore.deleteItemAsync(key),
    };

export async function saveTokens(tokens: AuthTokens) {
  await Promise.all([
    storage.setItem(ACCESS_TOKEN, tokens.access_token),
    storage.setItem(REFRESH_TOKEN, tokens.refresh_token),
    storage.setItem(JTI, tokens.jti),
  ]);
}

export async function getTokens(): Promise<AuthTokens | null> {
  const [access, refresh, jti] = await Promise.all([
    storage.getItem(ACCESS_TOKEN),
    storage.getItem(REFRESH_TOKEN),
    storage.getItem(JTI),
  ]);
  if (!access || !refresh || !jti) return null;
  return { access_token: access, refresh_token: refresh, jti };
}

export async function clearTokens() {
  await Promise.all([
    storage.removeItem(ACCESS_TOKEN),
    storage.removeItem(REFRESH_TOKEN),
    storage.removeItem(JTI),
  ]);
}
