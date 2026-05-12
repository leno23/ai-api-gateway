const ACCESS_TOKEN_KEY = "ai_gateway_admin_access_token";

export const authStorage = {
  getToken(): string | null {
    if (typeof window === "undefined") return null;
    return sessionStorage.getItem(ACCESS_TOKEN_KEY);
  },
  setToken(token: string) {
    sessionStorage.setItem(ACCESS_TOKEN_KEY, token);
  },
  clearToken() {
    sessionStorage.removeItem(ACCESS_TOKEN_KEY);
  },
};
