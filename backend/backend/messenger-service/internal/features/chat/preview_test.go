package chat

import (
	"strings"
	"testing"
)

func TestPreviewSnippet(t *testing.T) {
	t.Run("short plain text passes through", func(t *testing.T) {
		if got := PreviewSnippet("hello there"); got != "hello there" {
			t.Fatalf("PreviewSnippet() = %q, want %q", got, "hello there")
		}
	})

	t.Run("json payload passes through untouched", func(t *testing.T) {
		in := `{"from_display_name":"Alice","amount_wei":"100"}`
		if got := PreviewSnippet(in); got != in {
			t.Fatalf("PreviewSnippet() = %q, want raw JSON %q", got, in)
		}
	})

	t.Run("long plain text is truncated with ellipsis", func(t *testing.T) {
		in := strings.Repeat("a", previewSnippetLen+50)
		got := PreviewSnippet(in)
		wantRunes := previewSnippetLen + 1 // +1 for the ellipsis rune
		if n := len([]rune(got)); n != wantRunes {
			t.Fatalf("truncated length = %d runes, want %d", n, wantRunes)
		}
		if !strings.HasSuffix(got, "…") {
			t.Fatalf("expected ellipsis suffix, got %q", got)
		}
	})
}
