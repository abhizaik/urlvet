<script lang="ts">
  import { onMount } from "svelte";
  import { fly } from "svelte/transition";

  type ItemType = "tweet" | "paper" | "email" | "package" | "newsletter";

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
      type: "package",
      name: "Dika Ardianta",
      handle: "@DikaArdnt",
      date: "Apr 28, 2026",
      text: "I ported the detection engine to PHP and published it as an unofficial community Composer package on packagist.org",
      url: "https://packagist.org/packages/safesurf/safesurf",
      accent: "from-fuchsia-500 to-pink-600",
      flag: "🇮🇩",
    },
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
      type: "newsletter",
      name: "OSINTech",
      handle: "osintech.substack.com",
      date: "May 21, 2026",
      text: "URLvet. Open-source phishing detection engine — get a trust score, a fully explainable verdict, and a shareable security report with live page preview, all in real time.",
      url: "https://osintech.substack.com/p/osintechs-timeline-163-21052026?open=false#%C2%A7osint-tools-services-and-investigations",
      accent: "from-green-500 to-emerald-600",
      flag: "🇰🇿",
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
      type: "tweet",
      name: "Bryan",
      handle: "@so_sthbryan",
      date: "May 16, 2026",
      text: "Explainable phishing detection that scans URLs in real time. SafeSurf catches phishing sites before you interact with them.",
      url: "https://x.com/so_sthbryan/status/2055390764339974377",
      accent: "from-teal-500 to-cyan-600",
      flag: "🇺🇸",
    },
  ];

  const loop = [...items, ...items];

  // --- Ticker (desktop) ---
  let track: HTMLElement;
  let offset = 0;
  let halfWidth = 0;
  let isDragging = false;
  let isHovered = false;
  let dragStartX = 0;
  let dragStartOffset = 0;

  // --- Slideshow (mobile) ---
  let isMobile = false;
  let currentIndex = 0;
  let direction = 1; // 1 = forward (right→left), -1 = backward (left→right)
  let touchStartX = 0;
  let touchStartOffset = 0;

  function goTo(index: number) {
    direction = index > currentIndex ? 1 : -1;
    currentIndex = index;
  }

  onMount(() => {
    isMobile = window.innerWidth < 768;

    if (!isMobile) {
      halfWidth = track.scrollWidth / 2;
    }

    let slideTimer: ReturnType<typeof setInterval>;

    function startSlideshow() {
      slideTimer = setInterval(() => {
        direction = 1;
        currentIndex = (currentIndex + 1) % items.length;
      }, 4000);
    }

    if (isMobile) startSlideshow();

    let frame: number;
    function tick() {
      if (!isMobile && !isDragging && !isHovered) {
        offset += 0.5;
        if (offset >= halfWidth) offset -= halfWidth;
        track.style.transform = `translateX(-${offset}px)`;
      }
      frame = requestAnimationFrame(tick);
    }
    frame = requestAnimationFrame(tick);

    const onResize = () => {
      const nowMobile = window.innerWidth < 768;
      if (nowMobile === isMobile) return;
      isMobile = nowMobile;
      if (isMobile) {
        startSlideshow();
      } else {
        clearInterval(slideTimer);
        setTimeout(() => {
          if (track) halfWidth = track.scrollWidth / 2;
        }, 50);
      }
    };
    window.addEventListener("resize", onResize);

    return () => {
      cancelAnimationFrame(frame);
      clearInterval(slideTimer);
      window.removeEventListener("resize", onResize);
    };
  });

  // Ticker mouse handlers
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
    offset = ((offset % halfWidth) + halfWidth) % halfWidth;
    track.style.transform = `translateX(-${offset}px)`;
  }

  // Touch handlers (swipe on mobile, drag on desktop)
  function onTouchStart(e: TouchEvent) {
    touchStartX = e.touches[0].clientX;
    if (!isMobile) {
      touchStartOffset = offset;
      isDragging = true;
    }
  }
  function onTouchMove(e: TouchEvent) {
    if (isMobile || !isDragging) return;
    offset = touchStartOffset + (touchStartX - e.touches[0].clientX);
    offset = ((offset % halfWidth) + halfWidth) % halfWidth;
    track.style.transform = `translateX(-${offset}px)`;
  }
  function onTouchEnd(e: TouchEvent) {
    isDragging = false;
    if (!isMobile) return;
    const diff = touchStartX - (e.changedTouches[0]?.clientX ?? touchStartX);
    if (Math.abs(diff) > 40) {
      if (diff > 0) {
        direction = 1;
        currentIndex = (currentIndex + 1) % items.length;
      } else {
        direction = -1;
        currentIndex = (currentIndex - 1 + items.length) % items.length;
      }
    }
  }
</script>

