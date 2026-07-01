export type ApiTokenRow = {
  id: number;
  name: string;
  status: number;
  enabled: boolean;
  key_prefix: string;
  key_masked: string;
  token_group: string;
  quota_limit: number | null;
  used_quota: number;
  models: string[];
  ip_whitelist: string[];
  rate_limit: number;
  created_at: string;
};

export type TokenListResponse = {
  success: boolean;
  data: {
    items: ApiTokenRow[];
    total: number;
    page: number;
    page_size: number;
  };
};

export type TokenCreateResponse = {
  success: boolean;
  data: {
    token: ApiTokenRow;
    api_key: string;
  };
};
