<script lang="ts">
  import { onMount } from "svelte";

  type ItemType = "tweet" | "paper" | "email";

  type Item = {
    name: string;
    handle: string;
    date: string;
    text: string;
    url: string;
    type: ItemType;
    accent: string;
    flag: string;
  };

  const items: Item[] = [
    {
      type: "tweet",
      name: "Tom Dörr",
      handle: "@tom_doerr",
      date: "Jan 16, 2026",
      text: "Engine for phishing detection with a web UI and browser extension",
      url: "https://twitter.com/tom_doerr/status/2012050578721915177",
      accent: "from-blue-500 to-indigo-600",
      flag: "🇩🇪",
    },
    {
      type: "paper",
      name: "Academic Research",
      handle: "ijesr.org",
      date: "Apr 26, 2025",
      text: "An earlier version of this project was referenced and reproduced (👀) in an academic research paper.",
      url: "https://www.ijesr.org/index.php/ijesr/article/view/377",
      accent: "from-amber-500 to-orange-600",
      flag: "🇮🇳",
    },
    {
      type: "tweet",
      name: "NeoTeo.com",
      handle: "@NeoteoCom",
      date: "Jan 16, 2026",
      text: "SafeSurf, an open-source phishing detection engine with a web UI and browser extension. Easy to integrate to protect users and reduce fraud.",
      url: "https://twitter.com/NeoteoCom/status/2012102861807554843",
      accent: "from-violet-500 to-purple-600",
      flag: "🇪🇸",
    },
    {
      type: "email",
      name: "Raj",
      handle: "via email",
      date: "Oct 1, 2023",
      text: "I am using this project as reference for my BTech final year project.",
      url: "https://x.com/abhizaik/status/1708401691367022776",
      accent: "from-emerald-500 to-teal-600",
      flag: "🇮🇳",
    },
  ];

  const loop = [...items, ...items];

  let track: HTMLElement;
  let offset = 0; // px scrolled so far
  let halfWidth = 0; // pixel width of one copy — reset point
  let isDragging = false;
  let isHovered = false;
  let dragStartX = 0;
  let dragStartOffset = 0;
  const SPEED = 0.5; // px per frame

  onMount(() => {
    // Measure after render so max-content width is known
    halfWidth = track.scrollWidth / 2;

    let frame: number;
    function tick() {
      if (!isDragging && !isHovered) {
        offset += SPEED;
        if (offset >= halfWidth) offset -= halfWidth;
        track.style.transform = `translateX(-${offset}px)`;
      }
      frame = requestAnimationFrame(tick);
    }
    frame = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(frame);
  });

  function onMouseEnter() {
    isHovered = true;
  }
  function onMouseLeave() {
    isHovered = false;
    isDragging = false;
  }

  function onMouseDown(e: MouseEvent) {
    isDragging = true;
    dragStartX = e.clientX;
    dragStartOffset = offset;
  }

  function onMouseUp() {
    isDragging = false;
  }

  function onMouseMove(e: MouseEvent) {
    if (!isDragging) return;
    e.preventDefault();
    offset = dragStartOffset + (dragStartX - e.clientX);
    // keep in valid range
    offset = ((offset % halfWidth) + halfWidth) % halfWidth;
    track.style.transform = `translateX(-${offset}px)`;
  }

  // Touch support
  let touchStartX = 0;
  let touchStartOffset = 0;

  function onTouchStart(e: TouchEvent) {
    touchStartX = e.touches[0].clientX;
    touchStartOffset = offset;
    isDragging = true;
  }

  function onTouchMove(e: TouchEvent) {
    if (!isDragging) return;
    offset = touchStartOffset + (touchStartX - e.touches[0].clientX);
    offset = ((offset % halfWidth) + halfWidth) % halfWidth;
    track.style.transform = `translateX(-${offset}px)`;
  }

  function onTouchEnd() {
    isDragging = false;
  }
</script>

<div class="mt-12 pb-8 md:pb-0">
  <p class="text-xs font-semibold text-gray-500 uppercase tracking-widest mb-5">
    What people are saying
  </p>

  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="overflow-x-hidden overflow-y-visible relative py-2"
    on:mouseenter={onMouseEnter}
    on:mouseleave={onMouseLeave}
    on:mousedown={onMouseDown}
    on:mouseup={onMouseUp}
    on:mousemove={onMouseMove}
    on:touchstart={onTouchStart}
    on:touchmove={onTouchMove}
    on:touchend={onTouchEnd}
    role="region"
    aria-label="Testimonials"
  >
    <!-- Fade edges -->
    <div
      class="absolute left-0 top-0 h-full w-10 bg-gradient-to-r from-gray-950 to-transparent z-10 pointer-events-none"
    ></div>
    <div
      class="absolute right-0 top-0 h-full w-10 bg-gradient-to-l from-gray-950 to-transparent z-10 pointer-events-none"
    ></div>

    <div
      bind:this={track}
      class="flex gap-4 will-change-transform {isDragging ? 'cursor-grabbing' : 'cursor-grab'}"
      style="width: max-content;"
    >
      {#each loop as item}
        <a
          href={item.url}
          target="_blank"
          rel="noopener noreferrer"
          class="group flex flex-col gap-3 p-4 rounded-xl border border-gray-800 bg-gray-900/60 hover:border-gray-700 hover:bg-gray-900 hover:-translate-y-0.5 hover:shadow-lg hover:shadow-black/30 transition-all duration-200 flex-shrink-0 w-72 text-left"
          draggable="false"
        >
          <div class="flex items-start justify-between gap-2">
            <div class="flex items-center gap-2.5">
              <div
                class="w-8 h-8 rounded-full bg-gradient-to-br {item.accent} flex items-center justify-center text-white text-xs font-bold flex-shrink-0"
              >
                {item.name[0]}
              </div>
              <div class="leading-tight">
                <p class="text-sm font-semibold text-white">{item.name}</p>
                <p class="text-xs text-gray-500">
                  {item.handle} <span class="ml-0.5">{item.flag}</span>
                </p>
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
            {:else if item.type === "email"}
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
                  d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
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

          <p class="text-sm text-gray-300 leading-relaxed relative z-10">"{item.text}"</p>
          <p class="text-xs text-gray-600">{item.date}</p>
        </a>
      {/each}
    </div>
  </div>
</div>
