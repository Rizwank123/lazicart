import { router } from "expo-router";
import { Pressable, RefreshControl, StyleSheet, View } from "react-native";

import { Card } from "@/components/ui/card";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state-view";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useOrders } from "@/hooks/use-orders";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { useAuthStore } from "@/store/auth-store";
import { formatMoney, shortId, unitCount } from "@/utils/format";

export default function OrdersScreen() {
  const { colors } = useAppTheme();
  const { data, isLoading, isError, refetch, isRefetching, fetchNextPage, hasNextPage, isFetchingNextPage } = useOrders();
  const userResolved = useAuthStore((state) => state.userResolved);
  const { signedIn } = useRequireAuth();
  const orders = data?.pages.flatMap((page) => page) ?? [];

  // The orders query stays disabled until the canonical user id is resolved, so
  // wait for it before reporting that there are no orders.
  if (!signedIn) return <LoadingState label="Checking your session" />;
  if (isLoading || !userResolved) return <LoadingState label="Loading orders" />;
  if (isError) return <ErrorState body="Order service is unavailable." onRetry={refetch} />;
  if (!orders.length) return <EmptyState title="No orders yet" body="Your checkout history will appear here." />;
  return (
    <Screen refreshControl={<RefreshControl refreshing={isRefetching} onRefresh={refetch} tintColor={colors.primary} />}>
      <Text variant="display">Orders</Text>
      <View style={styles.list}>
        {orders.map((order) => (
          <Pressable key={order.id} onPress={() => router.push(`/order/${order.id}`)}>
            <Card style={styles.card}>
              <View>
                <Text variant="subtitle">Order {shortId(order.id)}</Text>
                <Text style={{ color: colors.muted }}>{unitCount(order.items ?? [])} items · {order.status}</Text>
              </View>
              <Text variant="label" style={{ color: colors.primary }}>{formatMoney(order.total, order.currency)}</Text>
            </Card>
          </Pressable>
        ))}
      </View>
      {hasNextPage && (
        <Pressable onPress={() => fetchNextPage()} style={styles.more}>
          <Text variant="label" style={{ color: colors.primary }}>
            {isFetchingNextPage ? "Loading..." : "Load more"}
          </Text>
        </Pressable>
      )}
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { gap: 12, marginTop: 18 },
  card: { flexDirection: "row", justifyContent: "space-between", gap: 12, alignItems: "center" },
  more: { alignItems: "center", paddingVertical: 18 },
});
