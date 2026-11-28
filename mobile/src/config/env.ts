import { Platform } from "react-native";

/**
 * The base URL of the API.
 *
 * On a device the app talks to the API gateway directly.
 *
 * On the web the browser talks to the development proxy instead. The gateway
 * does not answer the CORS preflight that the browser sends before a request
 * with a body, so a direct call would be blocked by the browser. The proxy
 * forwards the request to the gateway and adds the CORS headers. See
 * scripts/dev-api-proxy.js.
 */
function resolveApiBaseUrl(): string {
  if (Platform.OS === "web") {
    return process.env.EXPO_PUBLIC_WEB_API_BASE_URL ?? "http://localhost:7071";
  }
  return process.env.EXPO_PUBLIC_API_BASE_URL ?? "http://localhost:7070";
}

export const env = {
  apiBaseUrl: resolveApiBaseUrl(),
  defaultCurrency: process.env.EXPO_PUBLIC_DEFAULT_CURRENCY ?? "INR",
};
