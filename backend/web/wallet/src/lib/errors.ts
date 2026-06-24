import axios from 'axios';

// Backend error envelope. transaction-service sends { error: <code>, message: <human> };
// gateway/echo errors send only { message }. We prefer the human-readable message.
interface ApiErrorBody {
  error?: string;
  message?: string;
}

// isHttpError is true when the request reached the server and came back with a
// status (4xx/5xx). It is false for network failures, timeouts or aborts where
// no response exists — letting callers tell "server said no" from "request
// never completed".
export function isHttpError(err: unknown): boolean {
  return axios.isAxiosError(err) && err.response !== undefined;
}

// extractApiError pulls a user-facing message out of an unknown error, falling
// back to the caller-provided default for non-axios / bodyless errors.
export function extractApiError(err: unknown, fallback: string): string {
  if (axios.isAxiosError(err)) {
    const body = err.response?.data as ApiErrorBody | undefined;
    return body?.message ?? body?.error ?? fallback;
  }
  return fallback;
}

// extractApiErrorCode returns the machine error code (the `error` field) for
// callers that branch on it — e.g. ErrorPage maps codes like NO_INIT_DATA /
// SESSION_EXPIRED to copy. Prefers the code over the human message.
export function extractApiErrorCode(err: unknown, fallback: string): string {
  if (axios.isAxiosError(err)) {
    const body = err.response?.data as ApiErrorBody | undefined;
    return body?.error ?? fallback;
  }
  return fallback;
}
