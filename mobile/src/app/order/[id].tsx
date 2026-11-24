import { useQuery } from "@tanstack/react-query";
import { useLocalSearchParams } from "expo-router";
import { StyleSheet, View } from "react-native";

import { deliveryApi, orderApi, paymentApi } from "@/api/endpoints";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Screen } from "@/components/ui/screen";
import { ErrorState, LoadingState } from "@/components/ui/state-view";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { formatMoney, shortId } from "@/utils/format";

export default function OrderDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { colors } = useAppTheme();
  const order = useQuery({ queryKey: ["order", id], queryFn: () => orderApi.byId(id), enabled: !!id });
  const payment = useQuery({ queryKey: ["payment", id], queryFn: () => paymentApi.byOrder(id), enabled: !!id, retry: false });
  const shipment = useQuery({ queryKey: ["shipment", id], queryFn: () => deliveryApi.byOrder(id), enabled: !!id, retry: false });
  const cancel = useQuery({ queryKey: ["noop"], enabled: false, queryFn: async () => null });

  if (order.isLoading) return <LoadingState label="Loading order" />;
  if (order.isError || !order.data) return <ErrorState body="Order details are unavailable." onRetry={order.refetch} />;

  return (
    <Screen>
      <Text variant="display">Order {shortId(order.data.id)}</Text>
      <Card style={styles.section}>
        <View style={styles.row}><Text style={{ color: colors.muted }}>Status</Text><Text variant="label">{order.data.status}</Text></View>
        <View style={styles.row}><Text style={{ color: colors.muted }}>Total</Text><Text variant="title">{formatMoney(order.data.total, order.data.currency)}</Text></View>
      </Card>
      <Card style={styles.section}>
        <Text variant="subtitle">Items</Text>
        {(order.data.items ?? []).map((item) => (
          <View key={item.id} style={styles.row}>
            <Text style={{ color: colors.muted }}>{shortId(item.product_id)} x {item.quantity}</Text>
            <Text variant="label">{formatMoney(item.unit_price * item.quantity, item.currency)}</Text>
          </View>
        ))}
      </Card>
      <Card style={styles.section}>
        <Text variant="subtitle">Payment</Text>
        <Text style={{ color: colors.muted }}>{payment.data ? `${payment.data.provider} · ${payment.data.status}` : "No payment record yet"}</Text>
      </Card>
      <Card style={styles.section}>
        <Text variant="subtitle">Delivery</Text>
        <Text style={{ color: colors.muted }}>
          {shipment.data ? `${shipment.data.status}${shipment.data.tracking_number ? ` · ${shipment.data.tracking_number}` : ""}` : "Shipment not created yet"}
        </Text>
      </Card>
      {order.data.status !== "cancelled" && (
        <Button
          title="Cancel order"
          variant="danger"
          loading={cancel.isFetching}
          onPress={async () => {
            await orderApi.cancel(order.data.id);
            order.refetch();
          }}
        />
      )}
    </Screen>
  );
}

const styles = StyleSheet.create({
  section: { gap: 12, marginTop: 16 },
  row: { flexDirection: "row", justifyContent: "space-between", gap: 12, alignItems: "center" },
});
