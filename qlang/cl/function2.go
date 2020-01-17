package qlang

import (
	"github.com/qiniu/text/tpl/interpreter.util"
	"github.com/jmniu/workflow/qlang/exec"
)

// -----------------------------------------------------------------------------

func (p *Compiler) function(e interpreter.Engine) {

	fnb, _ := p.gstk.Pop()
	variadic := p.popArity()
	arity := p.popArity()
	args := p.gstk.PopFnArgs(arity)
	instr := p.code.Reserve()
	fnctx := p.fnctx
	p.exits = append(p.exits, func() {
		start, end, symtbl := p.clFunc(e, "doc", fnb, fnctx, args)
		instr.Set(exec.Func(nil, start, end, symtbl, args, variadic != 0))
	})
}

func (p *Compiler) anonymFn(e interpreter.Engine) {

	fnb, _ := p.gstk.Pop()
	instr := p.code.Reserve()
	fnctx := p.fnctx
	p.exits = append(p.exits, func() {
		start, end, symtbl := p.clFunc(e, "doc", fnb, fnctx, nil)
		instr.Set(exec.AnonymFn(start, end, symtbl))
	})
}

func (p *Compiler) fnReturn(e interpreter.Engine) {

	arity := p.popArity()
	p.code.Block(exec.Return(arity))
}

// Done completes all exit functions generated by `Cl`.
//
func (p *Compiler) Done() {

	for {
		n := len(p.exits)
		if n == 0 {
			break
		}
		onExit := p.exits[n-1]
		p.exits = p.exits[:n-1]
		onExit()
	}
}

func (p *Compiler) fnDefer(e interpreter.Engine) {

	src, _ := p.gstk.Pop()
	instr := p.code.Reserve()
	fnctx := p.fnctx
	p.exits = append(p.exits, func() {
		start, end := p.clBlock(e, "expr", src, fnctx)
		p.codeLine(src)
		instr.Set(exec.Defer(start, end))
	})
}

func (p *Compiler) fnRecover() {

	p.code.Block(exec.Recover)
}

func (p *Compiler) clFunc(
	e interpreter.Engine, g string, src interface{},
	parent *funcCtx, args []string) (start, end int, symtbl map[string]int) {

	ctx := newFuncCtx(parent, args)
	old := p.fnctx
	p.fnctx = ctx
	start, end = p.cl(e, g, src)
	symtbl = ctx.symtbl
	p.fnctx = old
	return
}

func (p *Compiler) clBlock(
	e interpreter.Engine, g string, src interface{}, parent *funcCtx) (start, end int) {

	old := p.fnctx
	p.fnctx = parent
	start, end = p.cl(e, g, src)
	p.fnctx = old
	return
}

func (p *Compiler) cl(e interpreter.Engine, g string, src interface{}) (start, end int) {

	start = p.code.Len()
	if src != nil {
		if err := e.EvalCode(p, g, src); err != nil {
			panic(err)
		}
	}
	end = p.code.Len()
	return
}

// -----------------------------------------------------------------------------
