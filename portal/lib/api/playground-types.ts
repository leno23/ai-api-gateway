export type PlaygroundRole = "user" | "assistant" | "system";

export type PlaygroundMessage = {
  id: string;
  role: PlaygroundRole;
  content: string;
  debug?: unknown;
};

export type PlaygroundConfig = {
  version: 1;
  token_group: string;
  model: string;
  custom_body: boolean;
  custom_body_raw: string;
  image_url?: string;
  max_tokens: number;
  temperature: number;
  temperature_enabled: boolean;
  top_p: number;
  top_p_enabled: boolean;
  frequency_penalty: number;
  frequency_penalty_enabled: boolean;
  presence_penalty: number;
  presence_penalty_enabled: boolean;
};

export type PlaygroundChatRequest = {
  model: string;
  messages: { role: PlaygroundRole; content: string }[];
  stream?: boolean;
  max_tokens?: number;
  temperature?: number;
  top_p?: number;
  frequency_penalty?: number;
  presence_penalty?: number;
  api_key_id?: number;
  token_group?: string;
  custom_body?: boolean;
  custom_body_raw?: unknown;
};
