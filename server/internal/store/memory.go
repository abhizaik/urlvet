package store

import (
	"sync"
	"time"
)

const (
	MaxRecentScans  = 100
	MaxRecentErrors = 50
)

type ScanRecord struct {
	URL      string    `json:"url"`
	Domain   string    `json:"domain"`
	Verdict  string    `json:"verdict"`
	Score    int       `json:"score"`
	Duration string    `json:"duration"`
	Time     time.Time `json:"time"`
	Cached   bool      `json:"cached"`
}

type ErrorRecord struct {
	Task  string    `json:"task"`
	Error string    `json:"error"`
	URL   string    `json:"url"`
	Time  time.Time `json:"time"`
}

var global = &MemoryStore{
	scans:  make([]ScanRecord, 0, MaxRecentScans),
	errors: make([]ErrorRecord, 0, MaxRecentErrors),
}

type MemoryStore struct {
	mu     sync.RWMutex
	scans  []ScanRecord
	errors []ErrorRecord
}

func AddScan(r ScanRecord) {
	global.mu.Lock()
	defer global.mu.Unlock()
	if len(global.scans) >= MaxRecentScans {
		global.scans = global.scans[1:]
	}
	global.scans = append(global.scans, r)
}

func AddError(r ErrorRecord) {
	global.mu.Lock()
	defer global.mu.Unlock()
	if len(global.errors) >= MaxRecentErrors {
		global.errors = global.errors[1:]
	}
	global.errors = append(global.errors, r)
}

func RecentScans() []ScanRecord {
	global.mu.RLock()
	defer global.mu.RUnlock()
	out := make([]ScanRecord, len(global.scans))
	copy(out, global.scans)
	// return newest first
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func RecentErrors() []ErrorRecord {
	global.mu.RLock()
	defer global.mu.RUnlock()
	out := make([]ErrorRecord, len(global.errors))
	copy(out, global.errors)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

type Stats struct {
	TotalScansToday int                `json:"total_scans_today"`
	TotalScansAll   int                `json:"total_scans_all"`
	CacheHits       int                `json:"cache_hits"`
	CacheMisses     int                `json:"cache_misses"`
	CacheHitRate    float64            `json:"cache_hit_rate"`
	AvgDurationMs   float64            `json:"avg_duration_ms"`
	TopDomains      []DomainCount      `json:"top_domains"`
	VerdictCounts   map[string]int     `json:"verdict_counts"`
}

type DomainCount struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

func GetStats() Stats {
	global.mu.RLock()
	defer global.mu.RUnlock()

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var todayCount, cacheHits, cacheMisses int
	var totalDurationMs float64
	domainMap := map[string]int{}
	verdictMap := map[string]int{}

	for _, s := range global.scans {
		if s.Time.After(todayStart) {
			todayCount++
		}
		if s.Cached {
			cacheHits++
		} else {
			cacheMisses++
		}
		domainMap[s.Domain]++
		verdictMap[s.Verdict]++

		// Parse duration string like "1.234s" or "234ms"
		if d, err := time.ParseDuration(s.Duration); err == nil {
			totalDurationMs += float64(d.Milliseconds())
		}
	}

	total := len(global.scans)
	var avgDur float64
	if total > 0 {
		avgDur = totalDurationMs / float64(total)
	}

	var hitRate float64
	if cacheHits+cacheMisses > 0 {
		hitRate = float64(cacheHits) / float64(cacheHits+cacheMisses) * 100
	}

	// Top 10 domains
	type kv struct {
		k string
		v int
	}
	var sorted []kv
	for k, v := range domainMap {
		sorted = append(sorted, kv{k, v})
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].v > sorted[i].v {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	top := []DomainCount{}
	for i, kv := range sorted {
		if i >= 10 {
			break
		}
		top = append(top, DomainCount{Domain: kv.k, Count: kv.v})
	}

	return Stats{
		TotalScansToday: todayCount,
		TotalScansAll:   total,
		CacheHits:       cacheHits,
		CacheMisses:     cacheMisses,
		CacheHitRate:    hitRate,
		AvgDurationMs:   avgDur,
		TopDomains:      top,
		VerdictCounts:   verdictMap,
	}
}
