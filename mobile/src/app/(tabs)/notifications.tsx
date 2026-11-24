import { useQuery } from "@tanstack/react-query";
import { StyleSheet, View } from "react-native";

import { notificationApi } from "@/api/endpoints";
import { Card } from "@/components/ui/card";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state-view";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useAuthStore } from "@/store/auth-store";

export default function NotificationsScreen() {
  const { colors } = useAppTheme();
  const userId = useAuthStore((state) => state.userId);
  const query = useQuery({
    queryKey: ["notifications", userId],
    queryFn: () => notificationApi.byUser(userId as string),
    enabled: !!userId,
  });

  if (query.isLoading) return <LoadingState label="Loading alerts" />;
  if (query.isError) return <ErrorState body="Notification service is unavailable." onRetry={query.refetch} />;
  if (!query.data?.length) return <EmptyState title="No notifications" body="Order and payment updates will show here." />;

  return (
    <Screen>
      <Text variant="display">Alerts</Text>
      <View style={styles.list}>
        {query.data.map((item) => (
          <Card key={item.id} style={styles.card}>
            <Text variant="subtitle">{item.subject || item.type}</Text>
            <Text style={{ color: colors.muted }}>{item.body}</Text>
            <Text variant="caption" style={{ color: colors.primary }}>{item.channel} · {item.status}</Text>
          </Card>
        ))}
      </View>
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { gap: 12, marginTop: 18 },
  card: { gap: 6 },
});
