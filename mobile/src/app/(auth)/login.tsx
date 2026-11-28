import { zodResolver } from "@hookform/resolvers/zod";
import { Link, router, useLocalSearchParams, type Href } from "expo-router";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { Pressable, StyleSheet, View } from "react-native";
import { z } from "zod";

import { getApiError } from "@/api/client";
import { authApi } from "@/api/endpoints";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useAuthStore } from "@/store/auth-store";

const schema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
});

export default function LoginScreen() {
  const { colors } = useAppTheme();
  const setSession = useAuthStore((state) => state.setSession);
  // The screen that sent the visitor here is returned to after signing in.
  const { from } = useLocalSearchParams<{ from?: string }>();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { control, handleSubmit, formState } = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: { email: "", password: "" },
  });

  const onSubmit = handleSubmit(async (values) => {
    setError(null);
    setLoading(true);
    try {
      const tokens = await authApi.login(values.email, values.password);
      await setSession(tokens);
      router.replace((from?.startsWith("/") ? from : "/(tabs)") as Href);
    } catch (err) {
      setError(getApiError(err, "Login failed"));
    } finally {
      setLoading(false);
    }
  });

  return (
    <Screen contentContainerStyle={styles.container}>
      <View style={styles.header}>
        <Text variant="display">LahziCart</Text>
        <Text style={{ color: colors.muted }}>Fast shopping, clean checkout, real backend integration.</Text>
      </View>
      <Card style={styles.form}>
        <Controller
          control={control}
          name="email"
          render={({ field }) => (
            <Input label="Email" autoCapitalize="none" keyboardType="email-address" value={field.value} onChangeText={field.onChange} error={formState.errors.email?.message} />
          )}
        />
        <Controller
          control={control}
          name="password"
          render={({ field }) => (
            <Input label="Password" secureTextEntry value={field.value} onChangeText={field.onChange} error={formState.errors.password?.message} />
          )}
        />
        {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
        <Button title="Sign in" loading={loading} onPress={onSubmit} />
        <Link href="/(auth)/otp" asChild>
          <Pressable><Text variant="label" style={{ color: colors.primary, textAlign: "center" }}>Use phone OTP</Text></Pressable>
        </Link>
      </Card>
      <Link href="/(auth)/register" asChild>
        <Pressable><Text style={{ color: colors.muted, textAlign: "center" }}>New here? Create an account</Text></Pressable>
      </Link>
    </Screen>
  );
}

const styles = StyleSheet.create({
  container: { flexGrow: 1, justifyContent: "center", gap: 22 },
  header: { gap: 8 },
  form: { gap: 14 },
});
