package main

import "testing"

func checkResult(result string, expectation string, t *testing.T) {
	if result != expectation {
		t.Errorf("Unexpected result, got: %s, want: %s.", result, expectation)
	}
}

func TestApplyRulesHeaders(t *testing.T) {
	checkResult(applyRules("# foo"), "<h5>foo</h5>", t)
	checkResult(applyRules("## foo"), "<h6>foo</h6>", t)
}

func TestApplyRulesLink(t *testing.T) {
	checkResult(applyRules("{foo, bar}"), "<p><a href='foo'>{bar}</a></p>", t)
	checkResult(applyRules("{foo, bar} {bop}"), "<p><a href='foo'>{bar}</a> <a href='bop'>{bop}</a></p>", t)
	checkResult(applyRules("{https://foo, bar}"), "<p><a href='https://foo' target='_blank'>{^bar}</a></p>", t)
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
}
