/** POST /auth/login */
export type LoginResponse = {
  access_token: string;
  token_type?: string;
  expires_in?: number;
};

/** GET /admin/channels */
export type ChannelListItem = {
  id: number;
  name: string;
  provider: string;
  base_url: string;
  api_key_preview: string;
  models: string[];
  model_mapping?: Record<string, string>;
  priority: number;
  weight: number;
  status: number;
  rate_limit: number;
};

export type ChannelsListResponse = { items: ChannelListItem[] };

/** POST/PUT channel body（与 OpenAPI ChannelPayload 对齐） */
export type ChannelPayload = {
  name: string;
  provider: string;
  base_url: string;
  api_key: string;
  models?: string[];
  model_mapping?: Record<string, string>;
  priority?: number;
  weight?: number;
  status?: number;
  rate_limit?: number;
};

/** 列表接口返回的兑换码行（Go 默认 JSON 为导出字段名 PascalCase） */
export type RedeemCodeRow = {
  ID: number;
  Code: string;
  Quota: number;
  UsedBy?: number | null;
  Status: number;
  ExpiresAt?: string | null;
  CreatedAt: string;
  UsedAt?: string | null;
};

export type RedeemListResponse = {
  items: RedeemCodeRow[];
  total: number;
  page: number;
  page_size: number;
};

export type RedeemStatsResponse = {
  by_status: Record<string, number>;
};

export type BatchRedeemResponse = {
  codes: string[];
  count: number;
};
