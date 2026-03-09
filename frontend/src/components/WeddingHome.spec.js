/**
 * Unit / component tests for WeddingHome.
 * Guestbook API is mocked so no real HTTP requests are made.
 */
import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import WeddingHome from './WeddingHome.vue'
import * as guestbook from '../api/guestbook.js'

vi.mock('../api/guestbook.js', () => ({
  listEntries: vi.fn(),
  createEntry: vi.fn(),
  deleteEntry: vi.fn(),
}))

/** Wait for async work (e.g. loadMessages) and Vue to re-render */
async function flush() {
  await new Promise((r) => setTimeout(r, 0))
  await new Promise((r) => setTimeout(r, 0))
}

describe('WeddingHome', () => {
  beforeEach(() => {
    vi.mocked(guestbook.listEntries).mockResolvedValue([])
  })

  it('renders the guestbook section and form', async () => {
    const wrapper = mount(WeddingHome)
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.message-section').exists()).toBe(true)
    expect(wrapper.find('.guestbook-form').exists()).toBe(true)
    expect(wrapper.find('#guestbook-name').exists()).toBe(true)
    expect(wrapper.find('#guestbook-message').exists()).toBe(true)
  })

  it('loads and displays entries on mount', async () => {
    const entries = [
      {
        id: 'e1',
        name: 'Alice',
        message: 'Hello!',
        createdAt: '2025-01-01T10:00:00Z',
      },
    ]
    vi.mocked(guestbook.listEntries).mockResolvedValue(entries)

    const wrapper = mount(WeddingHome)
    await flush()
    await wrapper.vm.$nextTick()

    expect(wrapper.findAll('.guestbook-entry')).toHaveLength(1)
    expect(wrapper.find('.guestbook-entry-name').text()).toBe('Alice')
    expect(wrapper.find('.guestbook-entry-message').text()).toBe('Hello!')
  })

  it('removes an entry from the DOM when the delete button is clicked', async () => {
    const entries = [
      { id: 'a', name: 'Alice', message: 'A', createdAt: '2025-01-01T10:00:00Z' },
      { id: 'b', name: 'Bob', message: 'B', createdAt: '2025-01-01T11:00:00Z' },
    ]
    vi.mocked(guestbook.listEntries).mockResolvedValue(entries)

    const wrapper = mount(WeddingHome)
    await flush()
    await wrapper.vm.$nextTick()
    expect(wrapper.findAll('.guestbook-entry')).toHaveLength(2)

    const deleteButtons = wrapper.findAll('.guestbook-entry-delete')
    await deleteButtons[0].trigger('click')
    await wrapper.vm.$nextTick()

    const items = wrapper.findAll('.guestbook-entry')
    expect(items).toHaveLength(1)
    expect(items[0].find('.guestbook-entry-name').text()).toBe('Bob')
  })

  it('shows empty state when there are no entries', async () => {
    const wrapper = mount(WeddingHome)
    await flush()
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.guestbook-empty').exists()).toBe(true)
    expect(wrapper.find('.guestbook-empty').text()).toContain('No messages yet')
  })
})
