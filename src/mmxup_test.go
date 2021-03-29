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
	checkResult(applyRules("{foo, bar}"), "<a href='foo'>bar</a>", t)
	checkResult(applyRules("{foo, bar} {bop}"), "<a href='foo'>bar</a> <a href='bop'>bop</a>", t)
}

func TestApplyRulesCode(t *testing.T) {
	checkResult(applyRules("```\nfoo()\n```"), "<pre><code>foo()</code></pre>", t)
}

func TestApplyRulesBlockquote(t *testing.T) {
	checkResult(applyRules(">>>\nquote\n>>>"), "<blockquote><p>quote</p></blockquote>", t)
	checkResult(applyRules(">>>\nline\nline\n>>>"), "<blockquote><p>line</p><p>line</p></blockquote>", t)
	checkResult(applyRules(">>>\nquote\n>>> citation"), "<blockquote><p>quote</p><cite>citation</cite></blockquote>", t)
}

func TestApplyRulesEm(t *testing.T) {
	checkResult(applyRules("_foo_"), "<em>foo</em>", t)
	checkResult(applyRules("_foo_ bar"), "<em>foo</em> bar", t)
	checkResult(applyRules("<a href='foo_bar_biff>foo</a>"), "<a href='foo_bar_biff>foo</a>", t)
}
