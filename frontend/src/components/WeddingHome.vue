<template>
  <main class="wedding-home">
    <header class="hero">
      <h1 class="couple-names">Emma & James</h1>
      <p class="date">15.06.2025</p>
      <p class="save-the-date">Save the date</p>
      <blockquote class="quote">
        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod
        tempor incididunt ut labore et dolore magna aliqua.
      </blockquote>
    </header>

    <div class="divider" aria-hidden="true">✦</div>

    <section class="section couple-section">
      <div class="half">
        <h2>The Bride</h2>
        <p class="name">Emma</p>
        <p class="title">Emma Louise Clarke</p>
        <p class="meta">Lorem ipsum dolor sit amet</p>
        <p class="parents">Daughter of Mr. John Clarke & Mrs. Mary Clarke</p>
      </div>
      <div class="half">
        <h2>The Groom</h2>
        <p class="name">James</p>
        <p class="title">James Robert Wright</p>
        <p class="meta">Consectetur adipiscing elit</p>
        <p class="parents">Son of Mr. Thomas Wright & Mrs. Sarah Wright</p>
      </div>
    </section>

    <div class="divider" aria-hidden="true">✦</div>

    <section class="section event-details">
      <h2>Event details</h2>
      <div class="event-block">
        <h3>Wedding ceremony</h3>
        <p>Saturday, 15 June 2025</p>
        <p>2:00 PM – 4:00 PM</p>
      </div>
      <div class="event-block">
        <h3>Reception</h3>
        <p>Saturday, 15 June 2025</p>
        <p>5:00 PM onwards</p>
      </div>
    </section>

    <div class="divider" aria-hidden="true">✦</div>

    <section class="section message-section">
      <h2>Leave us a message</h2>
      <p class="intro">We’d love to hear from you.</p>

      <form class="guestbook-form" @submit.prevent="onSubmit">
        <label class="form-label" for="guestbook-name">Your name</label>
        <input
          id="guestbook-name"
          v-model="form.name"
          type="text"
          class="form-input"
          placeholder="Your name"
          required
        />
        <label class="form-label" for="guestbook-message">Message</label>
        <textarea
          id="guestbook-message"
          v-model="form.message"
          class="form-textarea"
          placeholder="Your message..."
          rows="3"
          required
        />
        <p v-if="submitError" class="form-error">{{ submitError }}</p>
        <button type="submit" class="form-submit" :disabled="submitting">
          {{ submitting ? 'Sending…' : 'Send message' }}
        </button>
      </form>

      <div v-if="loading" class="guestbook-loading">Loading messages…</div>
      <div v-else-if="loadError" class="guestbook-error">{{ loadError }}</div>
      <ul v-else-if="entries.length" class="guestbook-list">
        <li v-for="entry in entries" :key="entry.id" class="guestbook-entry">
          <button
            type="button"
            class="guestbook-entry-delete"
            aria-label="Remove comment"
            @click="removeEntry(entry.id)"
          >
            ×
          </button>
          <p class="guestbook-entry-name">{{ entry.name }}</p>
          <p class="guestbook-entry-message">{{ entry.message }}</p>
          <time class="guestbook-entry-date" :datetime="entry.createdAt">{{ formatDate(entry.createdAt) }}</time>
        </li>
      </ul>
      <p v-else class="guestbook-empty">No messages yet. Be the first to leave one!</p>
    </section>

    <footer class="footer">
      <p>© 2025</p>
    </footer>
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listEntries, createEntry } from '../api/guestbook.js'

const entries = ref([])
const loading = ref(true)
const loadError = ref('')
const form = ref({ name: '', message: '' })
const submitting = ref(false)
const submitError = ref('')

function formatDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleDateString(undefined, { dateStyle: 'medium' }) + ' at ' + d.toLocaleTimeString(undefined, { timeStyle: 'short' })
}

async function loadMessages() {
  loading.value = true
  loadError.value = ''
  try {
    entries.value = await listEntries()
  } catch (e) {
    loadError.value = e.message || 'Failed to load messages.'
  } finally {
    loading.value = false
  }
}

function removeEntry(id) {
  entries.value = entries.value.filter((e) => e.id !== id)
}

async function onSubmit() {
  const name = (form.value.name || '').trim()
  const message = (form.value.message || '').trim()
  if (!name || !message) {
    submitError.value = 'Please enter your name and message.'
    return
  }
  submitError.value = ''
  submitting.value = true
  try {
    const created = await createEntry({ name, message })
    entries.value = [created, ...entries.value]
    form.value = { name: '', message: '' }
  } catch (e) {
    submitError.value = e.message || 'Failed to send message.'
  } finally {
    submitting.value = false
  }
}

onMounted(loadMessages)
</script>

<style scoped>
.wedding-home {
  max-width: 42rem;
  margin: 0 auto;
  padding: 2rem 1.5rem 3rem;
  text-align: center;
}

.hero {
  padding: 2rem 0;
}

