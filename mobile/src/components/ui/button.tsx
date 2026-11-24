import { ActivityIndicator, Pressable, PressableProps, StyleSheet, ViewStyle } from "react-native";

import { useAppTheme } from "@/hooks/use-app-theme";
import { Text } from "@/components/ui/text";

type Props = PressableProps & {
  title: string;
  loading?: boolean;
  variant?: "primary" | "secondary" | "ghost" | "danger";
  style?: ViewStyle;
};

export function Button({ title, loading, variant = "primary", disabled, style, ...props }: Props) {
  const { colors } = useAppTheme();
  const bg = {
    primary: colors.primary,
    secondary: colors.surfaceMuted,
    ghost: "transparent",
    danger: colors.danger,
  }[variant];
  const fg = variant === "primary" || variant === "danger" ? "#FFFFFF" : colors.text;
  return (
    <Pressable
      accessibilityRole="button"
      disabled={disabled || loading}
      style={({ pressed }) => [
        styles.button,
        { backgroundColor: bg, opacity: disabled ? 0.5 : pressed ? 0.82 : 1 },
        style,
      ]}
      {...props}
    >
      {loading ? <ActivityIndicator color={fg} /> : <Text variant="label" style={{ color: fg }}>{title}</Text>}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  button: {
    minHeight: 52,
    borderRadius: 8,
    alignItems: "center",
    justifyContent: "center",
    paddingHorizontal: 18,
  },
});
