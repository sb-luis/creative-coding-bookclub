package utils

import (
	"regexp"
)

// InnerMarkToHTML converts a very simplified version of Markdown to HTML.
// It only supports links (written like [text](href)), bold (written like **text**), and italic (written like _text_).
func InnerMarkToHTML(markdown string) string {
	if markdown == "" {
		return ""
	}

	// Convert links [text](href) to <a href="href">text</a>
	markdown = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`).ReplaceAllString(markdown, `<a href="$2" target="_blank" class="ccb-link">$1</a>`)

	// Convert bold **text** to <b>text</b>
	markdown = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(markdown, `<b>$1</b>`)

	// Convert italic _text_ to <i>text</i>
	markdown = regexp.MustCompile(`_(.+?)_`).ReplaceAllString(markdown, `<i>$1</i>`)

	return markdown
}
