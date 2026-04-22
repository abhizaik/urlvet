<script lang="ts">
  import { tick } from "svelte";
  export let text: string;

  let iconEl: HTMLElement;
  let tooltipEl: HTMLDivElement;
  let visible = false;
  let tipLeft = 0;
  let tipTop = 0;
  let arrowLeft = "50%";

  async function show() {
    visible = true;
    await tick();
    if (!iconEl || !tooltipEl) return;

    const icon = iconEl.getBoundingClientRect();
    const tw = tooltipEl.offsetWidth;
    const th = tooltipEl.offsetHeight;
    const vw = window.innerWidth;
    const padding = 10;

    const naturalLeft = icon.left + icon.width / 2 - tw / 2;
    const clampedLeft = Math.max(padding, Math.min(vw - padding - tw, naturalLeft));
    const shift = clampedLeft - naturalLeft;

    tipLeft = clampedLeft;
    tipTop = icon.top - th - 8;

    const arrowPos = tw / 2 - shift;
    arrowLeft = `${Math.max(12, Math.min(tw - 12, arrowPos))}px`;
  }

  function hide() {
    visible = false;
  }
</script>

<div class="relative inline-flex items-center" bind:this={iconEl}>
  <div
    class="w-4 h-4 flex items-center justify-center rounded-full
           bg-gray-700 text-gray-200 text-[10px] font-bold cursor-pointer
           hover:bg-gray-600 transition-colors duration-150 select-none"
    role="button"
    tabindex="0"
    aria-label={text}
    on:mouseenter={show}
    on:mouseleave={hide}
    on:focusin={show}
    on:focusout={hide}
  >
    i
  </div>
</div>

{#if visible}
  <div
    bind:this={tooltipEl}
    role="tooltip"
    class="fixed w-max max-w-[85vw] md:max-w-xs
           bg-gray-800 text-gray-100 text-xs px-3 py-1.5 rounded-lg shadow-lg
           border border-gray-700 pointer-events-none z-[99999] whitespace-normal"
    style="left: {tipLeft}px; top: {tipTop}px;"
  >
    {text}
    <div
      class="absolute top-full w-2 h-2 bg-gray-800 border-b border-r border-gray-700"
      style="left: {arrowLeft}; transform: translateX(-50%) rotate(45deg);"
    ></div>
  </div>
{/if}
