/** Format internal quota units per 1M tokens for display. */
export function formatQuotaPerMillion(value: number): string {
  if (value <= 0) return "—";
  if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(value % 1_000_000 === 0 ? 0 : 2)}M`;
  }
  if (value >= 1_000) {
    return `${(value / 1_000).toFixed(value % 1_000 === 0 ? 0 : 1)}K`;
  }
  return String(value);
}
