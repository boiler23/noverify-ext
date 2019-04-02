package main

// PSR-2 check for proper static modifier placement

import (
	"log"

	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/stmt"
	"github.com/z7zmey/php-parser/walker"
)

func main() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
	linter.RegisterRootChecker(func(ctx *linter.RootContext) linter.RootChecker { return &block{ctx: ctx} })
	cmd.Main()
}

type block struct {
	linter.RootCheckerDefaults
	ctx *linter.RootContext
}

func (b *block) BeforeEnterNode(w walker.Walkable) {
	switch w.(type) {
	case *stmt.ClassMethod:
		b.checkModifiersOrder(w.(*stmt.ClassMethod).Modifiers)

	case *stmt.PropertyList:
		b.checkModifiersOrder(w.(*stmt.PropertyList).Modifiers)
	}
}

func (b* block) checkModifiersOrder(modifiers []node.Node) {
	var nStatic node.Node = nil

	for _, m := range modifiers {
		switch v := m.(*node.Identifier).Value; v {
		case "static":
			nStatic = m
		case "public", "private", "protected":
			if nStatic != nil {
				b.ctx.Report(nStatic, linter.LevelWarning, "PSR-2", "static MUST be declared after the visibility.")
				return
			}
		}
	}
}
