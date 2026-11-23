/*
Added by SipToken, see SIPTOKEN.md. Licensed under the GNU Affero General
Public License v3.0, like the rest of this project.
*/
import { readCachedStatus } from '@/lib/status-query'

function field(source: unknown, key: string): unknown {
  return source !== null && typeof source === 'object'
    ? (source as Record<string, unknown>)[key]
    : undefined
}

function nonEmptyUrl(value: unknown): string | null {
  if (typeof value !== 'string') return null
  const url = value.trim().replace(/\/+$/, '')
  return url || null
}

/**
 * Base URL that API clients should call: used by the CC Switch and chat app
 * imports and by "Copy Connection Info".
 *
 * Upstream uses the server address, which is the console's own URL. When the
 * admin has listed API addresses (API info), the first one is where clients
 * should go, so it wins. Accepts the status data or the full /api/status
 * response; falls back to the cached status when none is given.
 */
export function getApiBaseAddress(status?: unknown): string {
  const raw = status ?? readCachedStatus()
  const data = field(raw, 'data') ?? raw

  if (field(data, 'api_info_enabled') !== false) {
    const apiInfo = field(data, 'api_info')
    if (Array.isArray(apiInfo)) {
      for (const item of apiInfo) {
        const url = nonEmptyUrl(field(item, 'url'))
        if (url) return url
      }
    }
  }

  const serverAddress =
    nonEmptyUrl(field(data, 'server_address')) ??
    nonEmptyUrl(field(data, 'serverAddress'))
  if (serverAddress) return serverAddress

  return typeof window !== 'undefined' ? window.location.origin : ''
}
