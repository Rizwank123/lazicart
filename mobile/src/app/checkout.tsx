import { router } from "expo-router";
import { useState } from "react";
import { StyleSheet, View } from "react-native";

import { getApiError } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Screen } from "@/components/ui/screen";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state-view";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useCart, useCartActions } from "@/hooks/use-cart";
import { useCreateCheckout } from "@/hooks/use-orders";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { useCheckoutStore } from "@/store/checkout-store";
import { cartCurrencies, formatMoney, totalsByCurrency, unitCount } from "@/utils/format";

export default function CheckoutScreen() {
  const { colors } = useAppTheme();
  const { data: cart, isLoading, isError, refetch } = useCart();
  const clear = useCartActions().clear;
  const checkout = useCreateCheckout();
  const { address, setAddress } = useCheckoutStore();
  const { signedIn } = useRequireAuth();
  const [error, setError] = useState<string | null>(null);
  const [warning, setWarning] = useState<string | null>(null);
  const [createdOrderId, setCreatedOrderId] = useState<string | null>(null);
  const items = cart?.items ?? [];
  const totals = totalsByCurrency(items);
  const currencies = cartCurrencies(items);
  const mixedCurrencies = currencies.length > 1;

  const addressComplete = Boolean(
    String(address.city ?? "").trim() &&
      String(address.state ?? "").trim() &&
      String(address.country ?? "").trim() &&
      String(address.pincode ?? "").trim(),
  );

  if (!signedIn) return <LoadingState label="Checking your session" />;
  if (isLoading) return <LoadingState label="Preparing checkout" />;
  if (isError) return <ErrorState body="The cart service is unavailable." onRetry={refetch} />;
  if (!items.length) return <EmptyState title="Nothing to checkout" body="Add something to your cart first." />;

  async function placeOrder() {
    setError(null);
    setWarning(null);
    setCreatedOrderId(null);
    try {
      const result = await checkout.mutateAsync({ address, items });
      // Clearing the cart is housekeeping. A failure here must not turn a
      // completed order into a failed checkout, so the order is kept and the
      // problem is reported separately.
      try {
        await clear.mutateAsync();
      } catch (err) {
        setWarning(`Your order was placed, but the cart could not be cleared: ${getApiError(err)}`);
        setCreatedOrderId(result.order.id);
        return;
      }
      router.replace(`/order/${result.order.id}`);
    } catch (err) {
      setError(getApiError(err, "Checkout failed"));
    }
  }

  return (
    <Screen>
      <Text variant="display">Checkout</Text>
      <Card style={styles.section}>
        <Text variant="subtitle">Delivery address</Text>
        <Input label="City" value={String(address.city ?? "")} onChangeText={(city) => setAddress({ ...address, city })} />
        <Input label="State" value={String(address.state ?? "")} onChangeText={(state) => setAddress({ ...address, state })} />
        <Input label="Country" value={String(address.country ?? "")} onChangeText={(country) => setAddress({ ...address, country })} />
        <Input label="Pincode" keyboardType="number-pad" value={String(address.pincode ?? "")} onChangeText={(pincode) => setAddress({ ...address, pincode })} />
        {!addressComplete && (
          <Text variant="caption" style={{ color: colors.muted }}>
            City, state, country and pincode are required to place the order.
          </Text>
        )}
      </Card>
      <Card style={styles.section}>
        <View style={styles.row}><Text style={{ color: colors.muted }}>Items</Text><Text variant="label">{unitCount(items)}</Text></View>
        <View style={styles.row}><Text style={{ color: colors.muted }}>Payment</Text><Text variant="label">Manual pending</Text></View>
        {totals.map((entry) => (
          <View key={entry.currency} style={styles.row}>
            <Text style={{ color: colors.muted }}>Total</Text>
            <Text variant="title">{formatMoney(entry.total, entry.currency)}</Text>
          </View>
        ))}
      </Card>
      {mixedCurrencies && (
        <Text variant="caption" style={{ color: colors.danger }}>
          This cart mixes {currencies.join(" and ")}. Remove the items you do not want to order before checking out.
        </Text>
      )}
      {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
      {!!warning && <Text variant="caption" style={{ color: colors.muted }}>{warning}</Text>}
      {!!createdOrderId && (
        <Button title="View my order" variant="secondary" onPress={() => router.replace(`/order/${createdOrderId}`)} />
      )}
      <Button
        title="Place order"
        loading={checkout.isPending}
        disabled={!addressComplete || mixedCurrencies}
        onPress={placeOrder}
      />
    </Screen>
  );
}

const styles = StyleSheet.create({
  section: { gap: 12, marginTop: 16 },
  row: { flexDirection: "row", justifyContent: "space-between", gap: 12 },
});
