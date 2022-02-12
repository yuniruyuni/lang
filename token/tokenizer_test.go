package token_test

import (
	"testing"

	"github.com/yuniruyuni/lang/token"
	"github.com/yuniruyuni/lang/token/kind"
	"gotest.tools/assert"
)

func TestTokenizer_Tokenize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		want []*token.Token
	}{
		{
			name: "empty code",
			code: "",
			want: []*token.Token{},
		},
		{
			name: "single string",
			code: "\"test\"",
			want: []*token.Token{
				{Kind: kind.String, Str: "test", Beg: 1, End: 5},
			},
		},
		{
			name: "with space",
			code: " \"test\"",
			want: []*token.Token{
				{Kind: kind.String, Str: "test", Beg: 2, End: 6},
			},
		},
		{
			name: "with tab",
			code: "\t\"test\"",
			want: []*token.Token{
				{Kind: kind.String, Str: "test", Beg: 2, End: 6},
			},
		},
		{
			name: "with space right",
			code: "\"test\" ",
			want: []*token.Token{
				{Kind: kind.String, Str: "test", Beg: 1, End: 5},
			},
		},
		{
			name: "with tab right",
			code: "\"test\"\t",
			want: []*token.Token{
				{Kind: kind.String, Str: "test", Beg: 1, End: 5},
			},
		},
		{
			name: "single integer",
			code: "1",
			want: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
			},
		},
		{
			name: "single integer",
			code: "10",
			want: []*token.Token{
				{Kind: kind.Integer, Str: "10", Beg: 0, End: 2},
			},
		},
		{
			name: "quoted integer is a string",
			code: `"1"`,
			want: []*token.Token{
				{Kind: kind.String, Str: "1", Beg: 1, End: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tk := token.Tokenizer{}
			got := tk.Tokenize(tt.code)
			assert.DeepEqual(t, tt.want, got)
		})
	}
}
