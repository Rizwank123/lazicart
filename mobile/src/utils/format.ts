import type { Review } from "@/types/api";

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

type Priced = { unit_price: number; quantity: number; currency?: string };

/**
 * Total the items of a cart, grouped by currency.
 *
 * Adding prices of different currencies together produces a meaningless
 * number, so every currency keeps its own total.
 */
export function totalsByCurrency(items: Priced[]) {
  const totals = new Map<string, number>();
  for (const item of items) {
    const currency = item.currency || "INR";
    totals.set(currency, (totals.get(currency) ?? 0) + item.unit_price * item.quantity);
  }
  return Array.from(totals, ([currency, total]) => ({ currency, total }));
}

/** The currencies present in a cart. More than one means the cart cannot be ordered. */
export function cartCurrencies(items: Priced[]) {
  return Array.from(new Set(items.map((item) => item.currency || "INR")));
}

/** The number of units in a cart, not the number of line items. */
export function unitCount(items: Priced[]) {
  return items.reduce((sum, item) => sum + item.quantity, 0);
}

/**
 * The reviews that may be shown.
 *
 * A review that is still pending moderation is not public, and it must not
 * influence the displayed rating either.
 */
export function visibleReviews(reviews?: Review[]) {
  return (reviews ?? []).filter((review) => review.is_approved !== false);
}

/** The average rating of the visible reviews. */
export function averageRating(reviews?: Review[]) {
  const shown = visibleReviews(reviews);
  if (!shown.length) return 0;
  return shown.reduce((sum, review) => sum + review.rating, 0) / shown.length;
}
