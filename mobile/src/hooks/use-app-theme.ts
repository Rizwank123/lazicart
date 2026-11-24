import { useColorScheme } from "react-native";

import { darkColors, lightColors } from "@/theme/colors";

export function useAppTheme() {
  const scheme = useColorScheme();
  const colors = scheme === "dark" ? darkColors : lightColors;
  return { colors, isDark: scheme === "dark" };
}
