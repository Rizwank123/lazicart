import { Image } from "expo-image";
import { router, useLocalSearchParams } from "expo-router";
import { Heart, Star } from "lucide-react-native";
import { Pressable, StyleSheet, View } from "react-native";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Screen } from "@/components/ui/screen";
import { ErrorState, LoadingState } from "@/components/ui/state-view";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useCartActions } from "@/hooks/use-cart";
import { useProduct } from "@/hooks/use-products";
import { useAuthStore } from "@/store/auth-store";
import { useWishlistStore } from "@/store/wishlist-store";
import { formatMoney } from "@/utils/format";

export default function ProductDetailsScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { colors } = useAppTheme();
  const userId = useAuthStore((state) => state.userId);
  const { data: product, isLoading, isError, refetch } = useProduct(id);
  const add = useCartActions().add;
  const wished = useWishlistStore((state) => state.has(id));
  const toggle = useWishlistStore((state) => state.toggle);

  if (isLoading) return <LoadingState label="Opening product" />;
  if (isError || !product) return <ErrorState body="Product details are unavailable." onRetry={refetch} />;

  const image = product.images?.[0]?.image_url || product.images?.[0]?.thumbnail_url;
  const reviews = product.review ?? [];
  const average = reviews.length ? reviews.reduce((sum, review) => sum + review.rating, 0) / reviews.length : 0;

  return (
    <Screen>
      <View style={styles.media}>
        {image ? <Image source={{ uri: image }} style={styles.image} contentFit="cover" /> : <Text>No image</Text>}
        <Pressable onPress={() => toggle(product.id)} style={[styles.wish, { backgroundColor: colors.surface }]}>
          <Heart color={wished ? colors.danger : colors.muted} fill={wished ? colors.danger : "transparent"} />
        </Pressable>
      </View>
      <View style={styles.titleRow}>
        <View style={{ flex: 1 }}>
          <Text variant="display">{product.name}</Text>
          <Text style={{ color: colors.muted }}>{product.brand || "LahziCart"} · {product.currency}</Text>
        </View>
        <View style={styles.rating}>
          <Star color={colors.accent} fill={colors.accent} size={18} />
          <Text variant="label">{average ? average.toFixed(1) : "New"}</Text>
        </View>
      </View>
      <Text variant="title" style={{ color: colors.primary }}>{formatMoney(product.price, product.currency)}</Text>
      {!!product.description && <Text style={{ color: colors.muted }}>{product.description}</Text>}
      {Object.keys(product.attributes ?? {}).length > 0 && (
        <Card style={styles.attributes}>
          <Text variant="subtitle">Details</Text>
          {Object.entries(product.attributes).map(([key, value]) => (
            <View key={key} style={styles.attrRow}>
              <Text style={{ color: colors.muted }}>{key}</Text>
              <Text variant="label">{String(value)}</Text>
            </View>
          ))}
        </Card>
      )}
      <Button
        title={userId ? "Add to cart" : "Sign in to add"}
        loading={add.isPending}
        onPress={() => {
          if (!userId) return router.push("/(auth)/login");
          add.mutate({ product_id: product.id, quantity: 1, unit_price: product.price, currency: product.currency });
        }}
      />
      <Button title="Go to cart" variant="secondary" onPress={() => router.push("/(tabs)/cart")} />
      <View style={styles.reviews}>
        <Text variant="subtitle">Reviews</Text>
        {reviews.length ? reviews.slice(0, 5).map((review) => (
          <Card key={review.id} style={styles.review}>
            <Text variant="label">{review.rating}/5</Text>
            <Text style={{ color: colors.muted }}>{review.comment || "No comment"}</Text>
          </Card>
        )) : <Text style={{ color: colors.muted }}>No reviews yet.</Text>}
      </View>
    </Screen>
  );
}

const styles = StyleSheet.create({
  media: { height: 320, borderRadius: 8, overflow: "hidden", alignItems: "center", justifyContent: "center", marginBottom: 18 },
  image: { width: "100%", height: "100%" },
  wish: { position: "absolute", top: 12, right: 12, width: 44, height: 44, borderRadius: 22, alignItems: "center", justifyContent: "center" },
  titleRow: { flexDirection: "row", gap: 12, alignItems: "flex-start" },
  rating: { flexDirection: "row", alignItems: "center", gap: 4 },
  attributes: { gap: 10, marginTop: 10 },
  attrRow: { flexDirection: "row", justifyContent: "space-between", gap: 12 },
  reviews: { gap: 10, marginTop: 12 },
  review: { gap: 4 },
});
