import { zodResolver } from "@hookform/resolvers/zod";
import { router } from "expo-router";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { StyleSheet } from "react-native";
import { z } from "zod";

import { getApiError } from "@/api/client";
import { authApi } from "@/api/endpoints";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";

const schema = z.object({
  first_name: z.string().min(2),
  last_name: z.string().min(1),
  email: z.string().email(),
  phone: z.string().min(7),
  password: z.string().min(6),
  city: z.string().min(2),
  pincode: z.string().min(4),
});

export default function RegisterScreen() {
  const { colors } = useAppTheme();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { control, handleSubmit, formState } = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: { first_name: "", last_name: "", email: "", phone: "", password: "", city: "", pincode: "" },
  });

  const onSubmit = handleSubmit(async (values) => {
    setLoading(true);
    setError(null);
    try {
      await authApi.register({
        email: values.email,
        phone: values.phone,
        password: values.password,
        first_name: values.first_name,
        last_name: values.last_name,
        dob: new Date("2000-01-01").toISOString(),
        gender: "",
        address: { city: values.city, country: "India", pincode: values.pincode },
      });
      router.replace("/(auth)/login");
    } catch (err) {
      setError(getApiError(err, "Registration failed"));
    } finally {
      setLoading(false);
    }
  });

  return (
    <Screen contentContainerStyle={styles.container}>
      <Text variant="display">Create account</Text>
      <Card style={styles.form}>
        {(["first_name", "last_name", "email", "phone", "password", "city", "pincode"] as const).map((name) => (
          <Controller
            key={name}
            control={control}
            name={name}
            render={({ field }) => (
              <Input
                label={name.replace("_", " ")}
                value={field.value}
                secureTextEntry={name === "password"}
                autoCapitalize={name === "email" ? "none" : "words"}
                keyboardType={name === "email" ? "email-address" : name === "phone" || name === "pincode" ? "number-pad" : "default"}
                onChangeText={field.onChange}
                error={formState.errors[name]?.message}
              />
            )}
          />
        ))}
        {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
        <Button title="Register" loading={loading} onPress={onSubmit} />
      </Card>
    </Screen>
  );
}

const styles = StyleSheet.create({
  container: { gap: 18 },
  form: { gap: 14 },
});
