import { getGatewayBaseUrl } from "@/lib/gateway-config";

import type { PlaygroundChatRequest } from "./playground-types";

export type PlaygroundStreamCallbacks = {
  onDelta: (text: string) => void;
  onDone: () => void;
  onError: (message: string) => void;
};

function parseSSELine(line: string): string | null {
  const trimmed = line.trim();
  if (!trimmed.startsWith("data:")) {
    return null;
  }
  const data = trimmed.slice(5).trim();
  if (data === "[DONE]") {
    return "";
  }
  try {
    const json = JSON.parse(data) as {
      choices?: { delta?: { content?: string } }[];
      error?: { message?: string };
    };
    if (json.error?.message) {
      throw new Error(json.error.message);
    }
    const piece = json.choices?.[0]?.delta?.content;
    return piece ?? null;
  } catch (e) {
    if (e instanceof Error && e.message !== "Unexpected end of JSON input") {
      throw e;
    }
    return null;
  }
}

export async function streamPlaygroundChat(
  token: string,
  body: PlaygroundChatRequest,
  callbacks: PlaygroundStreamCallbacks,
  signal?: AbortSignal,
): Promise<void> {
  const base = getGatewayBaseUrl().replace(/\/$/, "");
  const res = await fetch(`${base}/api/playground/chat`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ ...body, stream: true }),
    signal,
  });

  if (!res.ok) {
    const text = await res.text();
    let message = res.statusText;
    try {
      const json = JSON.parse(text) as { error?: string };
      if (json.error) message = json.error;
    } catch {
      if (text) message = text;
    }
    callbacks.onError(message);
    return;
  }

  const contentType = res.headers.get("Content-Type") ?? "";
  if (!contentType.includes("text/event-stream") && !res.body) {
    const text = await res.text();
    try {
      const json = JSON.parse(text) as {
        choices?: { message?: { content?: string } }[];
      };
      const content = json.choices?.[0]?.message?.content ?? "";
      if (content) callbacks.onDelta(content);
    } catch {
      callbacks.onError(text || "unexpected response");
      return;
    }
    callbacks.onDone();
    return;
  }

  const reader = res.body?.getReader();
  if (!reader) {
    callbacks.onError("no response body");
    return;
  }

  const decoder = new TextDecoder();
  let buffer = "";

  try {
    for (;;) {
      const { done, value } = await reader.read();
      if (done) break;
      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split("\n");
      buffer = lines.pop() ?? "";
      for (const line of lines) {
        if (line.trim() === "") continue;
        const piece = parseSSELine(line);
        if (piece === "") {
          callbacks.onDone();
          return;
        }
        if (piece) callbacks.onDelta(piece);
      }
    }
    callbacks.onDone();
  } catch (e) {
    if (signal?.aborted) return;
    callbacks.onError(e instanceof Error ? e.message : "stream failed");
  }
}
