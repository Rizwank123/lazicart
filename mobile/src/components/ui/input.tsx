import { forwardRef } from "react";
import { StyleSheet, TextInput, TextInputProps, View } from "react-native";

import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";

type Props = TextInputProps & { label: string; error?: string };

export const Input = forwardRef<TextInput, Props>(({ label, error, style, ...props }, ref) => {
  const { colors } = useAppTheme();
  return (
    <View style={styles.wrap}>
      <Text variant="caption" style={{ color: colors.muted }}>{label}</Text>
      <TextInput
        ref={ref}
        placeholderTextColor={colors.muted}
        style={[
          styles.input,
          { backgroundColor: colors.surface, borderColor: error ? colors.danger : colors.border, color: colors.text },
          style,
        ]}
        {...props}
      />
      {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
    </View>
  );
});

Input.displayName = "Input";

const styles = StyleSheet.create({
  wrap: { gap: 8 },
  input: {
    minHeight: 52,
    borderRadius: 8,
    borderWidth: 1,
    paddingHorizontal: 14,
    fontSize: 15,
    fontWeight: "600",
  },
});
