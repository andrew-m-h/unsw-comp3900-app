<template>
  <main class="example-page">
    <header class="page-header">
      <h1>UNSW Comp3900 App</h1>
      <p class="tagline">Vue.js frontend example</p>
    </header>

    <section class="card">
      <h2>Example page</h2>
      <p>
        This is a sample page built with Vue 3 and Vite. The build outputs
        static HTML, JS, and CSS to <code>frontend/dist/</code>.
      </p>
      <div class="actions">
        <button class="btn" @click="count++">Clicked {{ count }} times</button>
      </div>
    </section>

    <section class="card">
      <h2>Quick links</h2>
      <ul class="link-list">
        <li><a href="https://vuejs.org" target="_blank" rel="noopener">Vue.js</a></li>
        <li><a href="https://vitejs.dev" target="_blank" rel="noopener">Vite</a></li>
      </ul>
    </section>
  </main>
</template>

<script setup>
import { ref } from 'vue'

const count = ref(0)
</script>

<style scoped>
.example-page {
  max-width: 42rem;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.page-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 600;
  margin: 0 0 0.5rem;
  letter-spacing: -0.02em;
}

.tagline {
  margin: 0;
  color: var(--color-text-muted);
  font-size: 1.05rem;
}

.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.card h2 {
  font-size: 1.15rem;
  font-weight: 600;
  margin: 0 0 0.75rem;
}

.card p {
  margin: 0 0 1rem;
  color: var(--color-text-muted);
}

.card p:last-of-type {
  margin-bottom: 0;
}

.card code {
  font-size: 0.9em;
  padding: 0.15em 0.4em;
  background: var(--color-bg);
  border-radius: 6px;
  border: 1px solid var(--color-border);
}

.actions {
  margin-top: 1rem;
}

.btn {
  font: inherit;
  color: var(--color-bg);
  background: var(--color-accent);
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1.2rem;
  cursor: pointer;
  font-weight: 500;
}

.btn:hover {
  background: var(--color-accent-hover);
}

.link-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.link-list li {
  margin-bottom: 0.5rem;
}

.link-list a {
  color: var(--color-accent);
  text-decoration: none;
}

.link-list a:hover {
  text-decoration: underline;
}
</style>
