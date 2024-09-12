package ast

import (
	"bytes"
	"strings"

	"github.com/eugene-whitaker/writing-an-interpreter-in-go/token"
)

type Node interface {
	TokenLexeme() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLexeme() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLexeme()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLexeme() string {
	return ls.Token.Lexeme
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLexeme())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLexeme() string {
	return rs.Token.Lexeme
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLexeme())
	out.WriteString(" ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // The first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLexeme() string {
	return es.Token.Lexeme
}

func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type BlockStatement struct {
	Token      token.Token // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLexeme() string {
	return bs.Token.Lexeme
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLexeme() string {
	return i.Token.Lexeme
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLexeme() string {
	return il.Token.Lexeme
}

func (il *IntegerLiteral) String() string {
	return il.Token.Lexeme
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) TokenLexeme() string {
	return bl.Token.Lexeme
}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Lexeme
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLexeme() string {
	return fl.Token.Lexeme
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	ps := []string{}
	for _, p := range fl.Parameters {
		ps = append(ps, p.String())
	}

	out.WriteString(fl.TokenLexeme())
	out.WriteString("(")
	out.WriteString(strings.Join(ps, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLexeme() string {
	return sl.Token.Lexeme
}

func (sl *StringLiteral) String() string {
	return sl.Token.Lexeme
}

type ArrayLiteral struct {
	Token    token.Token // The '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLexeme() string {
	return al.Token.Lexeme
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	es := []string{}
	for _, e := range al.Elements {
		es = append(es, e.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(es, ", "))
	out.WriteString("]")

	return out.String()
}

type HashLiteral struct {
	Token token.Token // The '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}
func (hl *HashLiteral) TokenLexeme() string {
	return hl.Token.Lexeme
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	ps := []string{}
	for k, v := range hl.Pairs {
		ps = append(ps, k.String()+":"+v.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(ps, ", "))
	out.WriteString("}")

	return out.String()
}

type PrefixExpression struct {
	Token    token.Token // The operator token, e.g. '!'
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLexeme() string {
	return pe.Token.Lexeme
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. '+'
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLexeme() string {
	return ie.Token.Lexeme
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token // The 'if' Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLexeme() string {
	return ie.Token.Lexeme
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.TokenLexeme())
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLexeme() string {
	return ce.Token.Lexeme
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	as := []string{}
	for _, a := range ce.Arguments {
		as = append(as, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(as, ", "))
	out.WriteString(")")

	return out.String()
}

type IndexExpression struct {
	Token  token.Token // The '[' token
	Struct Expression  // ArrayLiteral or HashLiteral
	Index  Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLexeme() string {
	return ie.Token.Lexeme
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Struct.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type MacroExpression struct {
	Token      token.Token // The 'macro' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (me *MacroExpression) expressionNode() {}
func (me *MacroExpression) TokenLexeme() string {
	return me.Token.Lexeme
}

func (me *MacroExpression) String() string {
	var out bytes.Buffer

	ps := []string{}
	for _, p := range me.Parameters {
		ps = append(ps, p.String())
	}

	out.WriteString(me.TokenLexeme())
	out.WriteString("(")
	out.WriteString(strings.Join(ps, ", "))
	out.WriteString(")")
	out.WriteString(me.Body.String())

	return out.String()
}
