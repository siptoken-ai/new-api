package service

import (
	"errors"
	"net/http"

	"github.com/QuantumNous/new-api/i18n"
	kitutil "github.com/QuantumNous/new-api/relaykit/relayconvert/kitutil"
	"github.com/QuantumNous/new-api/relaykit/types"

	"github.com/gin-gonic/gin"
)

// upstreamAccountErrorCodes are 403 codes about the site's own account at the
// upstream rather than the request or the user.
var upstreamAccountErrorCodes = map[types.ErrorCode]bool{
	"account_restricted": true,
	"model_denylisted":   true,
}

// IsUpstreamAccountError reports whether an error returned by the upstream is
// about the site's own account there: its key (401), its balance (402) or a
// restriction on the whole account. Users can do nothing about these, and the
// upstream's wording names the upstream.
func IsUpstreamAccountError(err *types.NewAPIError) bool {
	if err == nil {
		return false
	}
	switch err.GetErrorType() {
	case types.ErrorTypeOpenAIError, types.ErrorTypeClaudeError:
	default:
		return false
	}
	switch err.StatusCode {
	case http.StatusUnauthorized, http.StatusPaymentRequired:
		return true
	case http.StatusForbidden:
		return upstreamAccountErrorCodes[err.GetErrorCode()]
	}
	return false
}

// HideUpstreamAccountError replaces an upstream account error of a request that
// reached a channel with a generic 503 when HIDE_UPSTREAM_IDENTITY is on. The
// server log keeps the original.
func HideUpstreamAccountError(c *gin.Context, err *types.NewAPIError) *types.NewAPIError {
	if !kitutil.HideUpstreamIdentity.Load() || len(c.GetStringSlice("use_channel")) == 0 || !IsUpstreamAccountError(err) {
		return err
	}
	return types.NewErrorWithStatusCode(
		errors.New(i18n.T(c, i18n.MsgRelayUpstreamUnavailable)),
		"service_unavailable",
		http.StatusServiceUnavailable,
		types.ErrOptionWithSkipRetry(),
	)
}

// upstreamAccountErrorLogText is the error log text users see for an upstream
// account error when HIDE_UPSTREAM_IDENTITY is on.
func upstreamAccountErrorLogText(c *gin.Context, err *types.NewAPIError) (string, bool) {
	if !kitutil.HideUpstreamIdentity.Load() || !IsUpstreamAccountError(err) {
		return "", false
	}
	return i18n.T(c, i18n.MsgRelayUpstreamUnavailable), true
}
