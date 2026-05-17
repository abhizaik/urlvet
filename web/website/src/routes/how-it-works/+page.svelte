<svelte:head>
  <title>How It Works — url.vet</title>
  <meta
    name="description"
    content="How url.vet detects phishing links: 18 concurrent analyzers, 33 signals across 7 categories, fully explainable scoring. No black-box ML."
  />
  <link rel="canonical" href="https://url.vet/how-it-works" />
</svelte:head>

<div class="min-h-screen bg-white dark:bg-gray-950 text-gray-900 dark:text-gray-100">
  <div class="max-w-3xl mx-auto px-6 py-16">
    <!-- Back link -->
    <a
      href="/"
      class="inline-flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors mb-10"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
      </svg>
      url.vet
    </a>

    <h1 class="text-3xl font-bold tracking-tight mb-3">How It Works</h1>
    <p class="text-gray-500 dark:text-gray-400 text-lg mb-8">
      18 concurrent analyzers. 33 signals. Every verdict is fully explainable, no black-box ML.
    </p>

    <div class="space-y-4 text-gray-600 dark:text-gray-400 text-sm leading-relaxed mb-12">
      <p>
        When you paste a URL, url.vet doesn't look it up in a database. It actively fetches and
        analyzes it in real time. Eighteen checks run in parallel, each probing a different
        dimension of the link: its structure, DNS records, TLS certificate, domain age, whether the
        page content matches what the domain implies, and whether it appears in any threat feeds.
      </p>
      <p>
        Each check emits a reason (good, bad, or neutral) and a numeric weight. Those weights feed
        into a single score using a transparent formula. You can see exactly what pushed the score
        up or down, which is the only way a safety verdict is actually useful.
      </p>
    </div>

    <!-- Pipeline image -->
    <div class="mb-14">
      <img
        src="/pipeline.png"
        alt="url.vet analyzer pipeline showing 18 checks across 7 signal categories feeding into a single trust score"
        class="w-full rounded-xl border border-gray-200 dark:border-gray-800 dark:invert"
        loading="lazy"
      />
      <p class="text-xs text-gray-400 dark:text-gray-500 mt-2 text-center">
        All 18 checks run concurrently via goroutines
      </p>
    </div>

    <!-- Scoring formula -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-4">The Score</h2>
      <div
        class="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-5 font-mono text-sm mb-5"
      >
        finalScore = clamp(50 + (trustScore − riskScore) × 0.5)
      </div>
      <p class="text-gray-600 dark:text-gray-400 text-sm leading-relaxed mb-5">
        50 is the neutral baseline. A URL with no signals at all scores exactly 50, landing in
        Suspicious. That's intentional: an unknown link should be treated as unknown, not safe.
        Trust signals pull the score up; risk signals pull it down. Each side is capped at 100
        individually, so a single catastrophic signal can't drown out everything else.
      </p>
      <div class="flex flex-wrap gap-3">
        <span
          class="px-3 py-1.5 rounded-full text-sm font-medium bg-red-100 dark:bg-red-900/40 text-red-700 dark:text-red-300"
        >
          Risky — score &lt; 30
        </span>
        <span
          class="px-3 py-1.5 rounded-full text-sm font-medium bg-yellow-100 dark:bg-yellow-900/40 text-yellow-700 dark:text-yellow-300"
        >
          Suspicious — 30 to 64
        </span>
        <span
          class="px-3 py-1.5 rounded-full text-sm font-medium bg-green-100 dark:bg-green-900/40 text-green-700 dark:text-green-300"
        >
          Safe — 65 and above
        </span>
      </div>
    </section>

    <!-- Signal categories -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-2">Signal Categories</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400 mb-8">
        Seven categories, 33 individual signals. Each runs independently; a failure in one never
        blocks another.
      </p>
      <div class="space-y-8">
        <div>
          <div class="flex items-center gap-3 mb-3">
            <span
              class="text-xs font-semibold uppercase tracking-wider text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/30 px-2.5 py-1 rounded-full"
              >URL Structure</span
            >
            <span class="text-xs text-gray-400">8 signals</span>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
            Analyzed before any network request is made.
          </p>
          <ol class="space-y-1.5 text-sm text-gray-700 dark:text-gray-300 list-decimal list-inside">
            <li>Raw IP address used as hostname instead of a domain name</li>
            <li>Punycode or IDN-encoded hostname (lookalike domain spoofing)</li>
            <li>URL shortener service detected</li>
            <li>Abnormally long URL</li>
            <li>Deeply nested path structure</li>
            <li>Phishing keywords in the path: login, verify, secure, update, confirm…</li>
            <li>Excessive subdomain depth</li>
            <li>
              Non-ASCII characters in hostname (IDN homograph attack, e.g. аpple.com with Cyrillic
              а)
            </li>
          </ol>
        </div>

        <div>
          <div class="flex items-center gap-3 mb-3">
            <span
              class="text-xs font-semibold uppercase tracking-wider text-purple-600 dark:text-purple-400 bg-purple-50 dark:bg-purple-900/30 px-2.5 py-1 rounded-full"
              >HTTP / Network</span
            >
            <span class="text-xs text-gray-400">4 signals</span>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
            One real HTTP request follows the full redirect chain.
          </p>
          <ol
            start="9"
            class="space-y-1.5 text-sm text-gray-700 dark:text-gray-300 list-decimal list-inside"
          >
            <li>Redirect chain length</li>
            <li>Final destination is a different domain than the original link</li>
            <li>HSTS header present</li>
            <li>HTTP status code</li>
          </ol>
        </div>

        <div>
          <div class="flex items-center gap-3 mb-3">
            <span
              class="text-xs font-semibold uppercase tracking-wider text-orange-600 dark:text-orange-400 bg-orange-50 dark:bg-orange-900/30 px-2.5 py-1 rounded-full"
              >DNS</span
            >
            <span class="text-xs text-gray-400">3 signals</span>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
            Checks that the domain resolves and has legitimate infrastructure in place.
          </p>
          <ol
            start="13"
            class="space-y-1.5 text-sm text-gray-700 dark:text-gray-300 list-decimal list-inside"
          >
            <li>NS record validity</li>
            <li>MX record validity</li>
            <li>IP resolution</li>
          </ol>
        </div>

        <div>
          <div class="flex items-center gap-3 mb-3">
            <span
              class="text-xs font-semibold uppercase tracking-wider text-green-600 dark:text-green-400 bg-green-50 dark:bg-green-900/30 px-2.5 py-1 rounded-full"
              >TLS / SSL</span
            >
            <span class="text-xs text-gray-400">2 signals</span>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
            One TLS handshake checks the full certificate chain.
          </p>
          <ol
            start="16"
            class="space-y-1.5 text-sm text-gray-700 dark:text-gray-300 list-decimal list-inside"
          >
            <li>TLS presence and hostname match</li>
            <li>
              Certificate chain validity: expiry, issuer, CT log inclusion, known-bad fingerprints
            </li>
          </ol>
        </div>

        <div>
          <div class="flex items-center gap-3 mb-3">
            <span
              class="text-xs font-semibold uppercase tracking-wider text-cyan-600 dark:text-cyan-400 bg-cyan-50 dark:bg-cyan-900/30 px-2.5 py-1 rounded-full"
              >Domain Intelligence</span
            >
            <span class="text-xs text-gray-400">6 signals</span>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
            Passive checks on the domain's history and registration.
          </p>
          <ol
            start="18"
            class="space-y-1.5 text-sm text-gray-700 dark:text-gray-300 list-decimal list-inside"
          >
            <li>Global traffic rank (position in the top-1M popularity list)</li>
            <li>TLD classification: trusted, risky, or ICANN-registered</li>
            <li>Domain registration age via WHOIS (newly registered domains carry high risk)</li>
            <li>DNSSEC</li>
            <li>Shannon entropy score (detects algorithmically generated domain names)</li>
            <li>Typosquatting and combo-squatting against 500+ known brands</li>
          </ol>
        </div>

        <div>
          <div class="flex items-center gap-3 mb-3">
            <span
              class="text-xs font-semibold uppercase tracking-wider text-rose-600 dark:text-rose-400 bg-rose-50 dark:bg-rose-900/30 px-2.5 py-1 rounded-full"
              >Content Analysis</span
            >
            <span class="text-xs text-gray-400">8 signals</span>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
            The page is fetched and parsed. Most expensive step, most revealing.
          </p>
          <ol
            start="24"
            class="space-y-1.5 text-sm text-gray-700 dark:text-gray-300 list-decimal list-inside"
          >
            <li>Login form on an unranked or newly registered domain</li>
            <li>Payment form with credit card or CVV fields</li>
            <li>Personal information collection form</li>
            <li>Hidden iframes (credential theft or clickjacking vector)</li>
            <li>1×1 tracking pixels</li>
            <li>
              Brand impersonation: brand name in page content doesn't match the hosting domain
            </li>
            <li>Form that submits data to a third-party domain</li>
            <li>Password field served over unencrypted HTTP</li>
          </ol>
        </div>

        <div>
          <div class="flex items-center gap-3 mb-3">
            <span
              class="text-xs font-semibold uppercase tracking-wider text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/30 px-2.5 py-1 rounded-full"
              >Threat Intelligence</span
            >
            <span class="text-xs text-gray-400">2 signals</span>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
            Cross-referenced against community-maintained phishing databases.
          </p>
          <ol
            start="32"
            class="space-y-1.5 text-sm text-gray-700 dark:text-gray-300 list-decimal list-inside"
          >
            <li>PhishTank confirmed phishing (community-verified)</li>
            <li>PhishTank reported but unverified (cached 3 hours)</li>
          </ol>
        </div>
      </div>
    </section>

    <!-- Request lifecycle -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-2">What Happens When You Submit a URL</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400 mb-6">
        From paste to verdict, typically in under 5 seconds.
      </p>
      <ol class="space-y-4">
        {#each [{ step: "Validate and normalize", desc: "The URL is validated and normalized. If the scheme is missing, https:// is inferred." }, { step: "Cache check", desc: "Results are cached for 24 hours. A cache hit returns the full report instantly with no re-analysis." }, { step: "18 goroutines launch", desc: "On a cache miss, all 18 checks start simultaneously via sync.WaitGroup. A panic in one task is recovered without affecting the others." }, { step: "Score aggregated", desc: "Each check returns a numeric weight and a reason string. Trust and risk scores are summed, clamped, and fed into the formula." }, { step: "Result returned", desc: "Trust score, verdict, every individual reason, redirect chain, page screenshot, and per-task timing in one response." }] as item, i}
          <li class="flex gap-4">
            <span
              class="flex-shrink-0 w-6 h-6 rounded-full bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 text-xs font-semibold flex items-center justify-center mt-0.5"
              >{i + 1}</span
            >
            <div>
              <p class="text-sm font-medium">{item.step}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{item.desc}</p>
            </div>
          </li>
        {/each}
      </ol>
    </section>

    <!-- Limitations -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-4">Limitations</h2>
      <div class="space-y-3 text-sm text-gray-600 dark:text-gray-400 leading-relaxed">
        <p>
          Heuristic detection means false positives are possible. A legitimate site that's new and
          unranked might score lower than it deserves. No tool is a guarantee.
        </p>
        <p>
          There's no ML model, and that's intentional. A model that can't explain itself isn't
          useful when you need to make a trust decision in the moment. Every signal here can be read
          and reasoned about.
        </p>
        <p>Use url.vet as one layer of defense, not the only one.</p>
      </div>
    </section>

    <!-- CTA -->
    <div
      class="border-t border-gray-200 dark:border-gray-800 pt-10 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
    >
      <p class="text-gray-500 dark:text-gray-400 text-sm">Got a link you're not sure about?</p>
      <a
        href="/"
        class="inline-flex items-center gap-2 bg-gray-900 dark:bg-white text-white dark:text-gray-900 px-5 py-2.5 rounded-lg text-sm font-medium hover:bg-gray-700 dark:hover:bg-gray-200 transition-colors"
      >
        Check it on url.vet
        <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
        </svg>
      </a>
    </div>
  </div>
</div>
