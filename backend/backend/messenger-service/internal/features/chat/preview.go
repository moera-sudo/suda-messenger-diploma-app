package chat

import "strings"

// previewSnippetLen — max rune length of a reply/forward preview snippet.
const previewSnippetLen = 120

// PreviewSnippet returns a length-bounded snippet of a message's content for a
// reply / forward preview. It performs NO localization — the backend only ships
// structured/raw data and the client renders it in the right language:
//   - structured/system messages store a JSON payload (content starts with '{'):
//     passed through untouched so the client can parse it by message `type`;
//   - plain text is truncated to a sane length.
func PreviewSnippet(content string) string {
	if strings.HasPrefix(strings.TrimSpace(content), "{") {
		return content
	}
	runes := []rune(content)
	if len(runes) > previewSnippetLen {
		return string(runes[:previewSnippetLen]) + "…"
	}
	return content
}
