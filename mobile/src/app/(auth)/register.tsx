import { zodResolver } from "@hookform/resolvers/zod";
import { router } from "expo-router";
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

const schema = z.object({
  first_name: z.string().min(2),
  last_name: z.string().min(1),
  email: z.string().email(),
  phone: z.string().min(7),
  password: z.string().min(6),
  dob: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Use the YYYY-MM-DD format"),
  gender: z.string().min(1),
  city: z.string().min(2),
  state: z.string().min(2),
  pincode: z.string().min(4),
});

const GENDERS = ["female", "male", "other"];

export default function RegisterScreen() {
  const { colors } = useAppTheme();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { control, handleSubmit, formState, watch, setValue } = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      first_name: "",
      last_name: "",
      email: "",
      phone: "",
      password: "",
      dob: "",
      gender: "",
      city: "",
      state: "",
      pincode: "",
    },
  });
  const gender = watch("gender");

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
        // The date of birth is sent as the actual value that was entered, not
        // as a placeholder.
        dob: new Date(`${values.dob}T00:00:00.000Z`).toISOString(),
        gender: values.gender,
        address: { city: values.city, state: values.state, country: "India", pincode: values.pincode },
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
        <Controller
          control={control}
          name="first_name"
          render={({ field }) => <Input label="first name" value={field.value} onChangeText={field.onChange} error={formState.errors.first_name?.message} />}
        />
        <Controller
          control={control}
          name="last_name"
          render={({ field }) => <Input label="last name" value={field.value} onChangeText={field.onChange} error={formState.errors.last_name?.message} />}
        />
        <Controller
          control={control}
          name="email"
          render={({ field }) => (
            <Input label="email" value={field.value} autoCapitalize="none" keyboardType="email-address" onChangeText={field.onChange} error={formState.errors.email?.message} />
          )}
        />
        <Controller
          control={control}
          name="phone"
          render={({ field }) => (
            <Input label="phone" value={field.value} keyboardType="phone-pad" onChangeText={field.onChange} error={formState.errors.phone?.message} />
          )}
        />
        <Controller
          control={control}
          name="password"
          render={({ field }) => (
            <Input label="password" value={field.value} secureTextEntry onChangeText={field.onChange} error={formState.errors.password?.message} />
          )}
        />
        <Controller
          control={control}
          name="dob"
          render={({ field }) => (
            <Input
              label="date of birth (YYYY-MM-DD)"
              value={field.value}
              placeholder="1995-04-12"
              autoCapitalize="none"
              onChangeText={field.onChange}
              error={formState.errors.dob?.message}
            />
          )}
        />
        <View style={styles.gender}>
          <Text variant="caption" style={{ color: colors.muted }}>gender</Text>
          <View style={styles.genderOptions}>
            {GENDERS.map((option) => (
              <Pressable
                key={option}
                onPress={() => setValue("gender", option, { shouldValidate: true })}
                style={[
                  styles.genderChip,
                  {
                    backgroundColor: gender === option ? colors.primary : colors.surfaceMuted,
                    borderColor: colors.border,
                  },
                ]}
              >
                <Text variant="caption" style={{ color: gender === option ? "#FFFFFF" : colors.text }}>{option}</Text>
              </Pressable>
            ))}
          </View>
          {!!formState.errors.gender?.message && (
            <Text variant="caption" style={{ color: colors.danger }}>{formState.errors.gender.message}</Text>
          )}
        </View>
        <Controller
          control={control}
          name="city"
          render={({ field }) => <Input label="city" value={field.value} onChangeText={field.onChange} error={formState.errors.city?.message} />}
        />
        <Controller
          control={control}
          name="state"
          render={({ field }) => <Input label="state" value={field.value} onChangeText={field.onChange} error={formState.errors.state?.message} />}
        />
        <Controller
          control={control}
          name="pincode"
          render={({ field }) => (
            <Input label="pincode" value={field.value} keyboardType="number-pad" onChangeText={field.onChange} error={formState.errors.pincode?.message} />
          )}
        />
        {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
        <Button title="Register" loading={loading} onPress={onSubmit} />
      </Card>
    </Screen>
  );
}

const styles = StyleSheet.create({
  container: { gap: 18 },
  form: { gap: 14 },
  gender: { gap: 8 },
  genderOptions: { flexDirection: "row", gap: 8 },
  genderChip: { paddingHorizontal: 14, paddingVertical: 10, borderRadius: 999, borderWidth: 1 },
});
