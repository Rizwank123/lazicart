import { PropsWithChildren } from "react";
import { StyleSheet, Text as RNText, TextProps } from "react-native";

import { useAppTheme } from "@/hooks/use-app-theme";

type Variant = "display" | "title" | "subtitle" | "body" | "caption" | "label";

export function Text({ children, style, variant = "body", ...props }: PropsWithChildren<TextProps & { variant?: Variant }>) {
  const { colors } = useAppTheme();
  return (
    <RNText style={[styles.base, styles[variant], { color: colors.text }, style]} {...props}>
      {children}
    </RNText>
  );
}

const styles = StyleSheet.create({
  base: { letterSpacing: 0 },
  display: { fontSize: 30, lineHeight: 36, fontWeight: "800" },
  title: { fontSize: 22, lineHeight: 28, fontWeight: "800" },
  subtitle: { fontSize: 17, lineHeight: 24, fontWeight: "700" },
  body: { fontSize: 15, lineHeight: 22, fontWeight: "500" },
  caption: { fontSize: 12, lineHeight: 17, fontWeight: "600" },
  label: { fontSize: 13, lineHeight: 18, fontWeight: "800" },
});