<div class="pb-8 md:pb-0">
  <p class="text-xs font-semibold text-gray-400 uppercase tracking-widest mb-5 text-center">
    What people are saying
  </p>

  {#if isMobile}
    <!-- Slideshow (mobile) -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div
      class="relative px-4"
      on:touchstart={onTouchStart}
      on:touchend={onTouchEnd}
      role="region"
      aria-label="Testimonials"
    >
      <div class="grid overflow-hidden">
        {#key currentIndex}
          <a
            href={items[currentIndex].url}
            target="_blank"
            rel="noopener noreferrer"
            class="col-start-1 row-start-1 group flex flex-col gap-3 p-4 rounded-xl border border-gray-300 dark:border-gray-800 bg-white dark:bg-gray-900/60 w-full text-left cursor-pointer h-[230px]"
            draggable="false"
            in:fly={{ x: direction * 300, duration: 350 }}
            out:fly={{ x: direction * -300, duration: 350 }}
          >
            <div class="flex items-start justify-between gap-2">
              <div class="flex items-center gap-2.5">
                <div
                  class="w-8 h-8 rounded-full bg-gradient-to-br {items[currentIndex]
                    .accent} flex items-center justify-center text-white text-xs font-bold flex-shrink-0"
                >
                  {items[currentIndex].name[0]}
                </div>
                <div class="leading-tight">
                  <p class="text-sm font-semibold text-gray-900 dark:text-white">
                    {items[currentIndex].name}
                  </p>
                  <p class="text-xs text-gray-500">
                    {items[currentIndex].handle}
                    <span class="ml-0.5">{items[currentIndex].flag}</span>
                  </p>
                </div>
              </div>
              {#if items[currentIndex].type === "tweet"}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 flex-shrink-0 mt-0.5"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.744l7.73-8.835L1.254 2.25H8.08l4.253 5.622 5.91-5.622zm-1.161 17.52h1.833L7.084 4.126H5.117z"
                  />
                </svg>
              {:else if items[currentIndex].type === "package"}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 flex-shrink-0 mt-0.5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3"
                  />
                </svg>
              {:else if items[currentIndex].type === "newsletter"}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 flex-shrink-0 mt-0.5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
                  />
                </svg>
              {:else}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 flex-shrink-0 mt-0.5"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25"
                  />
                </svg>
              {/if}
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-300 leading-relaxed line-clamp-4">
              "{items[currentIndex].text}"
            </p>
            <div class="flex items-center justify-between mt-auto">
              <p class="text-xs text-gray-400 dark:text-gray-600">{items[currentIndex].date}</p>
              <svg
                class="w-4 h-4 text-gray-300 dark:text-gray-700 group-hover:text-blue-500 dark:group-hover:text-blue-400 transition-colors"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                />
              </svg>
            </div>
          </a>
        {/key}
      </div>

      <!-- Dot indicators -->
      <div class="flex justify-center gap-2 mt-4">
        {#each items as _, i}
          <!-- svelte-ignore a11y_consider_explicit_label -->
          <button
            class="w-1.5 h-1.5 rounded-full transition-colors {i === currentIndex
              ? 'bg-gray-500 dark:bg-gray-400'
              : 'bg-gray-200 dark:bg-gray-700'}"
            on:click={() => goTo(i)}
            aria-label="Go to slide {i + 1}"
          ></button>
        {/each}
      </div>
    </div>
  {:else}
    <!-- Ticker (desktop) -->
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
        class="absolute left-0 top-0 h-full w-10 bg-gradient-to-r from-white dark:from-gray-950 to-transparent z-10 pointer-events-none"
      ></div>
      <div
        class="absolute right-0 top-0 h-full w-10 bg-gradient-to-l from-white dark:from-gray-950 to-transparent z-10 pointer-events-none"
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
            class="group flex flex-col gap-3 p-4 rounded-xl border border-gray-300 dark:border-gray-800 bg-white dark:bg-gray-900/60 hover:border-gray-400 dark:hover:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-900 hover:-translate-y-0.5 hover:shadow-lg hover:shadow-black/10 dark:hover:shadow-black/30 transition-all duration-200 flex-shrink-0 w-72 text-left cursor-pointer"
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
                  <p class="text-sm font-semibold text-gray-900 dark:text-white">{item.name}</p>
                  <p class="text-xs text-gray-500">
                    {item.handle} <span class="ml-0.5">{item.flag}</span>
                  </p>
                </div>
              </div>
              {#if item.type === "tweet"}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 group-hover:text-gray-600 dark:group-hover:text-gray-400 transition-colors flex-shrink-0 mt-0.5"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.744l7.73-8.835L1.254 2.25H8.08l4.253 5.622 5.91-5.622zm-1.161 17.52h1.833L7.084 4.126H5.117z"
                  />
                </svg>
              {:else if item.type === "email"}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 group-hover:text-gray-600 dark:group-hover:text-gray-400 transition-colors flex-shrink-0 mt-0.5"
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
              {:else if item.type === "package"}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 group-hover:text-gray-600 dark:group-hover:text-gray-400 transition-colors flex-shrink-0 mt-0.5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3"
                  />
                </svg>
              {:else if item.type === "newsletter"}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 group-hover:text-gray-600 dark:group-hover:text-gray-400 transition-colors flex-shrink-0 mt-0.5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
                  />
                </svg>
              {:else}
                <svg
                  class="w-4 h-4 text-gray-400 dark:text-gray-600 group-hover:text-gray-600 dark:group-hover:text-gray-400 transition-colors flex-shrink-0 mt-0.5"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25"
                  />
                </svg>
              {/if}
            </div>

            <p class="text-sm text-gray-600 dark:text-gray-300 leading-relaxed relative z-10">
              "{item.text}"
            </p>
            <div class="flex items-center justify-between mt-auto">
              <p class="text-xs text-gray-400 dark:text-gray-600">{item.date}</p>
              <svg
                class="w-4 h-4 text-gray-300 dark:text-gray-700 group-hover:text-blue-500 dark:group-hover:text-blue-400 transition-colors"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                />
              </svg>
            </div>
          </a>
        {/each}
      </div>
    </div>
  {/if}
</div>
