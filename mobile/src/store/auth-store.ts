import { create } from "zustand";

import { authApi, userApi } from "@/api/endpoints";
import { clearTokens, getTokens, saveTokens } from "@/services/token-storage";
import type { AuthTokens, User } from "@/types/api";
import { decodeUserId } from "@/utils/jwt";

type AuthState = {
  tokens: AuthTokens | null;
  user: User | null;
  userId: string | null;
  hydrated: boolean;
  setSession: (tokens: AuthTokens) => Promise<void>;
  hydrate: () => Promise<void>;
  loadUser: () => Promise<void>;
  logout: () => Promise<void>;
};

export const useAuthStore = create<AuthState>((set, get) => ({
  tokens: null,
  user: null,
  userId: null,
  hydrated: false,
  setSession: async (tokens) => {
    await saveTokens(tokens);
    set({ tokens, userId: decodeUserId(tokens.access_token) });
    await get().loadUser();
  },
  hydrate: async () => {
    const tokens = await getTokens();
    set({ tokens, userId: decodeUserId(tokens?.access_token), hydrated: true });
    if (tokens) await get().loadUser();
  },
  loadUser: async () => {
    const userId = get().userId;
    if (!userId) return;
    try {
      const user = await userApi.byId(userId);
      set({ user });
    } catch {
      set({ user: null });
    }
  },
  logout: async () => {
    const jti = get().tokens?.jti;
    if (jti) {
      try {
        await authApi.logout(jti);
      } catch {
        // Local logout should still complete if the backend token is already gone.
      }
    }
    await clearTokens();
    set({ tokens: null, userId: null, user: null });
  },
}));
