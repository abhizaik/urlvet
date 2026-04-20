<script lang="ts">
  import { onMount } from "svelte";

  type Item = {
    name: string;
    handle: string;
    date: string;
    text: string;
    url: string;
    type: "tweet" | "paper";
  };

  const items: Item[] = [
    {
      type: "tweet",
      name: "Tom Dörr",
      handle: "@tom_doerr",
      date: "Jan 16, 2026",
      text: "Engine for phishing detection with a web UI and browser extension",
      url: "https://twitter.com/tom_doerr/status/2012050578721915177",
    },
    {
      type: "tweet",
      name: "NeoTeo.com",
      handle: "@NeoteoCom",
      date: "Jan 16, 2026",
      text: "SafeSurf, an open-source phishing detection engine with a web UI and browser extension. Easy to integrate to protect users and reduce fraud.",
      url: "https://twitter.com/NeoteoCom/status/2012102861807554843",
    },
    {
      type: "paper",
      name: "Academic Research",
      handle: "ijesr.org",
      date: "2025",
      text: "An earlier version of this project was referenced and reproduced (👀) in an academic research paper.",
      url: "https://www.ijesr.org/index.php/ijesr/article/view/377",
    },
  ];

  // Duplicate for seamless loop reset
  const loop = [...items, ...items];

  let container: HTMLElement;
  let isDragging = false;
  let isHovered = false;
  let dragStartX = 0;
  let dragStartScroll = 0;
  let scrollPos = 0;
  const SPEED = 0.6; // px per frame

  onMount(() => {
    let frame: number;

    function tick() {
      if (!isDragging && !isHovered && container) {
        scrollPos += SPEED;
        // Seamless loop: reset when we've scrolled past the first copy
        const half = container.scrollWidth / 2;
        if (scrollPos >= half) scrollPos -= half;
        container.scrollLeft = scrollPos;
      }
      frame = requestAnimationFrame(tick);
    }

    frame = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(frame);
  });

  function onMouseEnter() { isHovered = true; }
  function onMouseLeave() {
    isHovered = false;
    if (isDragging) scrollPos = container.scrollLeft;
    isDragging = false;
  }

  function onMouseDown(e: MouseEvent) {
    isDragging = true;
    dragStartX = e.pageX - container.offsetLeft;
    dragStartScroll = container.scrollLeft;
  }

  function onMouseUp() {
    if (isDragging) scrollPos = container.scrollLeft;
    isDragging = false;
  }

  function onMouseMove(e: MouseEvent) {
    if (!isDragging) return;
    e.preventDefault();
    const x = e.pageX - container.offsetLeft;
    const newScroll = dragStartScroll - (x - dragStartX) * 1.5;
    container.scrollLeft = newScroll;
    scrollPos = newScroll;
  }
</script>

<div class="mt-12 pb-8 md:pb-0">
  <p class="text-xs font-semibold text-gray-500 uppercase tracking-widest mb-5">
    What people are saying
  </p>

  <div class="overflow-hidden relative">
    <!-- Fade edges -->
    <div class="absolute left-0 top-0 h-full w-10 bg-gradient-to-r from-gray-950 to-transparent z-10 pointer-events-none"></div>
    <div class="absolute right-0 top-0 h-full w-10 bg-gradient-to-l from-gray-950 to-transparent z-10 pointer-events-none"></div>

    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div
      class="overflow-x-auto scrollbar-hide select-none {isDragging ? 'cursor-grabbing' : 'cursor-grab'}"
      bind:this={container}
      on:mouseenter={onMouseEnter}
      on:mouseleave={onMouseLeave}
      on:mousedown={onMouseDown}
      on:mouseup={onMouseUp}
      on:mousemove={onMouseMove}
      role="region"
      aria-label="Testimonials"
    >
      <div class="flex gap-4" style="width: max-content;">
        {#each loop as item}
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            class="group flex flex-col gap-3 p-4 rounded-xl border border-gray-800 bg-gray-900/60 hover:border-gray-700 hover:bg-gray-900 transition-all duration-200 flex-shrink-0 w-72"
            draggable="false"
          >
            <div class="flex items-start justify-between gap-2">
              <div class="flex items-center gap-2.5">
                <div
                  class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-xs font-bold flex-shrink-0"
                >
                  {item.name[0]}
                </div>
                <div class="leading-tight">
                  <p class="text-sm font-semibold text-white">{item.name}</p>
                  <p class="text-xs text-gray-500">{item.handle}</p>
                </div>
              </div>
              {#if item.type === "tweet"}
                <svg
                  class="w-4 h-4 text-gray-600 group-hover:text-gray-400 transition-colors flex-shrink-0 mt-0.5"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.744l7.73-8.835L1.254 2.25H8.08l4.253 5.622 5.91-5.622zm-1.161 17.52h1.833L7.084 4.126H5.117z"
                  />
                </svg>
              {:else}
                <svg
                  class="w-4 h-4 text-gray-600 group-hover:text-gray-400 transition-colors flex-shrink-0 mt-0.5"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
              {/if}
            </div>

            <p class="text-sm text-gray-300 leading-relaxed">"{item.text}"</p>
            <p class="text-xs text-gray-600">{item.date}</p>
          </a>
        {/each}
      </div>
    </div>
  </div>
</div>

<style>
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
</style>
