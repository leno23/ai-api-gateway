export type PortalLocale = "zh-CN" | "en";

export type MessageTree = {
  nav: {
    home: string;
    console: string;
    pricing: string;
    docs: string;
    about: string;
    login: string;
    register: string;
    settings: string;
    logout: string;
    notifications: string;
    language: string;
  };
  home: {
    heroTitle: string;
    heroSubtitle: string;
    ctaRegister: string;
    ctaLogin: string;
    partnersTitle: string;
    feature1Title: string;
    feature1Desc: string;
    feature2Title: string;
    feature2Desc: string;
    feature3Title: string;
    feature3Desc: string;
  };
  about: {
    title: string;
    intro: string;
    missionTitle: string;
    mission: string;
    featureTitle: string;
    features: string[];
    contactTitle: string;
    contact: string;
  };
  docs: {
    title: string;
    intro: string;
    openExternal: string;
    quickStart: string;
    authTitle: string;
    authBody: string;
    baseUrlTitle: string;
    baseUrlBody: string;
    errorsTitle: string;
    errorsBody: string;
  };
  dashboard: {
    ping: string;
    pinging: string;
    pingOk: string;
    pingFail: string;
  };
};

export const messages: Record<PortalLocale, MessageTree> = {
  "zh-CN": {
    nav: {
      home: "首页",
      console: "控制台",
      pricing: "模型广场",
      docs: "文档",
      about: "关于",
      login: "登录",
      register: "注册",
      settings: "个人设置",
      logout: "退出登录",
      notifications: "通知",
      language: "语言",
    },
    home: {
      heroTitle: "极速连通，无界创造",
      heroSubtitle: "企业级多模型接入网关",
      ctaRegister: "立即注册",
      ctaLogin: "登录控制台",
      partnersTitle: "合作伙伴与模型生态",
      feature1Title: "高并发网关",
      feature1Desc: "无锁队列与连接池，稳定承载企业级调用峰值。",
      feature2Title: "极致性价比",
      feature2Desc: "按量计费、分组倍率透明，成本可控。",
      feature3Title: "多模型统一接入",
      feature3Desc: "OpenAI / Anthropic 兼容端点，一套密钥走天下。",
    },
    about: {
      title: "关于我们",
      intro:
        "AI API Gateway 是企业级多模型 API 聚合与转售运营平台，为团队提供统一鉴权、额度管理、渠道调度与可观测能力。",
      missionTitle: "产品使命",
      mission: "让开发者在同一套 OpenAI 兼容协议下，安全、低成本地调用全球主流大模型。",
      featureTitle: "核心能力",
      features: [
        "OpenAI / Anthropic 兼容网关与流式对话",
        "令牌分组、模型白名单、IP 限制与用量日志",
        "钱包、兑换码、邀请返利与运营公告",
        "多区域 API 节点与可用性探测",
      ],
      contactTitle: "联系我们",
      contact: "商务与合作请通过控制台公告或企业微信二维码（见首页弹窗）。",
    },
    docs: {
      title: "开发者文档",
      intro: "接入网关前，请先创建令牌并配置 Base URL。",
      openExternal: "打开完整文档站",
      quickStart: "快速开始",
      authTitle: "鉴权",
      authBody: "在请求头携带 Authorization: Bearer sk-xxx（于「令牌管理」创建）。",
      baseUrlTitle: "Base URL",
      baseUrlBody: "将 SDK 的 base_url 指向网关地址，例如 https://api.example.com/v1。",
      errorsTitle: "常见错误",
      errorsBody: "402 表示额度不足；403 可能为模型白名单或 IP 限制；429 为限流。",
    },
    dashboard: {
      ping: "测速",
      pinging: "测速中…",
      pingOk: "可用",
      pingFail: "不可用",
    },
  },
  en: {
    nav: {
      home: "Home",
      console: "Console",
      pricing: "Models",
      docs: "Docs",
      about: "About",
      login: "Sign in",
      register: "Register",
      settings: "Settings",
      logout: "Sign out",
      notifications: "Notifications",
      language: "Language",
    },
    home: {
      heroTitle: "Connect fast. Create without limits.",
      heroSubtitle: "Enterprise multi-model API gateway",
      ctaRegister: "Get started",
      ctaLogin: "Open console",
      partnersTitle: "Partners & model ecosystem",
      feature1Title: "High-throughput gateway",
      feature1Desc: "Lock-free queues and connection pooling for peak enterprise traffic.",
      feature2Title: "Cost efficiency",
      feature2Desc: "Pay-as-you-go with transparent group multipliers.",
      feature3Title: "Unified model access",
      feature3Desc: "OpenAI / Anthropic compatible endpoints with one API key.",
    },
    about: {
      title: "About",
      intro:
        "AI API Gateway is an enterprise platform for aggregating and reselling LLM APIs with unified auth, quotas, routing, and observability.",
      missionTitle: "Mission",
      mission:
        "Help teams call leading global models securely and affordably through one OpenAI-compatible protocol.",
      featureTitle: "Capabilities",
      features: [
        "OpenAI / Anthropic compatible gateway with streaming",
        "Token groups, model allowlists, IP rules, and usage logs",
        "Wallet, redeem codes, affiliate rebates, announcements",
        "Multi-region API nodes with health checks",
      ],
      contactTitle: "Contact",
      contact: "For business inquiries, see console announcements or the home modal QR code.",
    },
    docs: {
      title: "Documentation",
      intro: "Create an API token and set the gateway Base URL before calling models.",
      openExternal: "Open full documentation",
      quickStart: "Quick start",
      authTitle: "Authentication",
      authBody: "Send Authorization: Bearer sk-xxx (create keys under Token management).",
      baseUrlTitle: "Base URL",
      baseUrlBody: "Point your SDK base_url to the gateway, e.g. https://api.example.com/v1.",
      errorsTitle: "Common errors",
      errorsBody: "402 insufficient quota; 403 model/IP restriction; 429 rate limited.",
    },
    dashboard: {
      ping: "Ping",
      pinging: "Pinging…",
      pingOk: "OK",
      pingFail: "Down",
    },
  },
};

export function t(locale: PortalLocale): MessageTree {
  return messages[locale];
}
