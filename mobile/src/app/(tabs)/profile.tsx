import { router } from "expo-router";
import { useState } from "react";
import { StyleSheet, View } from "react-native";

import { getApiError } from "@/api/client";
import { userApi } from "@/api/endpoints";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Screen } from "@/components/ui/screen";
import { Text } from "@/components/ui/text";
import { useAppTheme } from "@/hooks/use-app-theme";
import { useAuthStore } from "@/store/auth-store";
import { initials } from "@/utils/format";

export default function ProfileScreen() {
  const { colors } = useAppTheme();
  const user = useAuthStore((state) => state.user);
  const userId = useAuthStore((state) => state.userId);
  const loadUser = useAuthStore((state) => state.loadUser);
  const logout = useAuthStore((state) => state.logout);
  const [firstName, setFirstName] = useState(user?.first_name ?? "");
  const [lastName, setLastName] = useState(user?.last_name ?? "");
  const [phone, setPhone] = useState(user?.phone ?? "");
  const [city, setCity] = useState(String(user?.address?.city ?? ""));
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

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
      <Card style={styles.form}>
        <Input label="First name" value={firstName} onChangeText={setFirstName} />
        <Input label="Last name" value={lastName} onChangeText={setLastName} />
        <Input label="Phone" value={phone} onChangeText={setPhone} keyboardType="phone-pad" />
        <Input label="City" value={city} onChangeText={setCity} />
        {!!error && <Text variant="caption" style={{ color: colors.danger }}>{error}</Text>}
        <Button title="Save profile" loading={saving} onPress={save} />
      </Card>
      <Card style={styles.form}>
        <Text variant="subtitle">Settings</Text>
        <Text style={{ color: colors.muted }}>Privacy policy, terms, support, and app preferences are prepared as frontend sections; the backend does not expose content APIs for them yet.</Text>
      </Card>
      <Button title="Sign out" variant="secondary" onPress={async () => { await logout(); router.replace("/(auth)/login"); }} />
    </Screen>
  );
}

const styles = StyleSheet.create({
  hero: { flexDirection: "row", alignItems: "center", gap: 14 },
  avatar: { width: 64, height: 64, borderRadius: 32, alignItems: "center", justifyContent: "center" },
  form: { gap: 12, marginTop: 18 },
});
