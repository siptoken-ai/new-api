package kitutil

import (
	"net"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/MoeForge/moego/regexp"
)

// HideUpstreamIdentity makes MaskSensitiveInfo also mask the upstream's site:
// every address on the registrable domain of an upstream host, and the site
// name taken from it. Upstream error messages often name the provider and link
// to pages on its main domain ("Get a key at https://example.com/api") rather
// than the API host, which the default masking misses. The host sets it from
// HIDE_UPSTREAM_IDENTITY at startup; it is off by default.
var HideUpstreamIdentity atomic.Bool

// upstreamSiteNameMinLength skips site names short enough to be ordinary words
// or parts of them.
const upstreamSiteNameMinLength = 4

// maskUpstreamSites masks the registrable domain of each upstream host, its
// subdomains and the site name, as in api.example.com -> example.com and
// "example". Example with upstream https://api.example.com:
// Example keys start with sk-. Get one at https://example.com/api.
// -> *** keys start with sk-. Get one at https://***.com/***.
func maskUpstreamSites(str string, upstreamHosts []string) string {
	for _, host := range upstreamHosts {
		if net.ParseIP(host) != nil {
			continue
		}
		parts := strings.Split(strings.ToLower(host), ".")
		tail := maskHostTail(parts)
		if len(parts) <= len(tail) {
			continue
		}
		name := parts[len(parts)-len(tail)-1]
		if len(name) < upstreamSiteNameMinLength {
			continue
		}
		domain := regexp.QuoteMeta(name + "." + strings.Join(tail, "."))

		urlPattern := regexp.MustCompile(`(?i)(?:https?|wss?)://(?:[^\s"<>{}/?#@]*@)?(?:[a-z0-9-]+\.)*` + domain + `\b[^\s"<>{}]*`)
		str = urlPattern.ReplaceAllStringFunc(str, func(match string) string {
			urlStr := strings.TrimRight(match, ".,;:!?')")
			suffix := match[len(urlStr):]
			u, err := url.Parse(urlStr)
			if err != nil {
				return "***" + suffix
			}
			result := u.Scheme + "://" + maskHostForURL(strings.TrimSuffix(strings.ToLower(u.Hostname()), "."))
			if u.Path != "" && u.Path != "/" {
				result += "/***"
			}
			return result + suffix
		})

		hostPattern := regexp.MustCompile(`(?i)\b(?:[a-z0-9-]+\.)*` + domain + `\b`)
		str = hostPattern.ReplaceAllStringFunc(str, func(match string) string {
			return maskHostForPlainDomain(strings.ToLower(match))
		})

		namePattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(name) + `\b`)
		str = namePattern.ReplaceAllString(str, "***")
	}
	return str
}
