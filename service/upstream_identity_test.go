package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/i18n"
	kitutil "github.com/QuantumNous/new-api/relaykit/relayconvert/kitutil"
	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIsUpstreamAccountError(t *testing.T) {
	tests := []struct {
		name string
		err  *types.NewAPIError
		want bool
	}{
		{"upstream rejects the site's key", types.WithOpenAIError(types.OpenAIError{Message: "Malformed API key", Code: "invalid_api_key"}, http.StatusUnauthorized), true},
		{"upstream balance used up", types.WithOpenAIError(types.OpenAIError{Message: "insufficient balance"}, http.StatusPaymentRequired), true},
		{"whole account restricted, Claude format", types.WithClaudeError(types.ClaudeError{Type: "account_restricted", Message: "restricted"}, http.StatusForbidden), true},
		{"model denied for the whole account", types.WithOpenAIError(types.OpenAIError{Code: "model_denylisted"}, http.StatusForbidden), true},
		{"one end user blocked", types.WithOpenAIError(types.OpenAIError{Code: "request_blocked"}, http.StatusForbidden), false},
		{"bad request", types.WithOpenAIError(types.OpenAIError{Code: "invalid_request"}, http.StatusBadRequest), false},
		{"the site's own quota error", types.NewErrorWithStatusCode(errors.New("quota"), types.ErrorCodeInsufficientUserQuota, http.StatusPaymentRequired), false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsUpstreamAccountError(tt.err))
		})
	}
}

func TestHideUpstreamAccountError(t *testing.T) {
	newContext := func(usedChannel bool) *gin.Context {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
		c.Request.Header.Set("Accept-Language", "en")
		if usedChannel {
			c.Set("use_channel", []string{"2"})
		}
		return c
	}
	upstreamErr := types.WithOpenAIError(types.OpenAIError{Message: "insufficient balance"}, http.StatusPaymentRequired)

	assert.Same(t, upstreamErr, HideUpstreamAccountError(newContext(true), upstreamErr), "off by default")

	kitutil.HideUpstreamIdentity.Store(true)
	t.Cleanup(func() { kitutil.HideUpstreamIdentity.Store(false) })

	assert.Same(t, upstreamErr, HideUpstreamAccountError(newContext(false), upstreamErr), "request reached no channel")

	hidden := HideUpstreamAccountError(newContext(true), upstreamErr)
	assert.Equal(t, http.StatusServiceUnavailable, hidden.StatusCode)
	assert.Equal(t, types.ErrorCode("service_unavailable"), hidden.GetErrorCode())
	assert.Equal(t, i18n.Translate("en", i18n.MsgRelayUpstreamUnavailable), hidden.Error())
	assert.NotContains(t, hidden.ToOpenAIError().Message, "insufficient balance")

	blocked := types.WithOpenAIError(types.OpenAIError{Code: "request_blocked"}, http.StatusForbidden)
	assert.Same(t, blocked, HideUpstreamAccountError(newContext(true), blocked))
}
