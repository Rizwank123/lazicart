import { Tabs } from "expo-router";
import { Bell, Home, Package, ShoppingBag, UserRound } from "lucide-react-native";

import { useAppTheme } from "@/hooks/use-app-theme";

export default function TabsLayout() {
  const { colors } = useAppTheme();
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: colors.primary,
        tabBarInactiveTintColor: colors.muted,
        tabBarStyle: { backgroundColor: colors.surface, borderTopColor: colors.border, height: 64 },
        tabBarLabelStyle: { fontSize: 11, fontWeight: "700" },
      }}
    >
      <Tabs.Screen name="index" options={{ title: "Home", tabBarIcon: ({ color }) => <Home color={color} size={22} /> }} />
      <Tabs.Screen name="cart" options={{ title: "Cart", tabBarIcon: ({ color }) => <ShoppingBag color={color} size={22} /> }} />
      <Tabs.Screen name="orders" options={{ title: "Orders", tabBarIcon: ({ color }) => <Package color={color} size={22} /> }} />
      <Tabs.Screen name="notifications" options={{ title: "Alerts", tabBarIcon: ({ color }) => <Bell color={color} size={22} /> }} />
      <Tabs.Screen name="profile" options={{ title: "Profile", tabBarIcon: ({ color }) => <UserRound color={color} size={22} /> }} />
    </Tabs>
  );
}
