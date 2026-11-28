import { router, useLocalSearchParams, type Href } from "expo-router";
import { useState } from "react";
import { StyleSheet } from "react-native";

import { getApiError } from "@/api/client";
import { authApi } from "@/api/endpoints";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useAuthStore } from "@/store/auth-store";

export default function OtpScreen() {
  const { colors } = useAppTheme();
  const setSession = useAuthStore((state) => state.setSession);
  // The screen that sent the visitor here is returned to after signing in.
  const { from } = useLocalSearchParams<{ from?: string }>();
  const [phone, setPhone] = useState("");
  const [otp, setOtp] = useState("");
  const [sent, setSent] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function sendOtp() {
    setLoading(true);
    setError(null);
    try {
      await authApi.sendOtp(phone);
      setSent(true);
    } catch (err) {
      setError(getApiError(err, "Could not send OTP"));
    } finally {
      setLoading(false);
    }
  }

  async function login() {
    setLoading(true);
    setError(null);
    try {
      const tokens = await authApi.loginWithPhone(phone, otp);
      await setSession(tokens);
      router.replace((from?.startsWith("/") ? from : "/(tabs)") as Href);
    } catch (err) {
      setError(getApiError(err, "OTP login failed"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <Screen contentContainerStyle={styles.container}>
      <Text variant="display">Phone sign in</Text>
      <Card style={styles.form}>
        <Input
          label="Phone"
          keyboardType="phone-pad"
          value={phone}
          // The phone number is locked once a code was sent, because the code
          // was delivered to that number. Changing it here would verify the
          // code against a different phone.
          editable={!sent}
          onChangeText={(value) => {
            setPhone(value);
            setOtp("");
            setSent(false);
          }}
        />
        {sent && <Input label="OTP" keyboardType="number-pad" value={otp} onChangeText={setOtp} maxLength={6} />}
        {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
        {sent ? (
          <>
            <Button title="Verify and continue" loading={loading} onPress={login} />
            <Button title="Resend OTP" variant="secondary" loading={loading} onPress={sendOtp} />
            <Button
              title="Use a different number"
              variant="ghost"
              onPress={() => {
                setSent(false);
                setOtp("");
              }}
            />
          </>
        ) : (
          <Button title="Send OTP" loading={loading} onPress={sendOtp} />
        )}
      </Card>
    </Screen>
  );
}

const styles = StyleSheet.create({
  container: { flexGrow: 1, justifyContent: "center", gap: 18 },
  form: { gap: 14 },
});
