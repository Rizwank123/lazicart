import { ActivityIndicator, StyleSheet, View } from "react-native";
import { AlertCircle, PackageOpen } from "lucide-react-native";

import { Button } from "@/components/ui/button";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";

export function LoadingState({ label = "Loading LahziCart" }: { label?: string }) {
  const { colors } = useAppTheme();
  return (
    <View style={styles.center}>
      <ActivityIndicator color={colors.primary} size="large" />
      <Text variant="caption" style={{ color: colors.muted }}>{label}</Text>
    </View>
  );
}

export function EmptyState({ title, body }: { title: string; body?: string }) {
  const { colors } = useAppTheme();
  return (
    <View style={styles.center}>
      <PackageOpen color={colors.muted} size={36} />
      <Text variant="subtitle">{title}</Text>
      {!!body && <Text style={{ color: colors.muted, textAlign: "center" }}>{body}</Text>}
    </View>
  );
}

export function ErrorState({ title = "Could not load", body, onRetry }: { title?: string; body?: string; onRetry?: () => void }) {
  const { colors } = useAppTheme();
  return (
    <View style={styles.center}>
      <AlertCircle color={colors.danger} size={36} />
      <Text variant="subtitle">{title}</Text>
      {!!body && <Text style={{ color: colors.muted, textAlign: "center" }}>{body}</Text>}
      {!!onRetry && <Button title="Try again" onPress={onRetry} variant="secondary" />}
    </View>
  );
}

const styles = StyleSheet.create({
  center: { flex: 1, alignItems: "center", justifyContent: "center", gap: 12, padding: 24 },
});
