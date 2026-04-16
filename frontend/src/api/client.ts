// 최소 fetch 래퍼 — /api/* 경로는 Vite dev proxy 가 Go 백엔드(:8080)로 프록시.

export type ApiError = {
  status: number;
  code?: string;
  message?: string;
};

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(path, {
    method,
    headers: body !== undefined ? { 'Content-Type': 'application/json' } : undefined,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });
  if (!res.ok) {
    const err: ApiError = { status: res.status };
    try {
      const parsed = (await res.json()) as { error?: { code?: string; message?: string } };
      err.code = parsed.error?.code;
      err.message = parsed.error?.message;
    } catch {
      /* ignore */
    }
    throw err;
  }
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}

export const api = {
  get: <T>(p: string) => request<T>('GET', p),
  post: <T>(p: string, b?: unknown) => request<T>('POST', p, b),
  patch: <T>(p: string, b?: unknown) => request<T>('PATCH', p, b),
  del: <T>(p: string) => request<T>('DELETE', p),
};
