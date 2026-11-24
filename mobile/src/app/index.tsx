import { Redirect } from "expo-router";

import { useAuthStore } from "@/store/auth-store";

export default function Index() {
  const tokens = useAuthStore((state) => state.tokens);
  return <Redirect href={tokens ? "/(tabs)" : "/(auth)/login"} />;
}
