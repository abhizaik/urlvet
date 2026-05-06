<div align="center">

# Phishing Detection
Open-source phishing detection engine for real-time URL analysis. Detect malicious links, explain every verdict, and generate a full security report in under a second.


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

 > Paste a URL → get a **trust score, verdict, and detailed report** in under a second.

![Phishing Detection Demo](assets/demo.gif)

Try phishing detection live: https://safesurf.xorwave.com


## Quick Start

```bash
git clone https://github.com/abhizaik/phishing-detection.git
cd phishing-detection
make build && make up
```

Open Web UI: **[localhost:3000](http://localhost:3000)** 

No external API keys required.

Detailed setup guide: [docs/setup.md](docs/setup.md) 



## At a Glance

- 300–700ms end-to-end scan time
- 20+ heuristic analyzers
- HTTP API + Web UI + Chrome extension
- Explainable scoring (no black-box ML)
- One-command Docker setup



## Why This Exists

Most phishing tools are either:

- **Closed commercial APIs:** expensive and opaque
- **Academic prototypes:** not built for production use

This project provides:

- **Explainable detection:** every verdict includes reasons
- **Real-time performance:** parallel analyzers
- **Flexible deployment:** API, UI, browser extension
- **Full control:** open-source and self-hostable


## How It Compares

| Feature | This Project | VirusTotal | Google Safe Browsing | URLVoid | PhishTank | Typical ML Demos |
|---------|--------------|------------|----------------------|---------|-----------|------------------|
| Open source | ✅ | ❌ | ❌ | ❌ | Partial | Usually ✅ |
| Self-hostable | ✅ | ❌ | ❌ | ❌ | ❌ | Usually ✅ |
| Explainable scoring | ✅ Full | ❌ | ❌ | ❌ | ❌ | Rare |
| Real-time heuristics | ✅ | Partial | ❌ | ❌ | ❌ | Varies |
| Live website preview | ✅ | Partial | ❌ | ❌ | ❌ | Rare |
| Beginner-friendly UI | ✅ | Partial | ❌ | Limited | ❌ | Rare |
| Browser extension | ✅ | ❌ | ❌ | ❌ | ❌ | Rare |
| HTTP API | ✅ | ✅ | ✅ | Partial | Limited | Rare |
| Web UI included | ✅ | ✅ | ❌ | ✅ | ❌ | Rare |
| Production ready | ✅ | ✅ | ✅ | ✅ | Database only | Rare |
| Private deployment | ✅ | ❌ | ❌ | ❌ | ❌ | Sometimes |
| Audit / modify scoring logic | ✅ | ❌ | ❌ | ❌ | ❌ | Partial |
| Multi-signal analysis | ✅ | ✅ | Database-based | Partial | Database-based | Usually single-model |
| No API key required | ✅ | ❌ | ❌ | ❌ | ❌ | Usually ✅ |


## Who This Is For

- End users checking suspicious links
- Developers integrating URL analysis
- Security teams building detection pipelines
- Researchers studying phishing infrastructure

Academic or research use of this project must cite this repository (see [CITATION.cff](CITATION.cff)).



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

The engine evaluates URLs using multiple independent analyzers, including:

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

Project is designed to reduce phishing risk, not guarantee safety.
Always combine automated detection with human review and layered defenses.


## Architecture
Modular analyzer-based architecture designed for parallel, real-time URL evaluation.

High level repository layout:

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

All documentation is under `docs/`. Start here [docs/README.md](docs/README.md) 

[Interactive API docs (Swagger UI)](https://api.safesurf.xorwave.com/swagger/index.html)

Full endpoint reference: [docs/api.md](docs/api.md)


## Contributing

- Report bugs → [Issues](https://github.com/abhizaik/phishing-detection/issues) 
- Discussions → [Discussions](https://github.com/abhizaik/phishing-detection/discussions) 
- Code contributions → [CONTRIBUTING.md](.github/CONTRIBUTING.md)
- Feedback is always welcome

**If you found this project helpful, consider giving it a star.** It directly helps visibility and continued development.


<div align="center">
  <a href="https://star-history.com/#abhizaik/phishing-detection&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api/star-history.com/svg?repos=abhizaik/phishing-detection&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date" />
    </picture>
  </a>
</div>

