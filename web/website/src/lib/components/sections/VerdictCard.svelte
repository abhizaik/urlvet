<script lang="ts">
  export let verdict: string | undefined;
  export let finalScore: number | undefined;
  export let unreachable = false;

  let displayedScore = 0;
  let prevScore: number | undefined = undefined;

  $: if (finalScore !== undefined && finalScore !== prevScore) {
    prevScore = finalScore;
    const target = finalScore;
    const duration = 900;
    const startTime = performance.now();
    const tick = (now: number) => {
      const t = Math.min((now - startTime) / duration, 1);
      const eased = 1 - Math.pow(1 - t, 3);
      displayedScore = Math.round(eased * target);
      if (t < 1) requestAnimationFrame(tick);
    };
    requestAnimationFrame(tick);
  }

  const STYLES: Record<
    string,
    {
      border: string;
      bg: string;
      shadow: string;
      badge: string;
      label: string;
      ringColor: string;
      scoreText: string;
    }
  > = {
    Safe: {
      border: "border-emerald-300 dark:border-emerald-500/30",
      bg: "bg-emerald-50 dark:bg-emerald-950/30",
      shadow: "shadow-emerald-500/10",
      badge:
        "bg-emerald-100 dark:bg-emerald-500/20 text-emerald-700 dark:text-emerald-300 border border-emerald-300 dark:border-emerald-500/30",
      label: "Trusted",
      ringColor: "#10b981",
      scoreText: "text-emerald-600 dark:text-emerald-400",
    },
    Risky: {
      border: "border-red-300 dark:border-red-500/30",
      bg: "bg-red-50 dark:bg-red-950/30",
      shadow: "shadow-red-500/10",
      badge:
        "bg-red-100 dark:bg-red-500/20 text-red-700 dark:text-red-300 border border-red-300 dark:border-red-500/30",
      label: "High Risk",
      ringColor: "#ef4444",
      scoreText: "text-red-600 dark:text-red-400",
    },
    Suspicious: {
      border: "border-yellow-300 dark:border-yellow-500/30",
      bg: "bg-amber-50 dark:bg-yellow-950/30",
      shadow: "shadow-yellow-500/10",
      badge:
        "bg-yellow-100 dark:bg-yellow-500/20 text-yellow-700 dark:text-yellow-300 border border-yellow-300 dark:border-yellow-500/30",
      label: "Be Cautious",
      ringColor: "#eab308",
      scoreText: "text-yellow-600 dark:text-yellow-400",
    },
  };

  const R = 36;
  const CIRC = 2 * Math.PI * R;

  const CONTEXT: Record<string, string> = {
    Safe: "Typical safe sites score above 75",
    Risky: "Scores below 40 indicate high threat likelihood",
    Suspicious: "May be safe, but use caution for score 40-75",
  };

  $: style = STYLES[verdict ?? ""] ?? STYLES.Suspicious;
  $: dashOffset = CIRC - ((finalScore ?? 0) / 100) * CIRC;
  $: context = CONTEXT[verdict ?? ""] ?? CONTEXT.Suspicious;
</script>

<div
  class={`flex flex-row items-center gap-4 p-4 sm:gap-6 sm:p-6 rounded-xl border shadow-lg ${style.border} ${style.bg} ${style.shadow}`}
>
  <!-- Verdict -->
  <div class="flex-1 flex flex-col gap-1.5 min-w-0">
    <span
      class="text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest"
      >Verdict</span
    >
    <div class="flex items-center gap-2 flex-wrap">
      <span
        class="text-2xl sm:text-3xl font-extrabold text-gray-900 dark:text-white tracking-tight leading-tight"
        >{verdict ?? "—"}</span
      >
      <span
        class={`px-2.5 py-0.5 rounded-full text-[10px] sm:text-xs font-semibold uppercase tracking-wide whitespace-nowrap ${style.badge}`}
      >
        {style.label}
      </span>
    </div>
    <p class="text-[10px] sm:text-[11px] text-gray-600 dark:text-gray-500 leading-relaxed">
      {context}
    </p>
    {#if unreachable}
      <p class="text-[10px] sm:text-[11px] text-red-400/80">
        Site may be unreachable or returning no content.
      </p>
    {/if}
  </div>

  <!-- Circular Score Ring -->
  <div class="flex flex-col items-center gap-1 flex-shrink-0">
    <span
      class="text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest mb-1"
      >Trust Score</span
    >
    <div class="relative w-16 h-16 sm:w-20 sm:h-20 md:w-24 md:h-24">
      <svg class="w-full h-full -rotate-90" viewBox="0 0 88 88">
        <circle cx="44" cy="44" r={R} fill="none" stroke="var(--ring-track)" stroke-width="7" />
        <circle
          cx="44"
          cy="44"
          r={R}
          fill="none"
          stroke={style.ringColor}
          stroke-width="7"
          stroke-linecap="round"
          stroke-dasharray={CIRC}
          stroke-dashoffset={dashOffset}
          style="transition: stroke-dashoffset 0.8s ease"
        />
      </svg>
      <div class="absolute inset-0 flex flex-col items-center justify-center">
        <span class={`text-xl font-extrabold leading-none ${style.scoreText}`}
          >{finalScore !== undefined ? displayedScore : "—"}</span
        >
        <span class="block w-5 border-t border-gray-500 my-0.5"></span>
        <span class="text-[9px] sm:text-[10px] text-gray-600 dark:text-gray-500 font-medium"
          >100</span
        >
      </div>
    </div>
  </div>
</div>
