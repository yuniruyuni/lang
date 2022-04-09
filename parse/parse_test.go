package parse_test

import (
	"math"
	"strconv"
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
				{Kind: kind.String, Str: `"abc"`, Beg: 0, End: 5},
			},
			want:    &ast.String{Word: "abc"},
			wantErr: false,
		},
		{
			name: "123+456 parses into Add(lhs:123, rhs:456)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Plus, Str: "+", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
			want: &ast.Add{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name: "123-456 parses into Sub(lhs:123, rhs:456)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Minus, Str: "-", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
			want: &ast.Sub{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name: "123*456 parses into Mul(lhs:123, rhs:456)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Multiply, Str: "*", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
			want: &ast.Mul{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name: "123<456 parses into Less(lhs:123, rhs:456)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Less, Str: "<", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
			want: &ast.Less{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name: "123==456 parses into Equal(lhs:123, rhs:456)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Equal, Str: "=", Beg: 3, End: 4},
				{Kind: kind.Equal, Str: "=", Beg: 4, End: 5},
				{Kind: kind.Integer, Str: "456", Beg: 5, End: 8},
			},
			want: &ast.Equal{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name: "1+2+3 parses into Add(lhs:Add(lhs:1, rhs:2), rhs:3))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
				{Kind: kind.Plus, Str: "+", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "3", Beg: 4, End: 5},
			},
			want: &ast.Add{
				LHS: &ast.Integer{Value: 1},
				RHS: &ast.Add{
					LHS: &ast.Integer{Value: 2},
					RHS: &ast.Integer{Value: 3},
				},
			},
		},
		{
			name: "1+2-3 parses into Sub(lhs:Add(lhs:1, rhs:2), rhs:3))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
				{Kind: kind.Minus, Str: "-", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "3", Beg: 4, End: 5},
			},
			want: &ast.Add{
				LHS: &ast.Integer{Value: 1},
				RHS: &ast.Sub{
					LHS: &ast.Integer{Value: 2},
					RHS: &ast.Integer{Value: 3},
				},
			},
		},
		{
			name: "1-2+3 parses into Add(lhs:Sub(lhs:1, rhs:2), rhs:3))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Minus, Str: "-", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
				{Kind: kind.Plus, Str: "+", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "3", Beg: 4, End: 5},
			},
			want: &ast.Sub{
				LHS: &ast.Integer{Value: 1},
				RHS: &ast.Add{
					LHS: &ast.Integer{Value: 2},
					RHS: &ast.Integer{Value: 3},
				},
			},
		},
		{
			name: "3*2+1 parses into Add(lhs: Mul(lhs:3, rhs:2), rhs:1))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "3", Beg: 0, End: 1},
				{Kind: kind.Divide, Str: "*", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
				{Kind: kind.Plus, Str: "+", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "1", Beg: 4, End: 5},
			},
			want: &ast.Add{
				LHS: &ast.Div{
					LHS: &ast.Integer{Value: 3},
					RHS: &ast.Integer{Value: 2},
				},
				RHS: &ast.Integer{Value: 1},
			},
		},
		{
			name: "1+2*3 parses into Add(lhs: 1, rhs:Mul(lhs:2, rhs:3))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
				{Kind: kind.Multiply, Str: "*", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "3", Beg: 4, End: 5},
			},
			want: &ast.Add{
				LHS: &ast.Integer{Value: 1},
				RHS: &ast.Mul{
					LHS: &ast.Integer{Value: 2},
					RHS: &ast.Integer{Value: 3},
				},
			},
		},
		{
			name: "1*2*3 parses into Mul(lhs: Mul(lhs:1, rhs:2), rhs:3)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Multiply, Str: "*", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
				{Kind: kind.Multiply, Str: "*", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "3", Beg: 4, End: 5},
			},
			want: &ast.Mul{
				LHS: &ast.Integer{Value: 1},
				RHS: &ast.Mul{
					LHS: &ast.Integer{Value: 2},
					RHS: &ast.Integer{Value: 3},
				},
			},
		},
		{
			name: "1+6/2 parses into Add(lhs: 1, rhs:Div(lhs:6, rhs:2))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "6", Beg: 2, End: 3},
				{Kind: kind.Divide, Str: "/", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "2", Beg: 4, End: 5},
			},
			want: &ast.Add{
				LHS: &ast.Integer{Value: 1},
				RHS: &ast.Div{
					LHS: &ast.Integer{Value: 6},
					RHS: &ast.Integer{Value: 2},
				},
			},
		},
		{
			name: "1-6/2 parses into Sub(lhs: 1, rhs:Div(lhs:6, rhs:2))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Minus, Str: "-", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "6", Beg: 2, End: 3},
				{Kind: kind.Divide, Str: "/", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "2", Beg: 4, End: 5},
			},
			want: &ast.Sub{
				LHS: &ast.Integer{Value: 1},
				RHS: &ast.Div{
					LHS: &ast.Integer{Value: 6},
					RHS: &ast.Integer{Value: 2},
				},
			},
		},
		{
			name: "6/2+1 parses into Add(lhs:Div(lhs:6, rhs:2), rhs:1)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "6", Beg: 0, End: 1},
				{Kind: kind.Divide, Str: "/", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
				{Kind: kind.Plus, Str: "+", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "1", Beg: 4, End: 5},
			},
			want: &ast.Add{
				LHS: &ast.Div{
					LHS: &ast.Integer{Value: 6},
					RHS: &ast.Integer{Value: 2},
				},
				RHS: &ast.Integer{Value: 1},
			},
		},
		{
			name: "6/2-1 parses into Sub(lhs:Div(lhs:6, rhs:2), rhs:1)",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "6", Beg: 0, End: 1},
				{Kind: kind.Divide, Str: "/", Beg: 1, End: 2},
				{Kind: kind.Integer, Str: "2", Beg: 2, End: 3},
				{Kind: kind.Minus, Str: "-", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "1", Beg: 4, End: 5},
			},
			want: &ast.Sub{
				LHS: &ast.Div{
					LHS: &ast.Integer{Value: 6},
					RHS: &ast.Integer{Value: 2},
				},
				RHS: &ast.Integer{Value: 1},
			},
		},
		{
			name: "(123+456) parses into Add(lhs:123, rhs:456)",
			tokens: []*token.Token{
				{Kind: kind.LeftParen, Str: "(", Beg: 0, End: 1},
				{Kind: kind.Integer, Str: "123", Beg: 1, End: 2},
				{Kind: kind.Plus, Str: "+", Beg: 2, End: 3},
				{Kind: kind.Integer, Str: "456", Beg: 3, End: 4},
				{Kind: kind.RightParen, Str: ")", Beg: 4, End: 5},
			},
			want: &ast.Add{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name: "(123-456) parses into Res(Sub(lhs:123, rhs:456))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Minus, Str: "-", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
			want: &ast.Sub{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name: "(123*456) parses into Res(child:Mul(lhs:123, rhs:456))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Multiply, Str: "*", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
			want: &ast.Mul{
				LHS: &ast.Integer{Value: 123},
				RHS: &ast.Integer{Value: 456},
			},
		},
		{
			name: "(123/456) parses into Res(child:Div(lhs:123, rhs:456))",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Divide, Str: "/", Beg: 3, End: 4},
				{Kind: kind.Integer, Str: "456", Beg: 4, End: 7},
			},
			want: &ast.Div{
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
		{
			name: "add cannot consume term",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
				{Kind: kind.RightParen, Str: ")", Beg: 2, End: 3},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sub cannot consume term",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Minus, Str: "-", Beg: 1, End: 2},
				{Kind: kind.RightParen, Str: ")", Beg: 2, End: 3},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "mul cannot consume term",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Multiply, Str: "*", Beg: 1, End: 2},
				{Kind: kind.RightParen, Str: ")", Beg: 2, End: 3},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "div cannot consume term",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Divide, Str: "/", Beg: 1, End: 2},
				{Kind: kind.RightParen, Str: ")", Beg: 2, End: 3},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "lack of clause right paren",
			tokens: []*token.Token{
				{Kind: kind.LeftParen, Str: "(", Beg: 0, End: 1},
				{Kind: kind.Integer, Str: "1", Beg: 1, End: 2},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid token for clause",
			tokens: []*token.Token{
				{Kind: kind.LeftParen, Str: "(", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "over than integer max",
			tokens: []*token.Token{
				{
					Kind: kind.Integer,
					// add one digit to over MaxInt64
					Str: strconv.Itoa(math.MaxInt64) + "0",
					Beg: 1, End: 2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when just if",
			tokens: []*token.Token{
				{
					Kind: kind.If,
					Str:  "if",
					Beg:  0, End: 2,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "let x = 10",
			tokens: []*token.Token{
				{
					Kind: kind.Let,
					Str:  "let",
					Beg:  0, End: 3,
				},
				{
					Kind: kind.Identifier,
					Str:  "x",
					Beg:  4, End: 5,
				},
				{
					Kind: kind.Equal,
					Str:  "=",
					Beg:  4, End: 5,
				},
				{
					Kind: kind.Integer,
					Str:  "1",
					Beg:  6, End: 7,
				},
			},
			want: &ast.Let{
				LHS: &ast.Variable{VarName: "x"},
				RHS: &ast.Integer{Value: 1},
			},
			wantErr: false,
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
