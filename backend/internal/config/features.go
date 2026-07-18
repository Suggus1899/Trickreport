package config

// FeatureFlags holds toggles for optional product capabilities.
//
// Flags are intentionally coarse-grained: each one gates an entire feature
// surface (MFA, attachments, automations, etc.) so operators can enable or
// disable functionality per environment without code changes.
//
// TODO(fase-12): wire FeatureFlags into Config.Load() so flags can be read from
// environment variables (e.g. FEATURE_ENABLE_MFA=true) and merged with
// DefaultFeatures(). Do NOT edit config.go as part of this stub — wire it in a
// follow-up change.
type FeatureFlags struct {
	// EnableMFA toggles multi-factor authentication on login.
	EnableMFA bool
	// EnableAttachments toggles ticket file attachments (upload/download).
	EnableAttachments bool
	// EnableAutomation toggles the SLA automation engine and its admin UI.
	EnableAutomation bool
	// EnableFullTextSearch toggles full-text search over tickets and KB articles.
	EnableFullTextSearch bool
	// EnableAnalytics toggles the analytics dashboard and metrics endpoints.
	EnableAnalytics bool
	// EnableEmailQueue toggles the asynchronous email queue/sender.
	EnableEmailQueue bool
}

// DefaultFeatures returns a FeatureFlags value with every feature enabled.
// This preserves backward compatibility: existing deployments keep all
// functionality on by default until operators explicitly disable a flag.
func DefaultFeatures() FeatureFlags {
	return FeatureFlags{
		EnableMFA:            true,
		EnableAttachments:    true,
		EnableAutomation:     true,
		EnableFullTextSearch: true,
		EnableAnalytics:      true,
		EnableEmailQueue:     true,
	}
}

// IsEnabled reports whether the named feature is enabled.
//
// Accepted names (case-insensitive):
//
//	mfa | attachments | automation | fulltextsearch | analytics | emailqueue
//
// Unknown names return false so callers can safely probe for features that may
// be added in the future without panicking.
func (f FeatureFlags) IsEnabled(name string) bool {
	switch normalize(name) {
	case "mfa":
		return f.EnableMFA
	case "attachments":
		return f.EnableAttachments
	case "automation":
		return f.EnableAutomation
	case "fulltextsearch":
		return f.EnableFullTextSearch
	case "analytics":
		return f.EnableAnalytics
	case "emailqueue":
		return f.EnableEmailQueue
	default:
		return false
	}
}

// normalize lowercases and strips common separators so callers can pass either
// "FullTextSearch" or "full_text_search" / "fulltext-search".
func normalize(name string) string {
	out := make([]byte, 0, len(name))
	for i := 0; i < len(name); i++ {
		c := name[i]
		switch c {
		case '_', '-', ' ', '.':
			continue
		default:
			if c >= 'A' && c <= 'Z' {
				c += 'a' - 'A'
			}
			out = append(out, c)
		}
	}
	return string(out)
}
