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
      <p class="placeholder">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut enim ad minim veniam.</p>
    </section>

    <footer class="footer">
      <p>© 2025</p>
    </footer>
  </main>
</template>

<script setup>
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
