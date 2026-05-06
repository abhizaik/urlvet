<div align="center">

# Phishing Detection
Open-source phishing detection engine for real-time URL analysis. Detect malicious links, explain every verdict, and generate a security report in real time.


[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go\&logoColor=white)](https://go.dev)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange?logo=svelte)](https://svelte.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/abhizaik/phishing-detection?style=social)](https://github.com/abhizaik/phishing-detection)
[![Last Commit](https://img.shields.io/github/last-commit/abhizaik/phishing-detection)](https://github.com/abhizaik/phishing-detection/commits/main)
[![Docs](https://img.shields.io/badge/Swagger-Docs-brightgreen)](https://api.safesurf.xorwave.com/swagger/index.html)

[⚡ Quick Start](#quick-start) · [⚙️ Detection Engine](#detection-engine) · [🏛 Architecture](#architecture) · [📚 Docs](#documentation) · [🤝 Contributing](#contributing)

</div>

---

## Phishing Detection Demo

 > Paste a URL → get a **trust score, verdict, and detailed report** in real time.

![Phishing Detection Demo](assets/demo.gif)

Live demo: https://safesurf.xorwave.com


## Quick Start

```bash
git clone https://github.com/abhizaik/phishing-detection.git
cd phishing-detection
make build && make up
```

Open Web UI: **[localhost:3000](http://localhost:3000)** 

Detailed setup guide: [docs/setup.md](docs/setup.md) 



## At a Glance

- Live scan, instant results
- 20+ heuristic analyzers
- HTTP API + Web UI + Chrome extension
- Explainable scoring (no black-box ML)
- One-command Docker setup


## How It Compares

| Feature | SafeSurf | VirusTotal | Google Safe Browsing | URLScan.io | CheckPhish |
|---------|----------|------------|----------------------|------------|------------|
| Live crawl, instant results | ✅ | Partial | ❌ | Partial | Partial |
| Explains every verdict | ✅ | Partial | ❌ | Partial | Partial |
| Beginner-friendly interface | ✅ | Partial | Partial | Partial | Partial |
| Credential form detection | ✅ | ❌ | ❌ | Partial | ✅ |
| Follows redirect chains | ✅ | ✅ | ❌ | ✅ | ✅ |
| Detailed technical insights | ✅ | ❌ | ❌ | ✅ | Partial |
| Live page preview | ✅ | ❌ | ❌ | ✅ | ✅ |
| Detection using AI/ML | ❌ | ✅ | ✅ | Partial | ✅ |
| Known phishing database coverage | Partial | ✅ | ✅ | Partial | Partial |
| Scan multiple URLs at once | ❌ | ✅ | ✅ | ✅ | ❌ |
| Browser protection | ✅ | ✅ | ✅ | ✅ | ❌ |
| Open source | ✅ | ❌ | ❌ | ❌ | ❌ |

Fast scanners give you a verdict with no context. Deep crawlers take too long. SafeSurf bridges the gap by doing live analysis with per-signal explanations in real time — and it's open-source.



## Who This Is For

- End users checking suspicious links
- Developers integrating URL analysis
- Security teams building detection pipelines
- Researchers



## API Example

Analyze a URL via HTTP:

```bash
curl -X GET http://localhost:8080/api/v1/analyze?url=https://example.com
```
**Sample Response:**
<pre><code class="language-json">
{
  "url": "https://example.com",
  "trust_score": 100,
  "verdict": "Safe",
  "reasons": {
    "good_reasons": [...]
  }
}
</code></pre>
Full response schema → [docs/api.md#example](docs/api.md#example) 

## Detection Engine

SafeSurf evaluates URLs using multiple independent analyzers, including:

- Domain reputation & age checks
- Suspicious URL patterns and homoglyphs
- Redirect chain analysis
- HTTPS / certificate anomalies
- Known phishing indicators and heuristics
- Content-based signals (HTML, scripts, forms)

Each analyzer contributes to a final trust score and verdict.



## Limitations
- Heuristic-based detection may produce false positives
- No ML model (intentional, prioritizes explainability and auditability)

Not a safety guarantee. Use alongside other defenses.


## Architecture
SafeSurf runs analyzers in parallel. High-level repo layout:

```text
server/               Go backend 
  cmd/safesurf        Backend entry point
  internal/           Analyzers, domaininfo, screenshot
web/website           SvelteKit UI
web/chrome-extension  Chrome extension
docker/               Dev & prod
docs/                 Setup, architecture, API, security, testing etc.
Makefile
```



## Documentation

- [docs/README.md](docs/README.md) — start here
- [docs/api.md](docs/api.md) — full API reference
- [Swagger UI](https://api.safesurf.xorwave.com/swagger/index.html) — interactive API docs


## Citation

If you use this project in academic or research work, please cite it — see [CITATION.cff](CITATION.cff).


## Contributing

- Found a bug? → [Open an issue](https://github.com/abhizaik/phishing-detection/issues)
- Have a question or idea? → [Start a discussion](https://github.com/abhizaik/phishing-detection/discussions)
- Want to contribute code? → [CONTRIBUTING.md](.github/CONTRIBUTING.md)

If you found this project helpful, consider giving it a star.




<div align="center">
  <a href="https://star-history.com/#abhizaik/phishing-detection&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date" />
    </picture>
  </a>
</div>