.couple-names {
  font-family: var(--font-display);
  font-size: clamp(2rem, 5vw, 2.75rem);
  font-weight: 400;
  letter-spacing: 0.02em;
  margin: 0 0 0.5rem;
}

.date {
  font-size: 1.25rem;
  color: var(--color-muted);
  margin: 0 0 1rem;
}

.save-the-date {
  font-size: 0.95rem;
  text-transform: uppercase;
  letter-spacing: 0.2em;
  margin: 0 0 2rem;
  color: var(--color-muted);
}

.quote {
  font-style: italic;
  color: var(--color-muted);
  margin: 0;
  padding: 0 1rem;
  font-size: 0.95rem;
  line-height: 1.7;
}

.divider {
  color: var(--color-accent);
  font-size: 1rem;
  margin: 2rem 0;
  opacity: 0.7;
}

.section {
  margin: 2rem 0;
}

.section h2 {
  font-family: var(--font-display);
  font-size: 1.35rem;
  font-weight: 400;
  margin: 0 0 1.5rem;
  letter-spacing: 0.02em;
}

.couple-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  text-align: center;
}

@media (max-width: 520px) {
  .couple-section {
    grid-template-columns: 1fr;
  }
}

.half h2 {
  font-size: 1.1rem;
  margin-bottom: 1rem;
}

.half .name {
  font-family: var(--font-display);
  font-size: 1.35rem;
  margin: 0 0 0.25rem;
}

.half .title {
  font-size: 0.95rem;
  color: var(--color-text);
  margin: 0 0 0.5rem;
}

.half .meta {
  font-size: 0.85rem;
  color: var(--color-muted);
  margin: 0 0 0.5rem;
}

.half .parents {
  font-size: 0.8rem;
  color: var(--color-muted);
  margin: 0;
}

.event-details .event-block {
  margin-bottom: 1.5rem;
}

.event-details .event-block:last-child {
  margin-bottom: 0;
}

.event-block h3 {
  font-size: 1rem;
  font-weight: 600;
  margin: 0 0 0.5rem;
}

.event-block p {
  margin: 0;
  font-size: 0.95rem;
  color: var(--color-muted);
}

.message-section .intro {
  margin: 0 0 0.5rem;
  font-size: 1rem;
}

.message-section .placeholder {
  margin: 0;
  font-size: 0.9rem;
  color: var(--color-muted);
}

.guestbook-form {
  text-align: left;
  max-width: 24rem;
  margin: 0 auto 2rem;
}

.form-label {
  display: block;
  font-size: 0.9rem;
  font-weight: 500;
  margin: 0 0 0.35rem;
  color: var(--color-text);
}

.form-input,
.form-textarea {
  width: 100%;
  font: inherit;
  font-size: 0.95rem;
  padding: 0.5rem 0.6rem;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  margin-bottom: 1rem;
  background: var(--color-surface);
  color: var(--color-text);
}

.form-textarea {
  resize: vertical;
  min-height: 4rem;
}

.form-error {
  margin: -0.5rem 0 0.75rem;
  font-size: 0.9rem;
  color: #c53030;
}

.form-submit {
  font: inherit;
  color: var(--color-bg);
  background: var(--color-accent);
  border: none;
  border-radius: 8px;
  padding: 0.5rem 1.25rem;
  cursor: pointer;
  font-weight: 500;
}

.form-submit:hover:not(:disabled) {
  background: var(--color-accent-hover);
}

.form-submit:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.guestbook-loading,
.guestbook-error,
.guestbook-empty {
  font-size: 0.95rem;
  color: var(--color-muted);
  margin: 0;
}

.guestbook-error {
  color: #c53030;
}

.guestbook-list {
  list-style: none;
  margin: 0;
  padding: 0;
  text-align: left;
}

.guestbook-entry {
  position: relative;
  padding: 1rem 1.75rem 1rem 0;
  border-bottom: 1px solid var(--color-border);
}

.guestbook-entry-delete {
  position: absolute;
  top: 0.5rem;
  right: 0;
  width: 1.5rem;
  height: 1.5rem;
  padding: 0;
  font-size: 1.25rem;
  line-height: 1;
  color: var(--color-muted);
  background: none;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.guestbook-entry-delete:hover {
  color: var(--color-text);
  background: var(--color-border);
}

.guestbook-entry:last-child {
  border-bottom: none;
}

.guestbook-entry-name {
  font-weight: 600;
  font-size: 0.95rem;
  margin: 0 0 0.25rem;
  color: var(--color-text);
}

.guestbook-entry-message {
  margin: 0 0 0.35rem;
  font-size: 0.95rem;
  color: var(--color-text);
  white-space: pre-wrap;
}

.guestbook-entry-date {
  font-size: 0.8rem;
  color: var(--color-muted);
}

.footer {
  margin-top: 3rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--color-border);
  font-size: 0.8rem;
  color: var(--color-muted);
}

.footer p {
  margin: 0;
}
</style>
