export type WalletSummaryResponse = {
  success: boolean;
  data: {
    quota: number;
    used_quota: number;
    request_count: number;
    invite_code: string;
    invite_url: string;
    affiliate_pending: number;
    recharge_enabled: boolean;
  };
};

export type AnnouncementItem = {
  id: number;
  title: string;
  content: string;
  level: string;
};
