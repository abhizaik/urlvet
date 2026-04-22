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
      border: "border-emerald-500/30",
      bg: "bg-emerald-950/30",
      shadow: "shadow-emerald-500/10",
      badge: "bg-emerald-500/20 text-emerald-300 border border-emerald-500/30",
      label: "Trusted",
      ringColor: "#10b981",
      scoreText: "text-emerald-400",
    },
    Risky: {
      border: "border-red-500/30",
      bg: "bg-red-950/30",
      shadow: "shadow-red-500/10",
      badge: "bg-red-500/20 text-red-300 border border-red-500/30",
      label: "High Risk",
      ringColor: "#ef4444",
      scoreText: "text-red-400",
    },
    Suspicious: {
      border: "border-yellow-500/30",
      bg: "bg-yellow-950/30",
      shadow: "shadow-yellow-500/10",
      badge: "bg-yellow-500/20 text-yellow-300 border border-yellow-500/30",
      label: "Be Cautious",
      ringColor: "#eab308",
      scoreText: "text-yellow-400",
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
  class={`flex flex-row items-center gap-6 p-6 rounded-xl border shadow-lg ${style.border} ${style.bg} ${style.shadow}`}
>
  <!-- Verdict -->
  <div class="flex-1 flex flex-col gap-2 min-w-0">
    <span class="text-xs font-semibold text-gray-400 uppercase tracking-widest">Verdict</span>
    <div class="flex items-center gap-3 flex-wrap">
      <span class="text-3xl font-extrabold text-white tracking-tight">{verdict ?? "—"}</span>
      <span
        class={`px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-wide ${style.badge}`}
      >
        {style.label}
      </span>
    </div>
    <p class="text-[11px] text-gray-500 mt-1">{context}</p>
    {#if unreachable}
      <p class="text-[11px] text-red-400/80 mt-0.5">
        Site may be unreachable or returning no content.
      </p>
    {/if}
  </div>

  <!-- Circular Score Ring -->
  <div class="flex flex-col items-center gap-1 flex-shrink-0">
    <span class="text-xs font-semibold text-gray-400 uppercase tracking-widest mb-1"
      >Trust Score</span
    >
    <div class="relative w-20 h-20 md:w-24 md:h-24">
      <svg class="w-full h-full -rotate-90" viewBox="0 0 88 88">
        <!-- Track -->
        <circle cx="44" cy="44" r={R} fill="none" stroke="#374151" stroke-width="7" />
        <!-- Progress -->
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
        <span class={`text-2xl font-extrabold leading-none ${style.scoreText}`}
          >{finalScore !== undefined ? displayedScore : "—"}</span
        >
        <span class="block w-6 border-t border-gray-500 my-0.5"></span>
        <span class="text-[10px] text-gray-500 font-medium">100</span>
      </div>
    </div>
  </div>
</div>
