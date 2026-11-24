import { router } from "expo-router";
import { Pressable, StyleSheet, View } from "react-native";

import { Card } from "@/components/ui/card";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state-view";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useOrders } from "@/hooks/use-orders";
import { formatMoney, shortId } from "@/utils/format";

export default function OrdersScreen() {
  const { colors } = useAppTheme();
  const { data, isLoading, isError, refetch } = useOrders();
  if (isLoading) return <LoadingState label="Loading orders" />;
  if (isError) return <ErrorState body="Order service is unavailable." onRetry={refetch} />;
  if (!data?.length) return <EmptyState title="No orders yet" body="Your checkout history will appear here." />;
  return (
    <Screen>
      <Text variant="display">Orders</Text>
      <View style={styles.list}>
        {data.map((order) => (
          <Pressable key={order.id} onPress={() => router.push(`/order/${order.id}`)}>
            <Card style={styles.card}>
              <View>
                <Text variant="subtitle">Order {shortId(order.id)}</Text>
                <Text style={{ color: colors.muted }}>{order.items?.length ?? 0} items · {order.status}</Text>
              </View>
              <Text variant="label" style={{ color: colors.primary }}>{formatMoney(order.total, order.currency)}</Text>
            </Card>
          </Pressable>
        ))}
      </View>
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { gap: 12, marginTop: 18 },
  card: { flexDirection: "row", justifyContent: "space-between", gap: 12, alignItems: "center" },
});
