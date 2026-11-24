import { PropsWithChildren } from "react";
import { StyleSheet, View, ViewProps } from "react-native";

import { useAppTheme } from "@/hooks/use-app-theme";

export function Card({ children, style, ...props }: PropsWithChildren<ViewProps>) {
  const { colors } = useAppTheme();
  return (
    <View style={[styles.card, { backgroundColor: colors.surface, borderColor: colors.border }, style]} {...props}>
      {children}
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    borderRadius: 8,
    borderWidth: StyleSheet.hairlineWidth,
    padding: 14,
  },
});
