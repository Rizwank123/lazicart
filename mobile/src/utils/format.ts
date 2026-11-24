export function formatMoney(value: number, currency = "INR") {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency,
    maximumFractionDigits: value % 1 === 0 ? 0 : 2,
  }).format(value);
}

export function initials(first?: string, last?: string, email?: string) {
  const source = [first, last].filter(Boolean).join(" ") || email || "LC";
  return source
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}

export function shortId(id?: string) {
  if (!id) return "";
  return id.slice(0, 8).toUpperCase();
}
