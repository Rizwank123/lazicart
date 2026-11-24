import * as SecureStore from "expo-secure-store";

import type { AuthTokens } from "@/types/api";

const ACCESS_TOKEN = "lahzicart.access_token";
const REFRESH_TOKEN = "lahzicart.refresh_token";
const JTI = "lahzicart.jti";

export async function saveTokens(tokens: AuthTokens) {
  await Promise.all([
    SecureStore.setItemAsync(ACCESS_TOKEN, tokens.access_token),
    SecureStore.setItemAsync(REFRESH_TOKEN, tokens.refresh_token),
    SecureStore.setItemAsync(JTI, tokens.jti),
  ]);
}

export async function getTokens(): Promise<AuthTokens | null> {
  const [access, refresh, jti] = await Promise.all([
    SecureStore.getItemAsync(ACCESS_TOKEN),
    SecureStore.getItemAsync(REFRESH_TOKEN),
    SecureStore.getItemAsync(JTI),
  ]);
  if (!access || !refresh || !jti) return null;
  return { access_token: access, refresh_token: refresh, jti };
}

export async function clearTokens() {
  await Promise.all([
    SecureStore.deleteItemAsync(ACCESS_TOKEN),
    SecureStore.deleteItemAsync(REFRESH_TOKEN),
    SecureStore.deleteItemAsync(JTI),
  ]);
}
