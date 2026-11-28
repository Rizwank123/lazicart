import { useQueries } from "@tanstack/react-query";
import { router } from "expo-router";
import { useState } from "react";
import { StyleSheet, View } from "react-native";

import { getApiError } from "@/api/client";
import { productApi } from "@/api/endpoints";
import { CartRow } from "@/components/product/cart-row";
import { Button } from "@/components/ui/button";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state-view";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useCart, useCartActions } from "@/hooks/use-cart";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { useAuthStore } from "@/store/auth-store";
import { cartCurrencies, formatMoney, totalsByCurrency } from "@/utils/format";

export default function CartScreen() {
  const { colors } = useAppTheme();
  const { data: cart, isLoading, isError, refetch } = useCart();
  const userResolved = useAuthStore((state) => state.userResolved);
  const { signedIn } = useRequireAuth();
  const clear = useCartActions().clear;
  const [error, setError] = useState<string | null>(null);
  const items = cart?.items ?? [];
  const products = useQueries({
    queries: items.map((item) => ({
      queryKey: ["product", item.product_id],
      queryFn: () => productApi.detail(item.product_id),
      staleTime: 120_000,
    })),
  });
  const totals = totalsByCurrency(items);
  const currencies = cartCurrencies(items);
  const mixedCurrencies = currencies.length > 1;

  // The cart query stays disabled until the canonical user id is resolved, so
  // wait for it before reporting the cart as empty.
  if (!signedIn) return <LoadingState label="Checking your session" />;
  if (isLoading || !userResolved) return <LoadingState label="Loading cart" />;
  if (isError) return <ErrorState body="Cart service is unavailable." onRetry={refetch} />;
  if (!items.length) return <EmptyState title="Your cart is empty" body="Find something excellent on the home tab." />;

  async function clearCart() {
    setError(null);
    try {
      await clear.mutateAsync();
    } catch (err) {
      setError(getApiError(err, "The cart could not be cleared"));
    }
  }

  return (
    <Screen>
      <Text variant="display">Cart</Text>
      <View style={styles.list}>
        {items.map((item, index) => (
          <CartRow key={item.id} item={item} product={products[index]?.data} />
        ))}
      </View>
      <View style={[styles.summary, { borderColor: colors.border }]}>
        <Text style={{ color: colors.muted }}>Subtotal</Text>
        <View style={styles.totals}>
          {totals.map((entry) => (
            <Text key={entry.currency} variant="title">{formatMoney(entry.total, entry.currency)}</Text>
          ))}
        </View>
      </View>
      {mixedCurrencies && (
        <Text variant="caption" style={{ color: colors.danger }}>
          This cart mixes {currencies.join(" and ")}. Remove the items you do not want to order before checking out.
        </Text>
      )}
      {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
      <Button title="Checkout" disabled={mixedCurrencies} onPress={() => router.push("/checkout")} />
      <Button title="Clear cart" variant="ghost" onPress={clearCart} loading={clear.isPending} />
    </Screen>
  );
}

const styles = StyleSheet.create({
  list: { gap: 12, marginTop: 18 },
  summary: { marginTop: 18, paddingTop: 18, borderTopWidth: 1, flexDirection: "row", justifyContent: "space-between", alignItems: "center" },
  totals: { alignItems: "flex-end" },
});
