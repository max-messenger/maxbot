package formatter

import (
	"testing"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestRestoreFormatting(t *testing.T) {
	markups := []model.Markup{
		{From: 0, Length: 5, Type: "strong"},
	}

	tests := []struct {
		name    string
		text    string
		markups []model.Markup
		format  model.TextFormat
		want    string
		wantErr bool
	}{
		{
			name:    "html format",
			text:    "hello",
			markups: markups,
			format:  model.FormatHTML,
			want:    "<b>hello</b>",
		},
		{
			name:    "markdown format",
			text:    "hello",
			markups: markups,
			format:  model.FormatMarkdown,
			want:    "**hello**",
		},
		{
			name:    "empty text",
			text:    "",
			markups: nil,
			format:  model.FormatHTML,
			want:    "",
		},
		{
			name:    "unsupported format",
			text:    "hello",
			markups: markups,
			format:  model.TextFormat("unknown"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RestoreFormatting(tt.text, tt.markups, tt.format)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMarkupToHTML(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		markups []model.Markup
		want    string
	}{
		{
			name: "no markups",
			text: "hello",
			want: "hello",
		},
		{
			name:    "empty text",
			text:    "",
			markups: []model.Markup{{From: 0, Length: 0, Type: "strong"}},
			want:    "",
		},
		{
			name:    "strong",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "strong"}},
			want:    "<b>hello</b>",
		},
		{
			name:    "emphasized",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "emphasized"}},
			want:    "<i>hello</i>",
		},
		{
			name:    "strikethrough",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "strikethrough"}},
			want:    "<s>hello</s>",
		},
		{
			name:    "underline",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "underline"}},
			want:    "<u>hello</u>",
		},
		{
			name:    "monospaced",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "monospaced"}},
			want:    "<code>hello</code>",
		},
		{
			name:    "highlighted",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "highlighted"}},
			want:    "<mark>hello</mark>",
		},
		{
			name:    "heading",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "heading"}},
			want:    "<h3>hello</h3>",
		},
		{
			name:    "quote",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "quote"}},
			want:    "<blockquote>hello</blockquote>",
		},
		{
			name:    "link",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "link", URL: "https://example.com"}},
			want:    `<a href="https://example.com">hello</a>`,
		},
		{
			name:    "link with special chars in URL",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "link", URL: `https://x.com/?a="b"&c<d>`}},
			want:    `<a href="https://x.com/?a=&#34;b&#34;&amp;c&lt;d&gt;">hello</a>`,
		},
		{
			name:    "user mention by id",
			text:    "@user",
			markups: []model.Markup{{From: 0, Length: 5, Type: "user_mention", UserID: 123}},
			want:    `<a href="max://user/123">@user</a>`,
		},
		{
			name:    "user mention by link",
			text:    "@user",
			markups: []model.Markup{{From: 0, Length: 5, Type: "user_mention", UserLink: "max://u/x"}},
			want:    `<a href="max://u/x">@user</a>`,
		},
		{
			name:    "user mention id takes precedence over link",
			text:    "@user",
			markups: []model.Markup{{From: 0, Length: 5, Type: "user_mention", UserID: 42, UserLink: "max://u/x"}},
			want:    `<a href="max://user/42">@user</a>`,
		},
		{
			name:    "user mention without id or link is skipped",
			text:    "@user",
			markups: []model.Markup{{From: 0, Length: 5, Type: "user_mention"}},
			want:    "@user",
		},
		{
			name:    "unsupported markup type is skipped",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "unknown"}},
			want:    "hello",
		},
		{
			name: "html escapes text",
			text: `a<b>&"c'`,
			want: "a&lt;b&gt;&amp;&#34;c&#39;",
		},
		{
			name:    "html escapes text inside markup",
			text:    "<b>",
			markups: []model.Markup{{From: 0, Length: 3, Type: "strong"}},
			want:    "<b>&lt;b&gt;</b>",
		},
		{
			name: "nested markup proper order",
			text: "hello world",
			markups: []model.Markup{
				{From: 0, Length: 11, Type: "emphasized"},
				{From: 0, Length: 5, Type: "strong"},
			},
			want: "<i><b>hello</b> world</i>",
		},
		{
			name: "nested markup partial overlap",
			text: "hello world",
			markups: []model.Markup{
				{From: 0, Length: 11, Type: "emphasized"},
				{From: 6, Length: 5, Type: "strong"},
			},
			want: "<i>hello <b>world</b></i>",
		},
		{
			name: "identical ranges",
			text: "abc",
			markups: []model.Markup{
				{From: 0, Length: 3, Type: "strong"},
				{From: 0, Length: 3, Type: "emphasized"},
			},
			want: "<b><i>abc</i></b>",
		},
		{
			name: "adjacent markups: close before open",
			text: "ab",
			markups: []model.Markup{
				{From: 0, Length: 1, Type: "strong"},
				{From: 1, Length: 1, Type: "emphasized"},
			},
			want: "<b>a</b><i>b</i>",
		},
		{
			name: "multiple disjoint markups",
			text: "abcdef",
			markups: []model.Markup{
				{From: 0, Length: 2, Type: "strong"},
				{From: 3, Length: 2, Type: "emphasized"},
			},
			want: "<b>ab</b>c<i>de</i>f",
		},
		{
			name:    "multi-byte text",
			text:    "Привет 🔥",
			markups: []model.Markup{{From: 0, Length: 8, Type: "strong"}}, // ← было 7
			want:    "<b>Привет 🔥</b>",
		},
		{
			name:    "invalid: negative from is skipped",
			text:    "hello",
			markups: []model.Markup{{From: -1, Length: 5, Type: "strong"}},
			want:    "hello",
		},
		{
			name:    "invalid: from out of range is skipped",
			text:    "hello",
			markups: []model.Markup{{From: 10, Length: 5, Type: "strong"}},
			want:    "hello",
		},
		{
			name:    "invalid: zero length is skipped",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 0, Type: "strong"}},
			want:    "hello",
		},
		{
			name:    "invalid: negative length is skipped",
			text:    "hello",
			markups: []model.Markup{{From: 2, Length: -1, Type: "strong"}},
			want:    "hello",
		},
		{
			name:    "end beyond text is clamped",
			text:    "hi",
			markups: []model.Markup{{From: 0, Length: 100, Type: "strong"}},
			want:    "<b>hi</b>",
		},
		{
			name: "valid markup still applied when another is invalid",
			text: "hello",
			markups: []model.Markup{
				{From: 100, Length: 1, Type: "strong"},
				{From: 0, Length: 5, Type: "emphasized"},
			},
			want: "<i>hello</i>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MarkupToHTML(tt.text, tt.markups)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMarkupToMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		markups []model.Markup
		want    string
	}{
		{
			name: "no markups",
			text: "hello",
			want: "hello",
		},
		{
			name: "plain text escapes markdown specials",
			text: "a*b_c~d`e[f]g\\h",
			want: `a\*b\_c\~d\` + "`" + `e\[f\]g\\h`,
		},
		{
			name:    "strong",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "strong"}},
			want:    "**hello**",
		},
		{
			name:    "emphasized",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "emphasized"}},
			want:    "*hello*",
		},
		{
			name:    "strikethrough",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "strikethrough"}},
			want:    "~~hello~~",
		},
		{
			name:    "underline",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "underline"}},
			want:    "++hello++",
		},
		{
			name:    "monospaced",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "monospaced"}},
			want:    "`hello`",
		},
		{
			name:    "highlighted",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "highlighted"}},
			want:    "^^hello^^",
		},
		{
			name:    "heading",
			text:    "title",
			markups: []model.Markup{{From: 0, Length: 5, Type: "heading"}},
			want:    "# title",
		},
		{
			name:    "quote",
			text:    "some quote",
			markups: []model.Markup{{From: 0, Length: 10, Type: "quote"}},
			want:    "> some quote",
		},
		{
			name:    "link",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "link", URL: "https://example.com"}},
			want:    "[hello](https://example.com)",
		},
		{
			name:    "link with ) in URL",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "link", URL: "https://x.com/a)b"}},
			want:    `[hello](https://x.com/a\)b)`,
		},
		{
			name:    "link with backslash in URL",
			text:    "hello",
			markups: []model.Markup{{From: 0, Length: 5, Type: "link", URL: `https://x.com/a\b`}},
			want:    `[hello](https://x.com/a\\b)`,
		},
		{
			name:    "user mention by id",
			text:    "@user",
			markups: []model.Markup{{From: 0, Length: 5, Type: "user_mention", UserID: 123}},
			want:    "[@user](max://user/123)",
		},
		{
			name:    "user mention by link",
			text:    "@user",
			markups: []model.Markup{{From: 0, Length: 5, Type: "user_mention", UserLink: "max://u/x"}},
			want:    "[@user](max://u/x)",
		},
		{
			name:    "user mention without id or link is skipped",
			text:    "@user",
			markups: []model.Markup{{From: 0, Length: 5, Type: "user_mention"}},
			want:    "@user",
		},
		{
			name:    "escapes inside markup",
			text:    "a*b",
			markups: []model.Markup{{From: 0, Length: 3, Type: "strong"}},
			want:    `**a\*b**`,
		},
		{
			name: "nested markup proper order",
			text: "hello world",
			markups: []model.Markup{
				{From: 0, Length: 11, Type: "emphasized"},
				{From: 0, Length: 5, Type: "strong"},
			},
			want: "***hello** world*",
		},
		{
			name:    "multi-byte text",
			text:    "Привет",
			markups: []model.Markup{{From: 0, Length: 6, Type: "strong"}},
			want:    "**Привет**",
		},
		{
			name:    "invalid range is skipped",
			text:    "hi",
			markups: []model.Markup{{From: 5, Length: 2, Type: "strong"}},
			want:    "hi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MarkupToMarkdown(tt.text, tt.markups)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEscapeMarkdown(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"hello", "hello"},
		{`a*b`, `a\*b`},
		{`a_b`, `a\_b`},
		{`a~b`, `a\~b`},
		{"a`b", "a\\`b"},
		{`a[b]c`, `a\[b\]c`},
		{`a\b`, `a\\b`},
		{`*_~` + "`" + `[]\`, `\*\_\~\` + "`" + `\[\]\\`},
		// Not escaped by design:
		{"(test)", "(test)"},
		{"# heading", "# heading"},
		{"> quote", "> quote"},
		{"+underline+", "+underline+"},
		{"!bang", "!bang"},
	}

	for _, tt := range tests {
		if got := escapeMarkdown(tt.in); got != tt.want {
			t.Errorf("escapeMarkdown(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEscapeMarkdownURL(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"https://x.com", "https://x.com"},
		{`https://x.com/a)b`, `https://x.com/a\)b`},
		{`a\b`, `a\\b`},
		{`a\)b`, `a\\\)b`},
		{`)))`, `\)\)\)`},
		{`\\\\`, `\\\\\\\\`},
	}

	for _, tt := range tests {
		if got := escapeMarkdownURL(tt.in); got != tt.want {
			t.Errorf("escapeMarkdownURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
