import { jwtDecode } from "jwt-decode";

type Claims = {
  sub?: string;
  role?: string;
  jti?: string;
  exp?: number;
};

export function decodeUserId(token?: string | null) {
  if (!token) return null;
  try {
    return jwtDecode<Claims>(token).sub ?? null;
  } catch {
    return null;
  }
}

export function isExpired(token?: string | null) {
  if (!token) return true;
  try {
    const exp = jwtDecode<Claims>(token).exp;
    return !exp || exp * 1000 < Date.now() + 30_000;
  } catch {
    return true;
  }
}
