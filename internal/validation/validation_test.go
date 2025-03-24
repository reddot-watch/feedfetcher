package validation

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateAndSanitizeHeadline(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		want      string
		wantErr   error
	}{
		{
			name:      "valid headline",
			input:     "  Test  headline  ",
			maxLength: 50,
			want:      "Test headline",
			wantErr:   nil,
		},
		{
			name:      "empty headline",
			input:     "",
			maxLength: 50,
			want:      "",
			wantErr:   ErrEmptyHeadline,
		},
		{
			name:      "headline too long",
			input:     "This is a very long headline that exceeds the maximum length",
			maxLength: 10,
			want:      "",
			wantErr:   ErrHeadlineTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAndSanitizeHeadline(tt.input, tt.maxLength)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValidateAndResolveURL(t *testing.T) {
	baseURL, err := url.Parse("https://example.com/feed")
	require.NoError(t, err)

	tests := []struct {
		name    string
		feedURL *url.URL
		rawURL  string
		want    string
		wantErr error
	}{
		{
			name:    "absolute url",
			feedURL: baseURL,
			rawURL:  "https://example.com/article",
			want:    "https://example.com/article",
			wantErr: nil,
		},
		{
			name:    "relative url",
			feedURL: baseURL,
			rawURL:  "/article",
			want:    "https://example.com/article",
			wantErr: nil,
		},
		{
			name:    "empty url not allowed",
			feedURL: baseURL,
			rawURL:  "",
			want:    "",
			wantErr: ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAndResolveURL(tt.feedURL, tt.rawURL)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
