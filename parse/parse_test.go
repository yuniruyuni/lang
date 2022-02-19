package parse_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/yuniruyuni/lang/ast"
	"github.com/yuniruyuni/lang/parse"
	"github.com/yuniruyuni/lang/token"
	"github.com/yuniruyuni/lang/token/kind"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []*token.Token
		want    ast.AST
		wantErr bool
	}{
		{
			name: `"abc" parses into String(word:"abc")`,
			tokens: []*token.Token{
				{Kind: kind.String, Str: "abc", Beg: 1, End: 4},
			},
			want:    &ast.String{Word: "abc"},
			wantErr: false,
		},
		{
			name: "123+456 parses into Sum(lhs:1, rhs:1)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Plus, Str: "+", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
			want: &ast.Sum{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name:    "root node not enough tokens",
			tokens:  []*token.Token{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid source code",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Plus, Str: "+", Beg: 3, End: 4},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse.Parse(tt.tokens)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.DeepEqual(t, tt.want, got)
		})
	}
}
