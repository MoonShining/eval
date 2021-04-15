package eval

import (
	"encoding/json"
	"errors"
	"github.com/MoonShining/eval/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func GetConditionAst(cond string)(*Condition,error) {
	ins := antlr.NewInputStream(cond)
	lexer := parser.Newcond_parserLexer(ins)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.Newcond_parserParser(stream)
	errLis := &ErrListener{}
	p.AddErrorListener(errLis)
	expr := p.Expr()
	if err := errLis.Err(); err != nil {
		return nil, err
	}

	vis := &CondVisiter{}
	expr.Accept(vis)
	return vis.Get(), nil
}

// parser tree visiter
// generate condition
type CondVisiter struct {
	stack []interface{}
	*parser.Basecond_parserVisitor
}

func (v *CondVisiter) Get() *Condition {
	ret := v.stack[0]
	v.stack = nil
	return ret.(*Condition)
}

func (v *CondVisiter) VisitLogicExpr(ctx *parser.LogicExprContext) interface{} {
	ctx.Expr(0).Accept(v)
	ctx.Expr(1).Accept(v)

	r, l := v.stack[len(v.stack)-1].(*Condition), v.stack[len(v.stack)-2].(*Condition)
	v.stack = v.stack[:len(v.stack)-2]

	cond := &Condition{}

	op := ctx.GetOp().GetText()
	if op == "&&" {
		cond.AND = append(cond.AND, l, r)
	} else {
		cond.OR = append(cond.OR, l, r)
	}

	v.stack = append(v.stack, cond)
	return nil
}

func (v *CondVisiter) VisitParenExpr(ctx *parser.ParenExprContext) interface{} {
	return ctx.Expr().Accept(v)
}

func (v *CondVisiter) VisitBinaryExpr(ctx *parser.BinaryExprContext) interface{} {
	return ctx.Binary().Accept(v)
}

func (v *CondVisiter) VisitCompareBinary(ctx *parser.CompareBinaryContext) interface{} {
	id := ctx.ID().GetText()
	op := ctx.GetOp().GetText()

	ctx.Atom().Accept(v)
	atomVal := v.stack[len(v.stack)-1]
	v.stack = v.stack[:len(v.stack)-1]

	v.stack = append(v.stack, &Condition{Left: id, Right: atomVal, Op: op})

	return nil
}

func (v *CondVisiter) VisitAtom(ctx *parser.AtomContext) interface{} {
	str := ctx.STRING()
	if str != nil {
		text := str.GetText()
		v.stack = append(v.stack, text[1:len(text)-1])
	}

	null := ctx.NULL()
	if null != nil {
		v.stack = append(v.stack, nil)
	}

	num := ctx.NUMBER()
	if num != nil {
		jsonNum := json.Number(num.GetText())
		v.stack = append(v.stack, jsonNum)
	}

	return nil
}


// syntax error handler
type ErrListener struct {
	err error
}

func (d *ErrListener) Err() error {
	return d.err
}

func (d *ErrListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	d.err = errors.New("SyntaxError")
}

func (d *ErrListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	d.err = errors.New("ambiguity")
}

func (d *ErrListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	d.err = errors.New("AttemptingFulContext")
}

func (d *ErrListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	d.err = errors.New("ContextSensitivity")
}

