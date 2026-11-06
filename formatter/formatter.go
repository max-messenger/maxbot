package formatter

import (
	"fmt"
	"html"
	"sort"
	"strings"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type event struct {
	pos   int
	open  bool
	order int
	tag   string
}

// ToHTML converts MAX markup into HTML.
func ToHTML(m model.MessageBody) string {
	return render(m.Text, m.MarkUp, true)
}

// ToMarkdown converts MAX markup into MAX Markdown.
func ToMarkdown(m model.MessageBody) string {
	return render(m.Text, m.MarkUp, false)
}

func render(text string, markups []model.Markup, isHTML bool) string {
	runes := []rune(text)

	// offsets[i] — смещение i-й руны в UTF-16 code units.
	// offsets[len(runes)] — общая длина текста в UTF-16.
	offsets := make([]int, len(runes)+1)
	for i, r := range runes {
		width := 1
		if r >= 0x10000 { // суррогатная пара в UTF-16
			width = 2
		}
		offsets[i+1] = offsets[i] + width
	}
	total := offsets[len(runes)]

	var events []event

	for i, m := range markups {
		// MAX отдаёт позиции в UTF-16 code units.
		u16Start := int(m.From)
		u16End := u16Start + int(m.Length)

		// Protect against invalid ranges.
		if m.Length <= 0 || u16Start < 0 || u16Start >= total {
			continue
		}
		if u16End > total {
			u16End = total
		}

		// UTF-16 → индексы рун. Если граница попала в середину
		// суррогатной пары, диапазон расширяется до целой руны:
		// начало сдвигается назад, конец — вперёд.
		start := sort.Search(len(offsets), func(k int) bool {
			return offsets[k] > u16Start
		}) - 1
		end := sort.SearchInts(offsets, u16End)

		open, tail, ok := delimiters(m, isHTML)
		if !ok {
			continue
		}

		events = append(events,
			event{pos: start, open: true, order: i, tag: open},
			event{pos: end, open: false, order: i, tag: tail},
		)
	}

	// At equal positions:
	// close inner elements first, open outer elements first.
	sort.SliceStable(events, func(i, j int) bool {
		a, b := events[i], events[j]

		if a.pos != b.pos {
			return a.pos < b.pos
		}

		if a.open != b.open {
			return !a.open
		}

		if a.open {
			return a.order < b.order
		}
		return a.order > b.order
	})

	var result strings.Builder
	eventIndex := 0

	for pos := 0; pos <= len(runes); pos++ {
		for eventIndex < len(events) &&
			events[eventIndex].pos == pos {
			result.WriteString(events[eventIndex].tag)
			eventIndex++
		}

		if pos < len(runes) {
			char := string(runes[pos])

			if isHTML {
				result.WriteString(html.EscapeString(char))
			} else {
				result.WriteString(escapeMarkdown(char))
			}
		}
	}

	return result.String()
}

func delimiters(m model.Markup, isHTML bool) (string, string, bool) {
	if isHTML {
		switch m.Type {
		case "strong":
			return "<b>", "</b>", true
		case "emphasized":
			return "<i>", "</i>", true
		case "strikethrough":
			return "<s>", "</s>", true
		case "underline":
			return "<u>", "</u>", true
		case "monospaced":
			return "<code>", "</code>", true
		case "highlighted":
			return "<mark>", "</mark>", true
		case "heading":
			return "<h3>", "</h3>", true
		case "quote":
			return "<blockquote>", "</blockquote>", true
		case "link":
			return `<a href="` + html.EscapeString(m.URL) + `">`,
				"</a>", true
		case "user_mention":
			if m.UserID != 0 {
				return `<a href="max://user/` +
					fmt.Sprint(m.UserID) + `">`, "</a>", true
			}
			if m.UserLink != "" {
				return `<a href="` +
						html.EscapeString(m.UserLink) + `">`,
					"</a>", true
			}
		}
	} else {
		switch m.Type {
		case "strong":
			return "**", "**", true
		case "emphasized":
			return "*", "*", true
		case "strikethrough":
			return "~~", "~~", true
		case "underline":
			return "++", "++", true
		case "monospaced":
			return "`", "`", true
		case "highlighted":
			return "^^", "^^", true
		case "heading":
			return "# ", "", true
		case "quote":
			return "> ", "", true
		case "link":
			return "[", "](" + escapeMarkdownURL(m.URL) + ")", true
		case "user_mention":
			if m.UserID != 0 {
				return "[", "](max://user/" +
					fmt.Sprint(m.UserID) + ")", true
			}
			if m.UserLink != "" {
				return "[", "](" +
					escapeMarkdownURL(m.UserLink) + ")", true
			}
		}
	}

	return "", "", false
}

func escapeMarkdown(s string) string {
	var result strings.Builder

	for _, r := range s {
		switch r {
		case '\\', '*', '_', '~', '`', '[', ']':
			result.WriteRune('\\')
		}
		result.WriteRune(r)
	}

	return result.String()
}

func escapeMarkdownURL(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(s, "\\", "\\\\"),
		")", "\\)",
	)
}
