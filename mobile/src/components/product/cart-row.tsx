import { Image } from "expo-image";
import { Minus, Plus, Trash2 } from "lucide-react-native";
import { Pressable, StyleSheet, View } from "react-native";

import { Card } from "@/components/ui/card";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import type { CartItem, Product } from "@/types/api";
import { formatMoney } from "@/utils/format";

type Props = {
  item: CartItem;
  product?: Product;
  onIncrement: () => void;
  onDecrement: () => void;
  onRemove: () => void;
};

export function CartRow({ item, product, onIncrement, onDecrement, onRemove }: Props) {
  const { colors } = useAppTheme();
  const image = product?.images?.[0]?.thumbnail_url || product?.images?.[0]?.image_url;
  return (
    <Card style={styles.card}>
      <View style={[styles.imageWrap, { backgroundColor: colors.surfaceMuted }]}>
        {image ? <Image source={{ uri: image }} style={styles.image} contentFit="cover" /> : <Text variant="caption">Item</Text>}
      </View>
      <View style={styles.body}>
        <Text variant="subtitle" numberOfLines={2}>{product?.name ?? `Product ${item.product_id.slice(0, 8)}`}</Text>
        <Text variant="label" style={{ color: colors.primary }}>{formatMoney(item.unit_price * item.quantity, item.currency)}</Text>
        <View style={styles.actions}>
          <View style={[styles.stepper, { borderColor: colors.border }]}>
            <Pressable onPress={onDecrement} style={styles.iconButton}><Minus size={16} color={colors.text} /></Pressable>
            <Text variant="label">{item.quantity}</Text>
            <Pressable onPress={onIncrement} style={styles.iconButton}><Plus size={16} color={colors.text} /></Pressable>
          </View>
          <Pressable onPress={onRemove} style={styles.iconButton}><Trash2 size={18} color={colors.danger} /></Pressable>
        </View>
      </View>
    </Card>
  );
}

const styles = StyleSheet.create({
  card: { flexDirection: "row", gap: 12 },
  imageWrap: { width: 92, height: 104, borderRadius: 8, overflow: "hidden", alignItems: "center", justifyContent: "center" },
  image: { width: "100%", height: "100%" },
  body: { flex: 1, gap: 8 },
  actions: { flexDirection: "row", justifyContent: "space-between", alignItems: "center" },
  stepper: { flexDirection: "row", alignItems: "center", gap: 10, borderWidth: 1, borderRadius: 8, paddingHorizontal: 6 },
  iconButton: { width: 34, height: 34, alignItems: "center", justifyContent: "center" },
});
