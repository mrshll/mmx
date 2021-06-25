package main

import (
	"fmt"
	"testing"
)

func checkResult(result string, expectation string, t *testing.T) {
	if result != expectation {
		t.Errorf("Unexpected result:\n%s\n%s", result, expectation)
	}
}

func TestApplyRulesHeaders(t *testing.T) {
	checkResult(applyRules("# foo"), "<h5>foo</h5>", t)
	checkResult(applyRules("## foo"), "<h6>foo</h6>", t)
}

func TestApplyRulesLink(t *testing.T) {
	checkResult(applyRules("{foo, bar}"), "<p><a href='foo'>&#123;bar&#125;</a></p>", t)
	checkResult(applyRules("{foo, bar} {bop}"), "<p><a href='foo'>&#123;bar&#125;</a> <a href='bop'>&#123;bop&#125;</a></p>", t)
	checkResult(applyRules("{https://foo, bar}"), "<p><a href='https://foo' target='_blank'>&#123;^bar&#125;</a></p>", t)
	checkResult(applyRules("# test ({foo})"), "<h5>test (<a href='foo'>&#123;foo&#125;</a>)</h5>", t)
	// link appears twice
	checkResult(applyRules("{Upstream Tech} ({Upstream Tech})"), "<p><a href='Upstream Tech'>&#123;Upstream Tech&#125;</a> (<a href='Upstream Tech'>&#123;Upstream Tech&#125;</a>)</p>", t)
}

func TestApplyRulesCode(t *testing.T) {

	codeblockTest := "```" +
		`
#!/bin/sh
bash build.sh
while inotifywait -qqre modify ./src ./links ./data; do
  bash build.sh
done
` + "```"
	codeblockTestExpectation := `<pre><code>#!/bin/sh
bash build.sh
while inotifywait -qqre modify ./src ./links ./data; do
  bash build.sh
done</code></pre>`

	checkResult(applyRules("```\nfoo()\n```"), "<pre><code>foo()</code></pre>", t)
	checkResult(applyRules(codeblockTest), codeblockTestExpectation, t)
}

func TestApplyRulesBlockquote(t *testing.T) {
	checkResult(applyRules(">>>\nquote\n>>>"), "<blockquote><p>quote</p></blockquote>", t)
	checkResult(applyRules(">>>\nline\nline\n>>>"), "<blockquote><p>line</p><p>line</p></blockquote>", t)
	checkResult(applyRules(">>>\nquote\n>>> citation"), "<blockquote><p>quote</p><cite>citation</cite></blockquote>", t)
}

func TestApplyRulesEm(t *testing.T) {
	checkResult(applyRules("_foo_"), "<p><em>foo</em></p>", t)
	checkResult(applyRules("_foo_ bar"), "<p><em>foo</em> bar</p>", t)
	checkResult(applyRules("<a href='foo_bar_biff>foo</a>"), "<p><a href='foo_bar_biff>foo</a></p>", t)
	long := "Unlike some of the most immersive works of fiction I read that reveal to me more about the human condition, this book also educated me to a period of history I have only experienced through works penned or framed by Americans: _The Things They Carried_, _Apocalpyse Now_, _Rescue Dawn_, _Full Metal Jacket_, _Platoon_, and the other cowboys-in-the-east movies."
	longExpected := "<p>Unlike some of the most immersive works of fiction I read that reveal to me more about the human condition, this book also educated me to a period of history I have only experienced through works penned or framed by Americans: <em>The Things They Carried</em>, <em>Apocalpyse Now</em>, <em>Rescue Dawn</em>, <em>Full Metal Jacket</em>, <em>Platoon</em>, and the other cowboys-in-the-east movies.</p>"
	checkResult(applyRules(long), longExpected, t)

	checkResult(applyRules("*A note about backlog issues*: We use the {#1-3-6, 1-3-6 document} as our way to track work that is longer term (six months in startup land is like 3 years in big company land). We try to keep GitLab issues _very fresh_, meaning issues that fall outside of the “this month” in the 1-3-6 document should likely {https://codeburst.io/life-without-bug-tracking-925668ed6842, not exist} in GitLab."), "<p><strong>A note about backlog issues</strong>: We use the <a href='#1-3-6' target='_blank'>&#123;^1-3-6 document&#125;</a> as our way to track work that is longer term (six months in startup land is like 3 years in big company land). We try to keep GitLab issues <em>very fresh</em>, meaning issues that fall outside of the “this month” in the 1-3-6 document should likely <a href='https://codeburst.io/life-without-bug-tracking-925668ed6842' target='_blank'>&#123;^not exist&#125;</a> in GitLab.</p>", t)

}

func TestApplyRulesList(t *testing.T) {
	listTest := `
- A
  - Aa
- B
  + B1
  + B2
- C
  - Ca
  - Cb
`
	// listResult := "<ul><li>A</li><li class='sublist-container'><ul><li>Aa</li></ul></li><li>B</li><li class='sublist-container'><ol><li>B1</li><li>B2</li></ol></li><li>C</li><li class='sublist-container'><ul><li>Ca</li><li>Cb</li></ul></li></ul>"
	listResult := "<ul><li>A<ul><li>Aa</li></ul></li><li>B<ol><li>B1</li><li>B2</li></ol></li><li>C<ul><li>Ca</li><li>Cb</li></ul></li></ul>"
	checkResult(applyRules(listTest), listResult, t)

}

func _treeToString(tree *MmxDocNode) string {
	val := fmt.Sprintf("%*d - %s", tree.Depth, tree.Depth, tree.Tag)
	for _, node := range tree.Children {
		val += "\n"
		val += _treeToString(node)
	}
	return val
}

func TestCreateBodyTree(t *testing.T) {
	testBody := "<a><b></b><c></c></a>"
	tree := createBodyTree(testBody)
	checkResult(_treeToString(tree), "1 - <a>\n 2 - <b>\n 2 - <c>", t)
	checkResult(tree.NodeContent, testBody, t)

	testBody = "<a><b/><c></c></a>"
	tree = createBodyTree(testBody)
	checkResult(_treeToString(tree), "1 - <a>\n 2 - <b/>\n 2 - <c>", t)
	checkResult(tree.NodeContent, testBody, t)
}

func TestFindMmxDocNode(t *testing.T) {
	testBody := "<a><b></b><c></c></a>"
	tree := createBodyTree(testBody)
	checkResult(findMmxDocNode(tree, "<c>").Tag, "<c>", t)

	if findMmxDocNode(tree, "<d>") != nil {
		panic("searching for a missing node should be nil")
	}
}
