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
			name: "single digit integer",
			code: "1",
			want: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
			},
		},
		{
			name: "multi digit integer",
			code: "10",
			want: []*token.Token{
				{Kind: kind.Integer, Str: "10", Beg: 0, End: 2},
			},
		},
		{
			name: "quoted digit string",
			code: `"1"`,
			want: []*token.Token{
				{Kind: kind.String, Str: "1", Beg: 1, End: 2},
			},
		},
		{
			name: "plus",
			code: `+`,
			want: []*token.Token{
				{Kind: kind.Plus, Str: "+", Beg: 0, End: 1},
			},
		},
		{
			name: "minus",
			code: `-`,
			want: []*token.Token{
				{Kind: kind.Minus, Str: "-", Beg: 0, End: 1},
			},
		},
		{
			name: "multiply",
			code: `*`,
			want: []*token.Token{
				{Kind: kind.Multiply, Str: "*", Beg: 0, End: 1},
			},
		},
		{
			name: "divide",
			code: `/`,
			want: []*token.Token{
				{Kind: kind.Divide, Str: "/", Beg: 0, End: 1},
			},
		},
		{
			name: "leftparen",
			code: `(`,
			want: []*token.Token{
				{Kind: kind.LeftParen, Str: "(", Beg: 0, End: 1},
			},
		},
		{
			name: "rightparen",
			code: `)`,
			want: []*token.Token{
				{Kind: kind.RightParen, Str: ")", Beg: 0, End: 1},
			},
		},
		{
			name: "quoted plus is a string",
			code: `"+"`,
			want: []*token.Token{
				{Kind: kind.String, Str: "+", Beg: 1, End: 2},
			},
		},
		{
			name: "expression plus",
			code: `1+2`,
			want: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
			},
		},
		{
			name: "expression minus",
			code: `1-2`,
			want: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Minus, Str: "-", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
			},
		},
		{
			name: "term multiply",
			code: `1*2`,
			want: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Multiply, Str: "*", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
			},
		},
		{
			name: "term divide",
			code: `1/2`,
			want: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Divide, Str: "/", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
			},
		},
		{
			name: "res",
			code: `(1+2)`,
			want: []*token.Token{
				{Kind: kind.LeftParen, Str: "(", Beg: 0, End: 1},
				{Kind: kind.Integer, Str: "1", Beg: 1, End: 2},
				{Kind: kind.Plus, Str: "+", Beg: 2, End: 3},
				{Kind: kind.Integer, Str: "2", Beg: 3, End: 4},
				{Kind: kind.RightParen, Str: ")", Beg: 4, End: 5},
			},
		},
		{
			name: "expression with space",
			code: "1 +\t2",
			want: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 2, End: 3},
				{Kind: kind.Integer, Str: "2", Beg: 4, End: 5},
			},
		},
		{
			name: "multidigit integer expression",
			code: "123+456",
			want: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Plus, Str: "+", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
		},
		{
			name: "string expression",
			code: `"abc"+"def"`,
			want: []*token.Token{
				{Kind: kind.String, Str: "abc", Beg: 1, End: 4},
				{Kind: kind.Plus, Str: "+", Beg: 5, End: 6},
				{Kind: kind.String, Str: "def", Beg: 7, End: 10},
			},
		},
		{
			name: "string after an integer",
			code: `123"def"`,
			want: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.String, Str: "def", Beg: 4, End: 7},
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

func Test_isDigit(t *testing.T) {
	tests := []struct {
		name string
		ch   rune
		want bool
	}{
		{name: "0", ch: '0', want: true},
		{name: "1", ch: '1', want: true},
		{name: "9", ch: '9', want: true},
		{name: "a", ch: 'a', want: false},
		{name: "z", ch: 'z', want: false},
		{name: "+", ch: '+', want: false},
		{name: `"`, ch: '"', want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.DeepEqual(t, tt.want, token.IsDigit(tt.ch))
		})
	}
}
