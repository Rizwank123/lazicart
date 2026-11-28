import { router, usePathname } from "expo-router";
import { useEffect } from "react";

import { useAuthStore } from "@/store/auth-store";

/**
 * Require a signed-in session for a protected screen.
 *
 * The cart, orders, notifications, profile, checkout and order detail screens
 * all call JWT protected endpoints. A signed-out visitor would only see error
 * states, so the screen sends them to the sign-in form instead. The current
 * path is kept so that the user returns to it after signing in.
 */
export function useRequireAuth() {
  const hydrated = useAuthStore((state) => state.hydrated);
  const tokens = useAuthStore((state) => state.tokens);
  const pathname = usePathname();

  useEffect(() => {
    if (!hydrated || tokens) return;
    router.replace({ pathname: "/(auth)/login", params: { from: pathname } });
  }, [hydrated, tokens, pathname]);

  return { ready: hydrated, signedIn: Boolean(tokens) };
}
