<div align="center">

# Phishing Detection Engine (Open-Source, Real-Time URL Scanner)

Fast, explainable phishing detection for URLs — real-time scoring, clear verdicts, full transparency.

Open-source and production-ready. Analyze URLs in under a second with transparent scoring and detailed reports.

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go\&logoColor=white)](https://go.dev)
[![Svelte](https://img.shields.io/badge/Svelte-5-orange?logo=svelte)](https://svelte.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/abhizaik/phishing-detection?style=social)](https://github.com/abhizaik/phishing-detection)

[⚡ Quick Start](#quick-start) · [🏛 Architecture](#architecture) · [📚 Docs](#documentation) · [🤝 Contributing](#contributing) · [🌍 Community](#community)

</div>

---

## Phishing Detection Demo

> Paste a URL → get a **trust score, verdict, and detailed report** in under a second.

![Phishing Detection Demo](assets/demo.gif)


https://safesurf.xorwave.com
---

## Quick Start

Full setup: [docs/setup.md](docs/setup.md) 

1. Clone the repo

```bash
git clone https://github.com/abhizaik/phishing-detection.git
cd phishing-detection
```
2. Start the application (backend + UI via Docker)

Prerequisite: Docker must be installed and running.

Windows: Use WSL or install make.
```bash
make build
make up
```
Web UI: **[localhost:3000](http://localhost:3000)** 


---

## Features

* Scans URLs for **phishing, malicious behavior, and unsafe redirects**
* Produces a **trust score, clear verdict, and detailed report**
* Supports **developers and non-technical users** via UI, API, and extension
* Uses **multiple independent heuristic analyzers** for accurate detection
* Built with **Go (backend)** and **Svelte (frontend)** for production use

---



## Use Cases

- Detect phishing links before users click them
- Scan URLs for malicious behavior
- Build anti-phishing browser extensions
- Integrate phishing detection into backend services
- Replace or supplement commercial phishing APIs



## Why Use This Tool?

Most phishing detection solutions are either **closed commercial APIs** or **academic ML demos**:

* **Commercial tools**: expensive, opaque, and impossible to audit
* **ML demos**: slow, fragile, and not built for real-world deployment

**Phishing remains a top cyber threat** because defenders lack **fast, explainable, and controllable detection systems**.
This engine fills that gap by providing:

* **Transparent, explainable analysis** — every verdict is backed by concrete signals
* **Fast, real-time scanning** — multiple analyzers run in parallel
* **Flexible integration** — web UI, HTTP API, browser extension
* **Full open-source control** — audit, modify, self-host, and scale



## Who This Is For

**Everyday users**

* Quickly check suspicious URLs in a website or browser extension

**Developers**

* Integrate phishing detection into applications or backend services
* Replace or supplement commercial phishing APIs

**Security engineers & SOC teams**

* Build explainable phishing detection pipelines
* Audit URLs with transparent, actionable signals

**Students & researchers**

* Use a real-world, production-grade reference for **academic or security projects**. Academic or research use of this project must cite this repository (see [CITATION.cff](CITATION.cff)).



## API Example

The phishing detection engine exposes a simple HTTP API for real-time URL analysis. 

Returns a detailed structured analysis including domain info, SSL, redirects, and final verdict.


Scan a URL using the API:

```bash
curl -X GET http://localhost:8080/api/v1/analyze?url=https://example.com
```
<details>
<summary>Example API response</summary>
 <pre><code class="language-json">
{
  "url": "https://example.com",
  "domain": "example.com",
  "features": {
    "rank": 175,
    "tld": {
      "tld": "com",
      "is_trusted_tld": false,
      "is_risky_tld": false,
      "is_icann": true
    },
    "url": {
      "url_shortener": false,
      "uses_ip": false,
      "contains_punycode": false,
      "too_long": false,
      "too_deep": false,
      "has_homoglyph": false,
      "subdomain_count": 0,
      "keywords": {
        "has_keywords": false,
        "found": [],
        "categories": {}
      }
    }
  },
  "infrastructure": {
    "ip_addresses": [
      "172.66.147.243",
      "104.20.23.154",
      "2606:4700:10::6814:179a",
      "2606:4700:10::ac42:93f3"
    ],
    "nameservers_valid": true,
    "ns_hosts": [
      "hera.ns.cloudflare.com."
    ],
    "mx_records_valid": false,
    "mx_hosts": [
      "."
    ]
  },
  "domain_info": {
    "domain": "EXAMPLE.COM",
    "registrar": "RESERVED-Internet Assigned Numbers Authority",
    "created": "1995-08-14T04:00:00Z",
    "updated": "2026-01-16T18:26:50Z",
    "expiry": "2026-08-13T04:00:00Z",
    "nameservers": [
      "ELLIOTT.NS.CLOUDFLARE.COM",
      "HERA.NS.CLOUDFLARE.COM"
    ],
    "status": [
      "client delete prohibited",
      "client transfer prohibited",
      "client update prohibited"
    ],
    "dnssec": true,
    "age_human": "30 years 8 months",
    "age_days": 11202,
    "raw": "{\"ldhName\":\"EXAMPLE.COM\",\"nameservers\":[{\"ldhName\":\"ELLIOTT.NS.CLOUDFLARE.COM\"},{\"ldhName\":\"HERA.NS.CLOUDFLARE.COM\"}],\"events\":[{\"eventAction\":\"registration\",\"eventDate\":\"1995-08-14T04:00:00Z\"},{\"eventAction\":\"expiration\",\"eventDate\":\"2026-08-13T04:00:00Z\"},{\"eventAction\":\"last changed\",\"eventDate\":\"2026-01-16T18:26:50Z\"},{\"eventAction\":\"last update of RDAP database\",\"eventDate\":\"2026-04-15T19:04:14Z\"}],\"entities\":[{\"roles\":[\"registrar\"],\"vcardArray\":[\"vcard\",[[\"version\",{},\"text\",\"4.0\"],[\"fn\",{},\"text\",\"RESERVED-Internet Assigned Numbers Authority\"]]]}],\"status\":[\"client delete prohibited\",\"client transfer prohibited\",\"client update prohibited\"],\"secureDNS\":{\"delegationSigned\":true}}",
    "source": "RDAP"
  },
  "analysis": {
    "redirection_result": {
      "is_redirected": false,
      "chain_length": 1,
      "chain": [
        "https://example.com"
      ],
      "final_url": "https://example.com",
      "final_url_domain": "example.com",
      "has_domain_jump": false
    },
    "http_status": {
      "code": 200,
      "text": "OK",
      "success": true,
      "is_redirect": false
    },
    "is_hsts_supported": false
  },
  "ssl_info": {
    "Domain": "example.com",
    "HasTLS": true,
    "ChainValid": true,
    "Issuer": "Cloudflare TLS Issuing ECC CA 1",
    "NotBefore": "2026-04-02T21:18:57Z",
    "NotAfter": "2026-07-01T21:24:46Z",
    "AgeDays": 12,
    "Fingerprint": "1AF627C6C2AC992E3C9102438F467C4C238D3112325AC7CF9003D77F75EFFFBA",
    "IsSuspicious": false,
    "Reasons": null,
    "CTLogged": true,
    "KnownBadChain": false
  },
  "tls_info": {
    "Present": true,
    "Issuer": "CLOUDFLARE, INC.",
    "AgeDays": 12,
    "HostnameMismatch": false
  },
  "content_data": {
    "url": "https://example.com",
    "title": "Example Domain",
    "has_forms": false,
    "has_login_form": false,
    "has_payment_form": false,
    "has_personal_form": false,
    "form_count": 0,
    "forms": null,
    "iframes": null,
    "has_hidden_iframe": false,
    "has_tracking": false,
    "fetch_duration": 137804093,
    "brand_check": {
      "brand_found": "",
      "is_mismatch": false,
      "detected_names": []
    }
  },
  "domain_randomness": {
    "Domain": "example.com",
    "Label": "example",
    "Length": 7,
    "Entropy": 2.521640636343318,
    "EntropyPerChar": 0.36023437662047403,
    "NormalizedEntropy": 0.06050092369175979,
    "VowelRatio": 0.42857142857142855,
    "DigitRatio": 0,
    "UniqueCharRatio": 0.8571428571428571,
    "LongestConsonantRun": 3,
    "BigramEnglishiness": 0.16666666666666666,
    "RandomnessScore": 0.3567918975896066,
    "IsSuspicious": false,
    "Reasons": []
  },
  "typosquat_result": {
    "is_suspicious": false
  },
  "phishing": {
    "in_database": true,
    "phish_id": 7366538,
    "phish_detail_page": "http://www.phishtank.com/phish_detail.php?phish_id=7366538",
    "verified": false,
    "verified_at": "",
    "valid": false,
    "target": "",
    "source": "phishtank",
    "from_cache": false,
    "raw_response": {
      "meta": {
        "timestamp": "2026-04-15T19:04:30+00:00",
        "serverid": "e5f3084e",
        "status": "success",
        "requestid": "172.17.128.1.69dfe13e5ee121.10644345"
      },
      "results": {
        "url": "https://example.com",
        "in_database": true,
        "phish_id": 7366538,
        "phish_detail_page": "http://www.phishtank.com/phish_detail.php?phish_id=7366538",
        "verified": false,
        "verified_at": null,
        "valid": false
      }
    }
  },
  "result": {
    "risk_score": 5,
    "trust_score": 100,
    "final_score": 98,
    "verdict": "Safe",
    "reasons": {
      "neutral_reasons": [
        "Standard, officially recognized domain extension.",
        "No email server configured for this domain."
      ],
      "good_reasons": [
        "Global Giant: Ranked #175 worldwide.",
        "Long-standing domain history (30 years 8 months).",
        "Advanced DNS security enabled (DNSSEC)."
      ],
      "bad_reasons": null
    }
  },
  "incomplete": false,
  "errors": null
}
</code></pre>
</details>


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




## Detection Engine

The engine evaluates URLs using multiple independent analyzers, including:

- Domain reputation & age checks
- Suspicious URL patterns and homoglyphs
- Redirect chain analysis
- HTTPS / certificate anomalies
- Known phishing indicators and heuristics
- Content-based signals (HTML, scripts, forms)

Each analyzer contributes to a final trust score and verdict.



## Performance

- Typical scan time: **~300–700 ms per URL**
- Designed to handle multiple concurrent scans efficiently
- Optimized for **real-time phishing detection at scale**

Exact performance depends on enabled analyzers and network conditions.



## Limitations
- Heuristic-based detection may produce false positives
- No ML model (intentional — prioritizes explainability and auditability)
- Accuracy depends on external signals (DNS, SSL etc.)

## Documentation

All documentation is under `docs/`. Start here [docs/README.md](docs/README.md) 

[Interactive API docs (Swagger UI)](https://api.safesurf.xorwave.com/swagger/index.html)

Full endpoint reference: [docs/api.md](docs/api.md)



## Contributing

Bug reports, feature requests, and pull requests are welcome.

Use [GitHub Issues](https://github.com/abhizaik/phishing-detection/issues) to report bugs or suggest features. For code contributions, see [CONTRIBUTING.md](.github/CONTRIBUTING.md).




## Community

**If you found this project helpful, consider giving it a star.** It directly helps visibility and continued development.

Have bugs, ideas, or feature requests?
Open an [issue](https://github.com/abhizaik/phishing-detection/issues) or start a [discussion](https://github.com/abhizaik/phishing-detection/discussions). Contributions and feedback are welcome.

Thanks for helping make the web safer.




<div align="center">
  <a href="https://star-history.com/#abhizaik/phishing-detection&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api/star-history.com/svg?repos=abhizaik/phishing-detection&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=Date" />
    </picture>
  </a>
</div>

<!-- [![Community Growth Trajectory](https://api.star-history.com/svg?repos=abhizaik/phishing-detection&type=date&legend=top-left)](https://www.star-history.com/#abhizaik/phishing-detection&type=date&legend=top-left) -->
