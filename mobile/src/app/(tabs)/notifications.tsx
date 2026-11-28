import { useInfiniteQuery } from "@tanstack/react-query";
import { Pressable, RefreshControl, StyleSheet, View } from "react-native";

import { notificationApi } from "@/api/endpoints";
import { Card } from "@/components/ui/card";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state-view";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { useAuthStore } from "@/store/auth-store";

const NOTIFICATIONS_PAGE_SIZE = 20;

export default function NotificationsScreen() {
  const { colors } = useAppTheme();
  const userId = useAuthStore((state) => state.userId);
  const userResolved = useAuthStore((state) => state.userResolved);
  const { signedIn } = useRequireAuth();
  const query = useInfiniteQuery({
    queryKey: ["notifications", userId],
    initialPageParam: 1,
    queryFn: ({ pageParam }) => notificationApi.byUser(userId as string, pageParam, NOTIFICATIONS_PAGE_SIZE),
    // The notification service returns a plain array, so a short page means
    // the end of the list was reached.
    getNextPageParam: (last, all) => (last.length < NOTIFICATIONS_PAGE_SIZE ? undefined : all.length + 1),
    enabled: !!userId,
  });
  const notifications = query.data?.pages.flatMap((page) => page) ?? [];

  // The query stays disabled until the canonical user id is resolved, so wait
  // for it before reporting that there are no notifications.
  if (!signedIn) return <LoadingState label="Checking your session" />;
  if (query.isLoading || !userResolved) return <LoadingState label="Loading alerts" />;
  if (query.isError) return <ErrorState body="Notification service is unavailable." onRetry={query.refetch} />;
  if (!notifications.length) return <EmptyState title="No notifications" body="Order and payment updates will show here." />;

  return (
    <Screen refreshControl={<RefreshControl refreshing={query.isRefetching} onRefresh={query.refetch} tintColor={colors.primary} />}>
      <Text variant="display">Alerts</Text>
      <View style={styles.list}>
        {notifications.map((item) => (
          <Card key={item.id} style={styles.card}>
            <Text variant="subtitle">{item.subject || item.type}</Text>
            <Text style={{ color: colors.muted }}>{item.body}</Text>
            <Text variant="caption" style={{ color: colors.primary }}>{item.channel} · {item.status}</Text>
          </Card>
        ))}
      </View>
      {query.hasNextPage && (
        <Pressable onPress={() => query.fetchNextPage()} style={styles.more}>
          <Text variant="label" style={{ color: colors.primary }}>
            {query.isFetchingNextPage ? "Loading..." : "Load more"}
          </Text>
        </Pressable>
      )}
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { gap: 12, marginTop: 18 },
  card: { gap: 6 },
  more: { alignItems: "center", paddingVertical: 18 },
});
