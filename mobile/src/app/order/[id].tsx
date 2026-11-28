import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import { useLocalSearchParams } from "expo-router";
import { useState } from "react";
import { StyleSheet, View } from "react-native";

import { deliveryApi, orderApi, paymentApi } from "@/api/endpoints";
import { getApiError } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Screen } from "@/components/ui/screen";
import { ErrorState, LoadingState } from "@/components/ui/state-view";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useCancelOrder } from "@/hooks/use-orders";
import { useRequireAuth } from "@/hooks/use-require-auth";
import type { Payment, Shipment } from "@/types/api";
import { formatMoney, shortId } from "@/utils/format";

export default function OrderDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { colors } = useAppTheme();
  // Order, payment and shipment routes are JWT protected. Without a session
  // the queries would just return 401, so send the user to sign in instead.
  const { signedIn } = useRequireAuth();
  const enabled = !!id && signedIn;
  const order = useQuery({ queryKey: ["order", id], queryFn: () => orderApi.byId(id), enabled });
  const payment = useQuery({ queryKey: ["payment", id], queryFn: () => paymentApi.byOrder(id), enabled, retry: false });
  const shipment = useQuery({ queryKey: ["shipment", id], queryFn: () => deliveryApi.byOrder(id), enabled, retry: false });
  const cancel = useCancelOrder(id);
  const [cancelError, setCancelError] = useState<string | null>(null);

  if (!signedIn) return <LoadingState label="Checking your session" />;
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
        <PaymentStatus payment={payment} />
      </Card>
      <Card style={styles.section}>
        <Text variant="subtitle">Delivery</Text>
        <ShipmentStatus shipment={shipment} />
      </Card>
      {order.data.status !== "cancelled" && (
        <>
          {!!cancelError && <Text variant="caption" style={{ color: colors.danger }}>{cancelError}</Text>}
          <Button
            title="Cancel order"
            variant="danger"
            loading={cancel.isPending}
            onPress={async () => {
              setCancelError(null);
              try {
                await cancel.mutateAsync(order.data.id);
              } catch (err) {
                setCancelError(getApiError(err, "Could not cancel the order"));
              }
            }}
          />
        </>
      )}
    </Screen>
  );
}

const styles = StyleSheet.create({
  section: { gap: 12, marginTop: 16 },
  row: { flexDirection: "row", justifyContent: "space-between", gap: 12, alignItems: "center" },
  statusBlock: { gap: 8, alignItems: "flex-start" },
});

/**
 * The payment record.
 *
 * A missing record and a failed request are different situations, so they
 * are reported separately instead of both looking like "no payment yet".
 */
function PaymentStatus({ payment }: { payment: UseQueryResult<Payment, Error> }) {
  const { colors } = useAppTheme();
  if (payment.isLoading) return <Text style={{ color: colors.muted }}>Loading payment...</Text>;
  if (payment.isError) {
    return (
      <View style={styles.statusBlock}>
        <Text style={{ color: colors.danger }}>{getApiError(payment.error, "The payment service did not respond.")}</Text>
        <Button title="Retry" variant="secondary" onPress={() => payment.refetch()} />
      </View>
    );
  }
  if (!payment.data) return <Text style={{ color: colors.muted }}>No payment record yet</Text>;
  return <Text style={{ color: colors.muted }}>{`${payment.data.provider} · ${payment.data.status}`}</Text>;
}

/** The shipment record. */
function ShipmentStatus({ shipment }: { shipment: UseQueryResult<Shipment, Error> }) {
  const { colors } = useAppTheme();
  if (shipment.isLoading) return <Text style={{ color: colors.muted }}>Loading shipment...</Text>;
  if (shipment.isError) {
    return (
      <View style={styles.statusBlock}>
        <Text style={{ color: colors.danger }}>{getApiError(shipment.error, "The delivery service did not respond.")}</Text>
        <Button title="Retry" variant="secondary" onPress={() => shipment.refetch()} />
      </View>
    );
  }
  if (!shipment.data) return <Text style={{ color: colors.muted }}>Shipment not created yet</Text>;
  return (
    <Text style={{ color: colors.muted }}>
      {`${shipment.data.status}${shipment.data.tracking_number ? ` · ${shipment.data.tracking_number}` : ""}`}
    </Text>
  );
}
