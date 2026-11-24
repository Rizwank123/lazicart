import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { DarkTheme, DefaultTheme, ThemeProvider } from "@react-navigation/native";
import { BottomSheetModalProvider } from "@gorhom/bottom-sheet";
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import type { ReactNode } from "react";
import { useEffect, useMemo } from "react";
import { Platform } from "react-native";
import { GestureHandlerRootView } from "react-native-gesture-handler";
import "react-native-reanimated";

import { LoadingState } from "@/components/ui/state-view";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useAuthStore } from "@/store/auth-store";

export default function RootLayout() {
  const { isDark } = useAppTheme();
  const hydrated = useAuthStore((state) => state.hydrated);
  const hydrate = useAuthStore((state) => state.hydrate);
  const queryClient = useMemo(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: { retry: 1, staleTime: 60_000 },
          mutations: { retry: 0 },
        },
      }),
    [],
  );

  useEffect(() => {
    hydrate();
  }, [hydrate]);

  return (
    <GestureHandlerRootView style={{ flex: 1 }}>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider value={isDark ? DarkTheme : DefaultTheme}>
          <Provider>
            {hydrated ? (
              <Stack screenOptions={{ headerShown: false }}>
                <Stack.Screen name="index" />
                <Stack.Screen name="(auth)" />
                <Stack.Screen name="(tabs)" />
                <Stack.Screen name="product/[id]" />
                <Stack.Screen name="checkout" />
                <Stack.Screen name="order/[id]" />
              </Stack>
            ) : (
              <LoadingState />
            )}
            <StatusBar style="auto" />
          </Provider>
        </ThemeProvider>
      </QueryClientProvider>
    </GestureHandlerRootView>
  );
}

function Provider({ children }: { children: ReactNode }) {
  if (Platform.OS === "web") return <>{children}</>;
  return <BottomSheetModalProvider>{children}</BottomSheetModalProvider>;
}
