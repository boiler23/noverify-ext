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
	linter.RegisterBlockChecker(func(ctx *linter.BlockContext) linter.BlockChecker { return &block{ctx: ctx} })
	cmd.Main()
}

type block struct {
	linter.BlockCheckerDefaults
	ctx *linter.BlockContext
}

func (b *block) BeforeEnterNode(w walker.Walkable) {
	switch w.(type) {
	case *stmt.Class:
		cls := w.(*stmt.Class)
		b.handleClassStatements(cls.Stmts)
	}
}

func (b *block) handleClassStatements(stmts []node.Node) {
	for _, s := range stmts {
		switch s.(type) {
		case *stmt.ClassMethod:
			method := s.(*stmt.ClassMethod)
			b.checkModifiersOrder(method, method.Modifiers)

		case *stmt.PropertyList:
			plist := s.(*stmt.PropertyList)
			b.checkModifiersOrder(plist, plist.Modifiers)
		}
	}
}

func (b* block) checkModifiersOrder(method node.Node, modifiers []node.Node) {
	visibilityModifiers := map[string]bool { "private": false, "protected": false, "public" : false }

	var mStatic node.Node = nil

	for _, m := range modifiers {
		mValue := m.(*node.Identifier).Value
		_, match := visibilityModifiers[mValue]
		if (match) {
			if mStatic != nil {
				b.ctx.Report(mStatic, linter.LevelWarning, "PSR-2", "static MUST be declared after the visibility.")
				return
			}
		} else if mValue == "static" {
			mStatic = m
		}
	}
}
