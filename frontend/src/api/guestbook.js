/**
 * Guestbook API client. Uses relative /api when VITE_API_URL is unset (same-origin).
 * Backend requires Accept: application/json and Content-Type: application/json for POST.
 */

const baseURL = (import.meta.env.VITE_API_URL ?? '').replace(/\/$/, '')

const jsonHeaders = {
  Accept: 'application/json',
  'Content-Type': 'application/json',
}

/**
 * @typedef {{ id: string, name: string, message: string, createdAt: string }} GuestbookEntry
 */

/**
 * Fetch guestbook entries (newest first).
 * @returns {Promise<GuestbookEntry[]>}
 */
export async function listEntries() {
  const res = await fetch(`${baseURL}/api/guestbook/`, {
    method: 'GET',
    headers: { Accept: 'application/json' },
  })
  if (!res.ok) {
    const err = new Error(`Guestbook list failed: ${res.status}`)
    err.status = res.status
    throw err
  }
  return res.json()
}

/**
 * Create a new guestbook entry.
 * @param {{ name: string, message: string }} body
 * @returns {Promise<GuestbookEntry>}
 */
export async function createEntry({ name, message }) {
  const res = await fetch(`${baseURL}/api/guestbook/`, {
    method: 'POST',
    headers: jsonHeaders,
    body: JSON.stringify({ name, message }),
  })
  if (!res.ok) {
    let message = `Guestbook create failed: ${res.status}`
    try {
      const data = await res.json()
      if (data?.error) message = data.error
    } catch (_) {}
    const err = new Error(message)
    err.status = res.status
    throw err
  }
  return res.json()
}
