package kitutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskUpstreamSites(t *testing.T) {
	hosts := []string{"api.relaysite.ai"}
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "site name and a link to the main domain",
			input: `Malformed API key. Relaysite API keys start with sk-; paste only the key itself. Get a key at https://relaysite.ai/api.`,
			want:  `Malformed API key. *** API keys start with sk-; paste only the key itself. Get a key at https://***.ai/***.`,
		},
		{
			name:  "link inside Chinese text",
			input: `无法识别模型 ID“x”。请前往 https://relaysite.ai/models 复制正确的模型 ID 后重试。`,
			want:  `无法识别模型 ID“x”。请前往 https://***.ai/*** 复制正确的模型 ID 后重试。`,
		},
		{
			name:  "site name next to Chinese characters",
			input: `请联系Relaysite客服`,
			want:  `请联系***客服`,
		},
		{
			name:  "plain domains and subdomains",
			input: `see docs.relaysite.ai or relaysite.ai, x-relaysite-request-id`,
			want:  `see ***.***.ai or ***.ai, x-***-request-id`,
		},
		{
			name:  "other sites and longer words stay",
			input: `https://docs.example.org/errors relaysites`,
			want:  `https://docs.example.org/errors relaysites`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, maskUpstreamSites(tt.input, hosts))
		})
	}
}

func TestMaskUpstreamSitesSkipsShortNamesAndIPs(t *testing.T) {
	input := `api.abc.ai said abc; 10.0.0.1 failed`
	assert.Equal(t, input, maskUpstreamSites(input, []string{"api.abc.ai", "10.0.0.1"}))
}

func TestMaskSensitiveInfoHidesUpstreamSiteOnlyWhenEnabled(t *testing.T) {
	input := `Relaysite: get a key at https://relaysite.ai/api`
	upstream := "https://api.relaysite.ai/v1"

	assert.Equal(t, input, MaskSensitiveInfo(input, upstream))

	HideUpstreamIdentity.Store(true)
	t.Cleanup(func() { HideUpstreamIdentity.Store(false) })
	assert.Equal(t, `***: get a key at https://***.ai/***`, MaskSensitiveInfo(input, upstream))
}
