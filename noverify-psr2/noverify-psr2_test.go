package main

import (
	"log"
	"strings"
	"testing"

	"github.com/VKCOM/noverify/src/linter"
	"github.com/VKCOM/noverify/src/meta"
	"github.com/z7zmey/php-parser/node"
)

func init() {
	linter.RegisterBlockChecker(func(ctx *linter.BlockContext) linter.BlockChecker { return &block{ctx: ctx} })
	go linter.MemoryLimiterThread()
}

func testParse(t *testing.T, filename string, contents string) (rootNode node.Node, w *linter.RootWalker) {
	var err error
	rootNode, w, err = linter.ParseContents(filename, []byte(contents), "UTF-8", nil)
	if err != nil {
		t.Errorf("Could not parse %s: %s", filename, err.Error())
		t.Fail()
	}

	if !meta.IsIndexingComplete() {
		w.UpdateMetaInfo()
	}

	return rootNode, w
}

func singleFileReports(t *testing.T, contents string) []*linter.Report {
	meta.ResetInfo()

	testParse(t, `test.php`, contents)
	meta.SetIndexingComplete(true)
	_, w := testParse(t, `test.php`, contents)

	return w.GetReports()
}

func TestStaticAfterVisibility(t *testing.T) {
	reports := singleFileReports(t, `<?
		class Test {
			protected $prop1 = 0;
			public static $prop2 = 1.0;
			static private $prop3 = "improper property declaration";
			protected static $prop4, $prop5; 

			final private function properMethod1() {
			}

			protected static function properMethod2() {
			}

			static public function improperMethod1() {
			}

			final static private function improperMethod2() {
			}
		}`)

	for _, r := range reports {
		log.Printf("%s", r)
	}

	if len(reports) != 3 {
		t.Errorf("Unexpected number of reports: expected 3, got %d", len(reports))
		if len(reports) < 3 {
			t.FailNow()
		}
	}

	for _, r := range reports {
		text := r.String()

		if !strings.Contains(text, "static MUST be declared after the visibility.") {
			t.Errorf("Wrong report text: expected 'static MUST be declared after the visibility.', got '%s'", text)
		}
	}
}
