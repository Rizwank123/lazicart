import { useQueries } from "@tanstack/react-query";
import { router } from "expo-router";
import { StyleSheet, View } from "react-native";

import { productApi } from "@/api/endpoints";
import { CartRow } from "@/components/product/cart-row";
import { Button } from "@/components/ui/button";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state-view";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useCart, useCartActions } from "@/hooks/use-cart";
import { formatMoney } from "@/utils/format";

export default function CartScreen() {
  const { colors } = useAppTheme();
  const { data: cart, isLoading, isError, refetch } = useCart();
  const actions = useCartActions();
  const items = cart?.items ?? [];
  const products = useQueries({
    queries: items.map((item) => ({
      queryKey: ["product", item.product_id],
      queryFn: () => productApi.detail(item.product_id),
      staleTime: 120_000,
    })),
  });
  const total = items.reduce((sum, item) => sum + item.unit_price * item.quantity, 0);
  const currency = items[0]?.currency ?? "INR";

  if (isLoading) return <LoadingState label="Loading cart" />;
  if (isError) return <ErrorState body="Cart service is unavailable." onRetry={refetch} />;
  if (!items.length) return <EmptyState title="Your cart is empty" body="Find something excellent on the home tab." />;

  return (
    <Screen>
      <Text variant="display">Cart</Text>
      <View style={styles.list}>
        {items.map((item, index) => (
          <CartRow
            key={item.id}
            item={item}
            product={products[index]?.data}
            onIncrement={() => actions.update.mutate({ itemId: item.id, quantity: item.quantity + 1 })}
            onDecrement={() => item.quantity > 1 ? actions.update.mutate({ itemId: item.id, quantity: item.quantity - 1 }) : actions.remove.mutate(item.id)}
            onRemove={() => actions.remove.mutate(item.id)}
          />
        ))}
      </View>
      <View style={[styles.summary, { borderColor: colors.border }]}>
        <Text style={{ color: colors.muted }}>Subtotal</Text>
        <Text variant="title">{formatMoney(total, currency)}</Text>
      </View>
      <Button title="Checkout" onPress={() => router.push("/checkout")} />
      <Button title="Clear cart" variant="ghost" onPress={() => actions.clear.mutate()} loading={actions.clear.isPending} />
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { gap: 12, marginTop: 18 },
  summary: { marginTop: 18, paddingTop: 18, borderTopWidth: 1, flexDirection: "row", justifyContent: "space-between", alignItems: "center" },
});
