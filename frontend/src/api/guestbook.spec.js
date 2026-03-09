/**
 * Unit tests for the guestbook API client.
 * fetch is mocked so no real HTTP requests are made.
 */
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { listEntries, createEntry } from './guestbook.js'

describe('guestbook API', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  describe('listEntries', () => {
    it('requests GET /api/guestbook/ and returns parsed JSON', async () => {
      const data = [
        { id: '1', name: 'Alice', message: 'Hi', createdAt: '2025-01-01T12:00:00Z' },
      ]
      fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(data),
      })

      const result = await listEntries()

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/guestbook/'),
        { method: 'GET', headers: { Accept: 'application/json' } }
      )
      expect(result).toEqual(data)
    })

    it('throws with status when response is not ok', async () => {
      fetch.mockResolvedValueOnce({ ok: false, status: 500 })

      await expect(listEntries()).rejects.toThrow('Guestbook list failed: 500')
    })
  })

  describe('createEntry', () => {
    it('sends POST with JSON body and returns created entry', async () => {
      const body = { name: 'Bob', message: 'Hello' }
      const created = { id: '2', ...body, createdAt: '2025-01-02T12:00:00Z' }
      fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(created),
      })

      const result = await createEntry(body)

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/guestbook/'),
        {
          method: 'POST',
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(body),
        }
      )
      expect(result).toEqual(created)
    })

    it('throws with server error message when response has data.error', async () => {
      fetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: () => Promise.resolve({ error: 'Name is required' }),
      })

      await expect(createEntry({ name: '', message: 'x' })).rejects.toThrow(
        'Name is required'
      )
    })

    it('throws generic message when response is not ok and no data.error', async () => {
      fetch.mockResolvedValueOnce({
        ok: false,
        status: 503,
        json: () => Promise.reject(new Error('not json')),
      })

      await expect(createEntry({ name: 'A', message: 'B' })).rejects.toThrow(
        'Guestbook create failed: 503'
      )
    })
  })
})
