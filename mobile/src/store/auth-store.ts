import { create } from "zustand";

import { getApiError, subscribeToTokenChanges } from "@/api/client";
import { authApi, userApi } from "@/api/endpoints";
import { clearTokens, getTokens, saveTokens } from "@/services/token-storage";
import { useCheckoutStore } from "@/store/checkout-store";
import { useWishlistStore } from "@/store/wishlist-store";
import type { AuthTokens, User } from "@/types/api";
import { decodeUserId } from "@/utils/jwt";

// The user profile is created asynchronously by the user-service from the
// user.created Kafka event, so it can be missing for a short while right after
// signup or login. Retry a few times before giving up.
const PROFILE_RESOLVE_ATTEMPTS = 5;
const PROFILE_RESOLVE_DELAY_MS = 600;
const PROFILE_RESOLVE_TIMEOUT_MS = 8_000;
// The whole resolution is bounded so that a failing user-service cannot keep
// the app waiting for minutes.
const PROFILE_RESOLVE_DEADLINE_MS = 20_000;

const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

type AuthState = {
  tokens: AuthTokens | null;
  user: User | null;
  /** auth-service identifier, carried in the JWT subject claim. */
  authId: string | null;
  /**
   * user-service identifier. The cart, order, payment, notification and
   * delivery APIs are all keyed by this id, not by the auth-service id.
   */
  userId: string | null;
  /** True once the canonical user id has been resolved (or definitively absent). */
  userResolved: boolean;
  /** Why the profile could not be resolved, if that happened. */
  userError: string | null;
  hydrated: boolean;
  /**
   * Incremented whenever the signed-in identity changes.
   *
   * A profile request that is still in flight belongs to the session that
   * started it. The generation lets the store drop such a stale response
   * instead of restoring the previous account.
   */
  sessionId: number;
  setSession: (tokens: AuthTokens) => Promise<void>;
  hydrate: () => Promise<void>;
  loadUser: () => Promise<void>;
  logout: () => Promise<void>;
};

export const useAuthStore = create<AuthState>((set, get) => ({
  tokens: null,
  user: null,
  authId: null,
  userId: null,
  userResolved: false,
  userError: null,
  hydrated: false,
  sessionId: 0,
  setSession: async (tokens) => {
    await saveTokens(tokens);
    set({
      tokens,
      authId: decodeUserId(tokens.access_token),
      userId: null,
      user: null,
      userResolved: false,
      userError: null,
      sessionId: get().sessionId + 1,
    });
    // The profile is resolved in the background. A slow user-service must not
    // keep the sign-in button disabled for a minute.
    void get().loadUser();
  },
  hydrate: async () => {
    try {
      const tokens = await getTokens();
      set({
        tokens,
        authId: decodeUserId(tokens?.access_token),
        userId: null,
        user: null,
        // Without a session there is no profile to resolve, so the app must
        // not wait for one that will never arrive.
        userResolved: !tokens,
        userError: null,
        hydrated: true,
        sessionId: get().sessionId + 1,
      });
      if (tokens) void get().loadUser();
    } catch {
      // A storage failure must not brick the app: drop the unusable session
      // and continue as a signed-out visitor.
      try {
        await clearTokens();
      } catch {
        // Nothing else can be done about an unreadable store.
      }
      set({
        tokens: null,
        authId: null,
        userId: null,
        user: null,
        userResolved: true,
        userError: null,
        hydrated: true,
        sessionId: get().sessionId + 1,
      });
    }
  },
  loadUser: async () => {
    const authId = get().authId;
    const sessionId = get().sessionId;
    if (!authId) {
      set({ user: null, userId: null, userResolved: true });
      return;
    }
    const deadline = Date.now() + PROFILE_RESOLVE_DEADLINE_MS;
    let message = "Your account details are unavailable right now.";
    for (let attempt = 0; attempt < PROFILE_RESOLVE_ATTEMPTS; attempt += 1) {
      const remaining = deadline - Date.now();
      if (remaining <= 0) break;
      try {
        // The user-service resolves either its primary key or the auth_id, so
        // the auth-service id from the token is accepted here.
        const user = await userApi.byId(authId, Math.min(PROFILE_RESOLVE_TIMEOUT_MS, remaining));
        // The session may have changed while the request was in flight.
        if (get().sessionId !== sessionId) return;
        set({ user, userId: user.id, userResolved: true, userError: null });
        return;
      } catch (err) {
        message = getApiError(err, message);
        if (get().sessionId !== sessionId) return;
        if (attempt < PROFILE_RESOLVE_ATTEMPTS - 1) {
          const delay = Math.min(PROFILE_RESOLVE_DELAY_MS * (attempt + 1), Math.max(0, deadline - Date.now()));
          if (delay > 0) await wait(delay);
        }
      }
    }
    if (get().sessionId !== sessionId) return;
    // The profile could not be resolved. Leave userId null so that no
    // user-scoped request is made against the wrong identifier.
    set({ user: null, userId: null, userResolved: true, userError: message });
  },
  logout: async () => {
    const jti = get().tokens?.jti;
    const sessionId = get().sessionId + 1;
    try {
      await clearTokens();
    } catch {
      // The in-memory session is cleared regardless.
    } finally {
      set({
        tokens: null,
        authId: null,
        userId: null,
        user: null,
        userResolved: true,
        userError: null,
        sessionId,
      });
      // Account specific local state must not leak into the next session.
      useCheckoutStore.getState().reset();
      useWishlistStore.getState().reset();
    }
    if (jti) {
      // Revoking the server session is best effort and must not delay the
      // sign-out.
      void authApi.logout(jti).catch(() => undefined);
    }
  },
}));

// The API client rotates tokens in the background when an access token
// expires. Mirror those changes so the store never holds a stale token or jti,
// and drop the session when the refresh fails.
subscribeToTokenChanges((tokens) => {
  if (tokens) {
    useAuthStore.setState({ tokens });
    return;
  }
  useAuthStore.setState({
    tokens: null,
    authId: null,
    userId: null,
    user: null,
    userResolved: true,
    userError: null,
    sessionId: useAuthStore.getState().sessionId + 1,
  });
});
