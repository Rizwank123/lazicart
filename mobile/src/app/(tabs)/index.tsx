import { FlashList } from "@shopify/flash-list";
import { router } from "expo-router";
import { Search, SlidersHorizontal } from "lucide-react-native";
import { useEffect, useMemo, useState } from "react";
import { Pressable, RefreshControl, StyleSheet, TextInput, View } from "react-native";

import { ProductCard } from "@/components/product/product-card";
import { Card } from "@/components/ui/card";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state-view";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { flattenProducts, useProducts } from "@/hooks/use-products";

export default function HomeScreen() {
  const { colors } = useAppTheme();
  const [search, setSearch] = useState("");
  const { data, isLoading, isError, refetch, fetchNextPage, hasNextPage, isFetchingNextPage, isRefetching } = useProducts();
  const products = flattenProducts(data?.pages);
  const needle = search.trim().toLowerCase();

  // The catalog service has no search endpoint, so the remaining pages are
  // loaded while a search is active. Otherwise a match on a later page would
  // be reported as missing.
  useEffect(() => {
    if (!needle) return;
    if (hasNextPage && !isFetchingNextPage) fetchNextPage();
  }, [needle, hasNextPage, isFetchingNextPage, fetchNextPage]);

  const filtered = useMemo(() => {
    if (!needle) return products;
    return products.filter((product) =>
      [product.name, product.brand, product.description].filter(Boolean).join(" ").toLowerCase().includes(needle),
    );
  }, [products, needle]);
  const brands = [...new Set(products.map((product) => product.brand).filter(Boolean))].slice(0, 6);

  if (isLoading) return <LoadingState label="Curating shelves" />;
  if (isError) return <ErrorState body="The product service did not respond." onRetry={refetch} />;

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <View style={styles.header}>
        <Text variant="display">LahziCart</Text>
        <Text style={{ color: colors.muted }}>Fresh finds from your backend catalog.</Text>
      </View>
      <Card style={styles.search}>
        <Search color={colors.muted} size={20} />
        <TextInput
          value={search}
          onChangeText={setSearch}
          placeholder="Search products, brands"
          placeholderTextColor={colors.muted}
          style={[styles.searchInput, { color: colors.text }]}
        />
        <SlidersHorizontal color={colors.primary} size={20} />
      </Card>
      {brands.length > 0 && (
        <View style={styles.chips}>
          {brands.map((brand) => (
            <Pressable key={brand} onPress={() => setSearch(String(brand))} style={[styles.chip, { backgroundColor: colors.surfaceMuted }]}>
              <Text variant="caption">{brand}</Text>
            </Pressable>
          ))}
        </View>
      )}
      <FlashList
        data={filtered}
        keyExtractor={(item) => item.id}
        numColumns={2}
        contentContainerStyle={styles.list}
        ItemSeparatorComponent={() => <View style={{ height: 12 }} />}
        ListEmptyComponent={
          needle ? (
            <EmptyState title="No matches" body={`Nothing in the catalog matches "${search.trim()}".`} />
          ) : (
            <EmptyState title="No products yet" body="Add products in the catalog service and they will appear here." />
          )
        }
        refreshControl={<RefreshControl refreshing={isRefetching} onRefresh={refetch} tintColor={colors.primary} />}
        onEndReached={() => {
          if (hasNextPage && !isFetchingNextPage) fetchNextPage();
        }}
        renderItem={({ item }) => <ProductCard product={item} onPress={() => router.push(`/product/${item.id}`)} />}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, paddingTop: 58 },
  header: { paddingHorizontal: 20, gap: 4 },
  search: { margin: 20, flexDirection: "row", alignItems: "center", gap: 10, paddingVertical: 10 },
  searchInput: { flex: 1, fontSize: 15, fontWeight: "600" },
  chips: { flexDirection: "row", gap: 8, paddingHorizontal: 20, paddingBottom: 12, flexWrap: "wrap" },
  chip: { paddingHorizontal: 12, paddingVertical: 8, borderRadius: 999 },
  list: { paddingHorizontal: 14, paddingBottom: 24 },
});
