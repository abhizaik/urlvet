<svelte:head>
  <title>How It Works — url.vet (URLvet)</title>
  <meta
    name="description"
    content="How url.vet scans links, what it checks, and what each result means."
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
    <p class="text-gray-500 dark:text-gray-400 text-lg mb-12">
      Paste a link. Get an answer. Here's what happens and how to read it.
    </p>

    <!-- Step 1: Using it -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-5">How to scan a link</h2>
      <ol class="space-y-5">
        {#each [{ n: "1", title: "Paste any link", desc: "Drop a URL into the input box. Full link, shortened link, or just a domain like example.com. url.vet handles the formatting." }, { n: "2", title: "Click Scan", desc: "The scan runs live. It takes a few seconds because url.vet fetches the page in real time rather than looking it up in a cached index." }, { n: "3", title: "Read the result", desc: "You'll get a trust score from 0 to 100, a verdict, and a breakdown of every check that ran." }] as item}
          <li class="flex gap-4">
            <span
              class="flex-shrink-0 w-7 h-7 rounded-full bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300 text-sm font-semibold flex items-center justify-center mt-0.5"
              >{item.n}</span
            >
            <div>
              <p class="text-sm font-semibold mb-0.5">{item.title}</p>
              <p class="text-base text-gray-500 dark:text-gray-400">{item.desc}</p>
            </div>
          </li>
        {/each}
      </ol>
    </section>

    <!-- Understanding the output -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-5">What you get back</h2>

      <!-- Score -->
      <div class="mb-8">
        <p class="text-sm font-semibold mb-2">The trust score</p>
        <p class="text-base text-gray-600 dark:text-gray-400 mb-4">
          A number from 0 to 100. Higher is safer. 50 is the neutral baseline. A brand new URL with
          no signals in either direction starts there. Most legitimate sites score above 65.
          Anything below 30 is worth treating as dangerous.
        </p>
        <div class="flex flex-wrap gap-3">
          <div
            class="flex items-center gap-2.5 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800/40 rounded-lg px-3.5 py-2.5"
          >
            <span class="w-2.5 h-2.5 rounded-full bg-red-500 flex-shrink-0"></span>
            <div>
              <p class="text-xs font-semibold text-red-700 dark:text-red-300">Risky</p>
              <p class="text-xs text-red-600/70 dark:text-red-400/70">score below 30</p>
            </div>
          </div>
          <div
            class="flex items-center gap-2.5 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800/40 rounded-lg px-3.5 py-2.5"
          >
            <span class="w-2.5 h-2.5 rounded-full bg-yellow-500 flex-shrink-0"></span>
            <div>
              <p class="text-xs font-semibold text-yellow-700 dark:text-yellow-300">Suspicious</p>
              <p class="text-xs text-yellow-600/70 dark:text-yellow-400/70">score 30 to 64</p>
            </div>
          </div>
          <div
            class="flex items-center gap-2.5 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800/40 rounded-lg px-3.5 py-2.5"
          >
            <span class="w-2.5 h-2.5 rounded-full bg-green-500 flex-shrink-0"></span>
            <div>
              <p class="text-xs font-semibold text-green-700 dark:text-green-300">Safe</p>
              <p class="text-xs text-green-600/70 dark:text-green-400/70">score 65 and above</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Breakdown -->
      <div class="mb-8">
        <p class="text-sm font-semibold mb-2">The signal breakdown</p>
        <p class="text-base text-gray-600 dark:text-gray-400">
          Below the score you'll see every check that ran, grouped into sections: URL structure,
          DNS, TLS, domain intelligence, content, and threat feeds. Each one shows a green flag or a
          red flag. Red flags push the risk score up. Green flags build trust.
        </p>
      </div>

      <!-- Screenshot -->
      <div class="mb-8">
        <p class="text-sm font-semibold mb-2">The page preview</p>
        <p class="text-base text-gray-600 dark:text-gray-400">
          url.vet takes a live screenshot of the page. It's one of the fastest ways to spot a
          phishing site. If the page looks like your bank's login screen but the domain has nothing
          to do with your bank, that's something no automated check can fully catch.
        </p>
      </div>

      <!-- Sharing -->
      <div>
        <p class="text-sm font-semibold mb-2">Sharing a result</p>
        <p class="text-base text-gray-600 dark:text-gray-400">
          Every scan has a permanent shareable URL. Use it to send a result to a colleague, post it
          in a security thread, or report a suspicious link to someone who needs context. The link
          includes the verdict and score so whoever you send it to sees the result without having to
          re-scan.
        </p>
      </div>
    </section>

    <!-- Under the hood (condensed) -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-2">The checks</h2>
      <p class="text-base text-gray-600 dark:text-gray-400">
        18 checks run at the same time the moment you submit. Each one is independent, so a timeout
        or failure in one never holds up the rest.
      </p>
      <br />

      <!-- Pipeline image -->
      <div class="mb-4">
        <img
          src="/pipeline.png"
          alt="url.vet analyzer pipeline showing 18 checks across 7 signal categories"
          class="w-full rounded-xl border border-gray-200 dark:border-gray-800 dark:invert"
          loading="lazy"
        />
      </div>
      <p class="text-base text-gray-600 dark:text-gray-400">
        Checks run in parallel, not in sequence. The slowest check sets the total scan time, not the
        sum of all checks. Scores are calculated once all signals are collected. The score reflects
        the weight of everything combined.
      </p>

      <br />

      <div class="space-y-3">
        {#each [{ label: "URL structure", desc: "Inspects the link before making any network request. Checks for IP addresses used as hostnames, URL shorteners, suspicious keywords in the path, IDN homograph attacks, and unusually deep subdomains." }, { label: "HTTP / Network", desc: "Makes one real request and follows every redirect. Checks HSTS, status code, and whether the final destination differs from the link you pasted." }, { label: "DNS", desc: "Verifies NS and MX records exist and that the domain resolves to a real IP." }, { label: "TLS / SSL", desc: "Checks certificate validity, expiry, issuer, Certificate Transparency log inclusion, and known-bad fingerprints." }, { label: "Domain intelligence", desc: "Looks up domain age via WHOIS, global traffic rank, TLD classification, DNSSEC status, character randomness in the domain name, and typosquatting similarity against 500+ brands." }, { label: "Content analysis", desc: "Fetches and parses the page. Detects login and payment forms on suspicious domains, hidden iframes, brand impersonation, and forms that submit data to external servers." }, { label: "Threat intelligence", desc: "Checks the URL against PhishTank's databases of confirmed and reported phishing links." }] as item}
          <div class="flex gap-3 py-3 border-b border-gray-100 dark:border-gray-800 last:border-0">
            <p class="text-sm font-medium w-40 flex-shrink-0 text-gray-700 dark:text-gray-300">
              {item.label}
            </p>
            <p class="text-base text-gray-500 dark:text-gray-400">{item.desc}</p>
          </div>
        {/each}
      </div>
    </section>

    <!-- Limitations -->
    <section class="mb-14">
      <h2 class="text-xl font-semibold mb-4">Limitations</h2>
      <div class="space-y-3 text-base text-gray-600 dark:text-gray-400 leading-relaxed">
        <p>
          Heuristic detection means false positives are possible. A legitimate site that's new and
          unranked might score lower than it deserves.
        </p>
        <p>
          There's no ML model as of now, and that's intentional. A model that can't explain its
          reasoning isn't useful when trust is on the line.
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
