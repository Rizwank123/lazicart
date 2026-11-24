import { router } from "expo-router";
import { useState } from "react";
import { StyleSheet, View } from "react-native";

import { getApiError } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Screen } from "@/components/ui/screen";
import { EmptyState, LoadingState } from "@/components/ui/state-view";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useCart, useCartActions } from "@/hooks/use-cart";
import { useCreateCheckout } from "@/hooks/use-orders";
import { useCheckoutStore } from "@/store/checkout-store";
import { formatMoney } from "@/utils/format";

export default function CheckoutScreen() {
  const { colors } = useAppTheme();
  const { data: cart, isLoading } = useCart();
  const clear = useCartActions().clear;
  const checkout = useCreateCheckout();
  const { address, setAddress } = useCheckoutStore();
  const [error, setError] = useState<string | null>(null);
  const items = cart?.items ?? [];
  const total = items.reduce((sum, item) => sum + item.unit_price * item.quantity, 0);
  const currency = items[0]?.currency ?? "INR";

  if (isLoading) return <LoadingState label="Preparing checkout" />;
  if (!items.length) return <EmptyState title="Nothing to checkout" />;

  async function placeOrder() {
    setError(null);
    try {
      const result = await checkout.mutateAsync({ address, items });
      await clear.mutateAsync();
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
      </Card>
      <Card style={styles.section}>
        <View style={styles.row}><Text style={{ color: colors.muted }}>Items</Text><Text variant="label">{items.length}</Text></View>
        <View style={styles.row}><Text style={{ color: colors.muted }}>Payment</Text><Text variant="label">Manual pending</Text></View>
        <View style={styles.row}><Text variant="subtitle">Total</Text><Text variant="title">{formatMoney(total, currency)}</Text></View>
      </Card>
      {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
      <Button title="Place order" loading={checkout.isPending || clear.isPending} onPress={placeOrder} />
    </Screen>
  );
}

const styles = StyleSheet.create({
  section: { gap: 12, marginTop: 16 },
  row: { flexDirection: "row", justifyContent: "space-between", gap: 12 },
});
