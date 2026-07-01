export type DashboardStatsResponse = {
  success: boolean;
  data: {
    account: { quota: number; used_quota: number; group: string };
    usage: { request_count_24h: number; total_tokens_24h: number };
    consumption: { cost_quota_24h: number };
    performance: { avg_latency_ms: number; rpm: number; tpm: number };
    sparklines: {
      requests: { t: string; v: number }[];
      cost: { t: string; v: number }[];
    };
  };
};

export type DashboardChartsResponse = {
  success: boolean;
  data: {
    consumption_by_model: { model: string; cost: number; count: number }[];
    call_trend: { hour: string; count: number; cost: number; tokens: number }[];
    ranking: { model: string; cost: number; count: number }[];
  };
};

export type PortalNode = {
  name: string;
  url: string;
  region: string;
};

export type NodePingResult = PortalNode & {
  ok: boolean;
  latency_ms: number;
  status?: string;
};

export type NodesPingResponse = {
  success: boolean;
  data: NodePingResult[];
};

export type UsageLogRow = {
  id: number;
  request_id: string;
  token_name: string;
  token_group: string;
  model: string;
  request_path: string;
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  cost_quota: number;
  latency_ms: number | null;
  time_to_first_ms: number | null;
  billing_detail: unknown;
  status_code: number | null;
  created_at: string;
};

export type TaskLogRow = {
  id: number;
  task_id: string;
  platform: string;
  type: string;
  status: string;
  progress: number;
  detail: string;
  submitted_at: string;
  finished_at: string | null;
  duration_ms: number | null;
};
