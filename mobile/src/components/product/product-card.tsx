import { Image } from "expo-image";
import { Heart, Star } from "lucide-react-native";
import { Pressable, StyleSheet, View } from "react-native";

import { Card } from "@/components/ui/card";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useWishlistStore } from "@/store/wishlist-store";
import type { Product } from "@/types/api";
import { formatMoney } from "@/utils/format";

type Props = {
  product: Product;
  onPress: () => void;
};

export function ProductCard({ product, onPress }: Props) {
  const { colors } = useAppTheme();
  const wished = useWishlistStore((state) => state.has(product.id));
  const toggle = useWishlistStore((state) => state.toggle);
  const image = product.images?.[0]?.thumbnail_url || product.images?.[0]?.image_url;
  const reviews = product.review ?? [];
  const avg = reviews.length ? reviews.reduce((sum, review) => sum + review.rating, 0) / reviews.length : 0;

  return (
    <Pressable onPress={onPress} style={styles.press}>
      <Card style={styles.card}>
        <View style={[styles.imageWrap, { backgroundColor: colors.surfaceMuted }]}>
          {image ? (
            <Image source={{ uri: image }} style={styles.image} contentFit="cover" transition={180} />
          ) : (
            <Text variant="caption" style={{ color: colors.muted }}>No image</Text>
          )}
          <Pressable accessibilityRole="button" onPress={() => toggle(product.id)} style={[styles.heart, { backgroundColor: colors.surface }]}>
            <Heart size={17} color={wished ? colors.danger : colors.muted} fill={wished ? colors.danger : "transparent"} />
          </Pressable>
        </View>
        <View style={styles.info}>
          <Text variant="caption" numberOfLines={1} style={{ color: colors.muted }}>{product.brand || "LahziCart"}</Text>
          <Text variant="subtitle" numberOfLines={2} style={styles.name}>{product.name}</Text>
          <View style={styles.row}>
            <Text variant="label" style={{ color: colors.primary }}>{formatMoney(product.price, product.currency)}</Text>
            <View style={styles.rating}>
              <Star size={13} color={colors.accent} fill={colors.accent} />
              <Text variant="caption" style={{ color: colors.muted }}>{avg ? avg.toFixed(1) : "New"}</Text>
            </View>
          </View>
        </View>
      </Card>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  press: { flex: 1 },
  card: { padding: 8, gap: 10 },
  imageWrap: { height: 138, borderRadius: 8, alignItems: "center", justifyContent: "center", overflow: "hidden" },
  image: { width: "100%", height: "100%" },
  heart: { position: "absolute", top: 8, right: 8, width: 34, height: 34, borderRadius: 17, alignItems: "center", justifyContent: "center" },
  info: { gap: 4, paddingHorizontal: 2, minHeight: 92 },
  name: { minHeight: 48 },
  row: { flexDirection: "row", alignItems: "center", justifyContent: "space-between", gap: 8 },
  rating: { flexDirection: "row", alignItems: "center", gap: 4 },
});
