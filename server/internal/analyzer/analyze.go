package analyzer

import (
	"context"
	"log"
	"time"

	"github.com/abhizaik/SafeSurf/internal/metrics"
	"github.com/abhizaik/SafeSurf/internal/service/cache"
	"github.com/abhizaik/SafeSurf/internal/service/checks"
)

// Analyze runs all tasks and builds the final response
func Analyze(ctx context.Context, rawURL string) (Response, []error) {
	// Validate and normalize the URL (adds https:// if scheme is missing).
	// Using the normalized form as the canonical URL ensures consistent cache
	// keys regardless of how the caller formatted the input.
	parsedURL, isValid, err := checks.IsValidURL(rawURL)
	if err != nil || !isValid {
		return Response{}, []error{ErrInvalidURL}
	}
	normalizedURL := parsedURL.String()

	domain, err := checks.GetDomain(normalizedURL)
	if err != nil {
		return Response{}, []error{err}
	}

	// Initialize cache (non-blocking - if cache fails, continue without it)
	var cacheInstance CacheInterface
	cacheConn, cacheErr := cache.New()
	if cacheErr != nil {
		log.Printf("Warning: Failed to initialize cache: %v. Continuing without cache.", cacheErr)
	} else {
		cacheInstance = cacheConn
		defer cacheConn.Close()
	}

	in := &Input{URL: normalizedURL, Domain: domain, Cache: cacheInstance}

	tasks := []Task{
		rankTask{},
		httpCombinedTask{}, // Optimized: combines redirects, HSTS, and status code
		ipCheckTask{},
		ipResolveTask{},
		punycodeTask{},
		tldTask{},
		shortenerTask{},
		structureTask{},
		keywordsTask{},
		dnsValidityTask{},
		subdomainTask{},
		whoisTask{},
		tlsCombinedTask{}, // Optimized: combines TLS and SSL checks
		entropyTask{},
		contentTask{},
		homoglyphTask{},
		phishtankTask{},
		typosquatTask{},
	}

	// Start timing right before tasks run
	start := time.Now()
	out, errs := runTasks(ctx, in, tasks)

	resp := Response{
		URL:    normalizedURL,
		Domain: domain,
		Features: Features{
			Rank: out.Rank,
			TLD: TLDInfo{
				TLD:       out.TLD,
				IsTrusted: out.TLDTrusted,
				IsRisky:   out.TLDRisky,
				IsICANN:   out.TLDICANN,
			},
			URL: URLChecks{
				IsURLShortener:   out.URLIsShortener,
				UsesIP:           out.URLUsesIP,
				ContainsPunycode: out.URLContainsPuny,
				TooLong:          out.URLTooLong,
				TooDeep:          out.URLTooDeep,
				SubdomainCount:   out.URLSubdomainCount,
				HasHomoglyph:     out.HomoglyphPresent,
				Keywords: Keywords{
					HasKeywords: out.URLKeywordsPresent,
					Found:       out.URLKeywordMatches,
					Categories:  out.URLKeywordCats,
				},
			},
		},
		Infrastructure: Infrastructure{
			IPAddresses:      out.IPs,
			NameserversValid: out.NSValid,
			NSHosts:          out.NSHosts,
			MXRecordsValid:   out.MXValid,
			MXHosts:          out.MXHosts,
		},
		DomainInfo: out.DomainInfo,
		Analysis: Analysis{
			RedirectionResult: out.RedirectionResult,
			SupportsHSTS:      out.SupportsHSTS,
			HTTPStatus: HTTPStatus{
				Code:                 out.StatusCode,
				Text:                 out.StatusText,
				Success:              out.StatusSuccess,
				IsRedirectStatusCode: out.StatusIsRedirect,
			},
		},
		SSLInfo:          out.SSLInfo,
		TLSInfo:          out.TLSInfo,
		ContentData:      out.ContentData,
		DomainRandomness: out.DomainRandomness,
		TyposquatResult:  out.TyposquatResult,
		ThreatIntel: ThreatIntel{
			PhishTank: out.PhishTank,
		},
		Performance: Performance{
			TotalTime: time.Since(start).String(),
			Timings:   ConvertTimings(out.Timings),
		},
	}

	result := GenerateResult(resp)
	resp.Result = result
	metrics.RiskScore.Observe(float64(result.RiskScore))
	metrics.TrustScore.Observe(float64(result.TrustScore))

	if len(errs) > 0 {
		resp.Incomplete = true
		for _, e := range errs {
			resp.Errors = append(resp.Errors, e.Error())
		}
	}

	return resp, errs
}
