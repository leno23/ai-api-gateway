export function getDocsUrl(): string {
  return (process.env.NEXT_PUBLIC_DOCS_URL ?? "").trim();
}

export function getPortalOrigin(): string {
  return (process.env.NEXT_PUBLIC_PORTAL_ORIGIN ?? "").trim();
}
