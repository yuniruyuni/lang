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

func TestParseExecute(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []*token.Token
		want    ast.AST
		wantErr bool
		invalid bool
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
			invalid: true,
		},
		{
			name: "add cannot consume term",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
				{Kind: kind.RightParen, Str: ")", Beg: 2, End: 3},
			},
			want:    nil,
			invalid: true,
		},
		{
			name: "sub cannot consume term",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Minus, Str: "-", Beg: 1, End: 2},
				{Kind: kind.RightParen, Str: ")", Beg: 2, End: 3},
			},
			want:    nil,
			invalid: true,
		},
		{
			name: "mul cannot consume term",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Multiply, Str: "*", Beg: 1, End: 2},
				{Kind: kind.RightParen, Str: ")", Beg: 2, End: 3},
			},
			want:    nil,
			invalid: true,
		},
		{
			name: "div cannot consume term",
			tokens: []*token.Token{
				{Kind: kind.Integer, Str: "1", Beg: 0, End: 1},
				{Kind: kind.Divide, Str: "/", Beg: 1, End: 2},
				{Kind: kind.RightParen, Str: ")", Beg: 2, End: 3},
			},
			want:    nil,
			invalid: true,
		},
		{
			name: "lack of clause right paren",
			tokens: []*token.Token{
				{Kind: kind.LeftParen, Str: "(", Beg: 0, End: 1},
				{Kind: kind.Integer, Str: "1", Beg: 1, End: 2},
			},
			want:    nil,
			wantErr: true,
			invalid: true,
		},
		{
			name: "invalid token for clause",
			tokens: []*token.Token{
				{Kind: kind.LeftParen, Str: "(", Beg: 0, End: 1},
				{Kind: kind.Plus, Str: "+", Beg: 1, End: 2},
			},
			want:    nil,
			wantErr: true,
			invalid: true,
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
			invalid: true,
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
			invalid: true,
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
		{
			name: "x = 1",
			tokens: []*token.Token{
				{
					Kind: kind.Identifier,
					Str:  "x",
					Beg:  0, End: 2,
				},
				{
					Kind: kind.Equal,
					Str:  "=",
					Beg:  3, End: 4,
				},
				{
					Kind: kind.Integer,
					Str:  "1",
					Beg:  5, End: 6,
				},
			},
			want: &ast.Assign{
				LHS: &ast.Variable{VarName: "x"},
				RHS: &ast.Integer{Value: 1},
			},
			wantErr: false,
		},
		{
			name: "while 1 { 2 }",
			tokens: []*token.Token{
				{
					Kind: kind.While,
					Str:  "while",
					Beg:  0, End: 5,
				},
				{
					Kind: kind.Integer,
					Str:  "1",
					Beg:  6, End: 7,
				},
				{
					Kind: kind.LeftCurly,
					Str:  "{",
					Beg:  9, End: 10,
				},
				{
					Kind: kind.Integer,
					Str:  "2",
					Beg:  11, End: 12,
				},
				{
					Kind: kind.RightCurly,
					Str:  "}",
					Beg:  13, End: 14,
				},
			},
			want: &ast.While{
				Cond: &ast.Integer{Value: 1},
				Proc: &ast.Integer{Value: 2},
			},
			wantErr: false,
		},
		{
			name: `f(123,456,) parses into Call("f", Params(123,456))`,
			tokens: []*token.Token{
				{Kind: kind.Identifier, Str: "f", Beg: 0, End: 3},
				{Kind: kind.LeftParen, Str: "(", Beg: 0, End: 3},
				{Kind: kind.Integer, Str: "123", Beg: 0, End: 3},
				{Kind: kind.Comma, Str: ",", Beg: 0, End: 3},
				{Kind: kind.Integer, Str: "456", Beg: 0, End: 3},
				{Kind: kind.Comma, Str: ",", Beg: 0, End: 3},
				{Kind: kind.RightParen, Str: ")", Beg: 0, End: 3},
			},
			want: &ast.Call{
				FuncName: &ast.FuncName{FuncName: "f"},
				Args: &ast.Args{
					Values: []ast.AST{
						&ast.Integer{Value: 123},
						&ast.Integer{Value: 456},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parse.New(tt.tokens)
			nx, got, err := p.Execute(0)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.invalid && !p.End(nx) {
				t.Errorf("parser.Execute() doesn't consume all tokens")
			}

			if !tt.invalid {
				assert.DeepEqual(t, tt.want, got)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []*token.Token
		want    ast.AST
		wantErr bool
	}{
		{
			name: `func main(){1} parses into DefFunc("main", Integer(1))`,
			tokens: []*token.Token{
				{Kind: kind.Func, Str: "func", Beg: 0, End: 4},
				{Kind: kind.Identifier, Str: "main", Beg: 6, End: 7},
				{Kind: kind.LeftParen, Str: "(", Beg: 7, End: 8},
				{Kind: kind.RightParen, Str: ")", Beg: 8, End: 9},
				{Kind: kind.LeftCurly, Str: "{", Beg: 9, End: 10},
				{Kind: kind.Integer, Str: "1", Beg: 10, End: 11},
				{Kind: kind.RightCurly, Str: "}", Beg: 11, End: 12},
			},
			want: &ast.Definitions{
				Defs: []ast.AST{
					&ast.Func{
						FuncName: &ast.FuncName{FuncName: "main"},
						Params:   &ast.Params{Vars: []ast.AST{}},
						Execute:  &ast.Integer{Value: 1},
					},
				},
			},
		},
		{
			name: `func test(x,){x} can parse properly`,
			tokens: []*token.Token{
				{Kind: kind.Func, Str: "func", Beg: 0, End: 4},
				{Kind: kind.Identifier, Str: "main", Beg: 6, End: 7},
				{Kind: kind.LeftParen, Str: "(", Beg: 7, End: 8},
				{Kind: kind.Identifier, Str: "x", Beg: 6, End: 7},
				{Kind: kind.Comma, Str: ",", Beg: 6, End: 7},
				{Kind: kind.RightParen, Str: ")", Beg: 8, End: 9},
				{Kind: kind.LeftCurly, Str: "{", Beg: 9, End: 10},
				{Kind: kind.Identifier, Str: "x", Beg: 10, End: 11},
				{Kind: kind.RightCurly, Str: "}", Beg: 11, End: 12},
			},
			want: &ast.Definitions{
				Defs: []ast.AST{
					&ast.Func{
						FuncName: &ast.FuncName{FuncName: "main"},
						Params: &ast.Params{
							Vars: []ast.AST{
								&ast.Param{VarName: "x"},
							},
						},
						Execute: &ast.Variable{VarName: "x"},
					},
				},
			},
		},
		{
			name: `func f(x,){x} func g(x,){x} can parse properly`,
			tokens: []*token.Token{
				{Kind: kind.Func, Str: "func", Beg: 0, End: 4},
				{Kind: kind.Identifier, Str: "f", Beg: 6, End: 7},
				{Kind: kind.LeftParen, Str: "(", Beg: 7, End: 8},
				{Kind: kind.Identifier, Str: "x", Beg: 6, End: 7},
				{Kind: kind.Comma, Str: ",", Beg: 6, End: 7},
				{Kind: kind.RightParen, Str: ")", Beg: 8, End: 9},
				{Kind: kind.LeftCurly, Str: "{", Beg: 9, End: 10},
				{Kind: kind.Identifier, Str: "x", Beg: 10, End: 11},
				{Kind: kind.RightCurly, Str: "}", Beg: 11, End: 12},
				{Kind: kind.Func, Str: "func", Beg: 0, End: 4},
				{Kind: kind.Identifier, Str: "g", Beg: 6, End: 7},
				{Kind: kind.LeftParen, Str: "(", Beg: 7, End: 8},
				{Kind: kind.Identifier, Str: "x", Beg: 6, End: 7},
				{Kind: kind.Comma, Str: ",", Beg: 6, End: 7},
				{Kind: kind.RightParen, Str: ")", Beg: 8, End: 9},
				{Kind: kind.LeftCurly, Str: "{", Beg: 9, End: 10},
				{Kind: kind.Identifier, Str: "x", Beg: 10, End: 11},
				{Kind: kind.RightCurly, Str: "}", Beg: 11, End: 12},
			},
			want: &ast.Definitions{
				Defs: []ast.AST{
					&ast.Func{
						FuncName: &ast.FuncName{FuncName: "f"},
						Params: &ast.Params{
							Vars: []ast.AST{
								&ast.Param{VarName: "x"},
							},
						},
						Execute: &ast.Variable{VarName: "x"},
					},
					&ast.Func{
						FuncName: &ast.FuncName{FuncName: "g"},
						Params: &ast.Params{
							Vars: []ast.AST{
								&ast.Param{VarName: "x"},
							},
						},
						Execute: &ast.Variable{VarName: "x"},
					},
				},
			},
		},
		{
			name: `func main(){f("%d",1+1,)} will be parsed properly`,
			tokens: []*token.Token{
				{Kind: kind.Func, Str: "func", Beg: 0, End: 4},
				{Kind: kind.Identifier, Str: "main", Beg: 6, End: 7},
				{Kind: kind.LeftParen, Str: "(", Beg: 7, End: 8},
				{Kind: kind.RightParen, Str: ")", Beg: 8, End: 9},
				{Kind: kind.LeftCurly, Str: "{", Beg: 9, End: 10},
				{Kind: kind.Identifier, Str: "f", Beg: 10, End: 11},
				{Kind: kind.LeftParen, Str: `(`, Beg: 11, End: 12},
				{Kind: kind.String, Str: `"%d"`, Beg: 12, End: 15},
				{Kind: kind.Comma, Str: ",", Beg: 15, End: 16},
				{Kind: kind.Integer, Str: "1", Beg: 16, End: 17},
				{Kind: kind.Plus, Str: "+", Beg: 17, End: 18},
				{Kind: kind.Integer, Str: "1", Beg: 18, End: 19},
				{Kind: kind.Comma, Str: ",", Beg: 19, End: 20},
				{Kind: kind.RightParen, Str: `)`, Beg: 20, End: 21},
				{Kind: kind.RightCurly, Str: "}", Beg: 21, End: 22},
			},
			want: &ast.Definitions{
				Defs: []ast.AST{
					&ast.Func{
						FuncName: &ast.FuncName{FuncName: "main"},
						Params:   &ast.Params{Vars: []ast.AST{}},
						Execute: &ast.Call{
							FuncName: &ast.FuncName{FuncName: "f"},
							Args: &ast.Args{
								Values: []ast.AST{
									&ast.String{Word: "%d"},
									&ast.Add{
										LHS: &ast.Integer{Value: 1},
										RHS: &ast.Integer{Value: 1},
									},
								},
							},
						},
					},
				},
			},
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
