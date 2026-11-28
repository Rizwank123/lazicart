import { useQueryClient } from "@tanstack/react-query";
import { router } from "expo-router";
import { useEffect, useState } from "react";
import { StyleSheet, View } from "react-native";

import { getApiError } from "@/api/client";
import { userApi } from "@/api/endpoints";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useRequireAuth } from "@/hooks/use-require-auth";
import { useAuthStore } from "@/store/auth-store";
import { initials } from "@/utils/format";

export default function ProfileScreen() {
  const { colors } = useAppTheme();
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const userId = useAuthStore((state) => state.userId);
  const userError = useAuthStore((state) => state.userError);
  const loadUser = useAuthStore((state) => state.loadUser);
  const logout = useAuthStore((state) => state.logout);
  const { signedIn } = useRequireAuth();
  const [firstName, setFirstName] = useState(user?.first_name ?? "");
  const [lastName, setLastName] = useState(user?.last_name ?? "");
  const [phone, setPhone] = useState(user?.phone ?? "");
  const [city, setCity] = useState(String(user?.address?.city ?? ""));
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  // The profile is fetched asynchronously, so the form has to be re-synced
  // once it arrives. Without this the fields stay empty because the
  // useState initializers only run on the first render.
  useEffect(() => {
    if (!user) return;
    setFirstName(user.first_name ?? "");
    setLastName(user.last_name ?? "");
    setPhone(user.phone ?? "");
    setCity(String(user.address?.city ?? ""));
  }, [user]);

  async function save() {
    if (!userId) return;
    setSaving(true);
    setError(null);
    try {
      await userApi.update(userId, {
        first_name: firstName,
        last_name: lastName,
        phone,
        address: { ...(user?.address ?? {}), city },
      });
      await loadUser();
    } catch (err) {
      setError(getApiError(err, "Profile update failed"));
    } finally {
      setSaving(false);
    }
  }

  async function signOut() {
    await logout();
    // The cached user data belongs to the account that just signed out.
    queryClient.clear();
    router.replace("/(auth)/login");
  }

  if (!signedIn) return null;

  return (
    <Screen>
      <View style={styles.hero}>
        <View style={[styles.avatar, { backgroundColor: colors.primary }]}>
          <Text variant="title" style={{ color: "#FFFFFF" }}>{initials(user?.first_name, user?.last_name, user?.email)}</Text>
        </View>
        <View style={{ flex: 1 }}>
          <Text variant="display">{user?.first_name || "Profile"}</Text>
          <Text style={{ color: colors.muted }}>{user?.email || userId}</Text>
        </View>
      </View>
      {!!userError && (
        <Text variant="caption" style={[styles.warning, { color: colors.danger }]}>
          {userError} Pull to retry from the orders tab.
        </Text>
      )}
      <Card style={styles.form}>
        <Input label="First name" value={firstName} onChangeText={setFirstName} />
        <Input label="Last name" value={lastName} onChangeText={setLastName} />
        <Input label="Phone" value={phone} onChangeText={setPhone} keyboardType="phone-pad" />
        <Input label="City" value={city} onChangeText={setCity} />
        {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
        <Button title="Save profile" loading={saving} onPress={save} disabled={!userId} />
      </Card>
      <Card style={styles.form}>
        <Text variant="subtitle">Settings</Text>
        <Text style={{ color: colors.muted }}>Privacy policy, terms, support, and app preferences are prepared as frontend sections; the backend does not expose content APIs for them yet.</Text>
      </Card>
      <Button title="Sign out" variant="secondary" onPress={signOut} />
    </Screen>
  );
}

const styles = StyleSheet.create({
  hero: { flexDirection: "row", alignItems: "center", gap: 14 },
  avatar: { width: 64, height: 64, borderRadius: 32, alignItems: "center", justifyContent: "center" },
  form: { gap: 12, marginTop: 18 },
  warning: { marginTop: 12 },
});
