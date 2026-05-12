<script lang="ts">
  import { theme } from "$lib/theme";
</script>

<div class="fixed top-4 right-4 z-50">
  <input
    id="theme-toggle"
    type="checkbox"
    class="toggle"
    role="switch"
    aria-label="Toggle dark mode"
    checked={$theme === "dark"}
    on:change={theme.toggle}
  />
  <label for="theme-toggle" class="sr-only">Dark mode</label>
</div>

<style>
  .toggle {
    --dur: 0.3s;
    font-size: 20px;
    appearance: none;
    -webkit-appearance: none;
    -webkit-tap-highlight-color: transparent;
    background: currentColor;
    border-radius: 0.75em;
    color: #4b5563;
    cursor: pointer;
    display: block;
    height: 1.5em;
    outline: none;
    overflow: hidden;
    position: relative;
    width: 3em;
  }

  :global(.dark) .toggle {
    color: #29254d;
  }

  .toggle:before,
  .toggle:after {
    background: #f3f4f6;
    border-radius: 50%;
    content: "";
    display: block;
    height: 1em;
    left: 0.25em;
    position: absolute;
    top: 0.25em;
    width: 1em;
    transition: transform var(--dur);
  }

  .toggle:before,
  .toggle:checked:after {
    transition-delay: 0s;
    transition-timing-function: ease-out;
  }

  .toggle:checked:before,
  .toggle:after {
    transition-delay: calc(var(--dur) * 0.75);
    transition-timing-function: cubic-bezier(0.3, 1.6, 0.5, 0.7);
  }

  /* unchecked: before is hidden left, after is visible */
  .toggle:before {
    transform: translateX(-1.5em);
  }

  /* checked: before bounces in from left */
  .toggle:checked:before {
    transform: translateX(1.5em);
  }

  /* checked: after exits right */
  .toggle:checked:after {
    transform: translateX(3em);
  }
</style>
