import { PropsWithChildren } from "react";
import { ScrollView, ScrollViewProps, StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { useAppTheme } from "@/hooks/use-app-theme";

type Props = PropsWithChildren<ScrollViewProps & { scroll?: boolean }>;

export function Screen({ children, scroll = true, contentContainerStyle, ...props }: Props) {
  const { colors } = useAppTheme();
  if (!scroll) {
    return (
      <SafeAreaView style={[styles.safe, { backgroundColor: colors.background }]}>
        <View style={styles.fill}>{children}</View>
      </SafeAreaView>
    );
  }
  return (
    <SafeAreaView style={[styles.safe, { backgroundColor: colors.background }]}>
      <ScrollView
        showsVerticalScrollIndicator={false}
        contentContainerStyle={[styles.content, contentContainerStyle]}
        {...props}
      >
        {children}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1 },
  fill: { flex: 1 },
  content: { padding: 20, paddingBottom: 40 },
});
