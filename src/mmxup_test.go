package main

import "testing"

func checkResult(result string, expectation string, t *testing.T) {
	if result != expectation {
		t.Errorf("Unexpected result for headers, got: %s, want: %s.", result, expectation)
	}
}

func TestApplyRulesHeaders(t *testing.T) {
	checkResult(applyRules("# foo"), "<h1>foo</h1>", t)
	checkResult(applyRules("## foo"), "<h2>foo</h2>", t)
}

func TestApplyRulesLink(t *testing.T) {
	checkResult(applyRules("{foo, bar}"), "<a href='foo.html'>bar</a>", t)
	checkResult(applyRules("{foo, bar} {bop}"), "<a href='foo.html'>bar</a> <a href='bop.html'>bop</a>", t)
}

func TestApplyRulesCode(t *testing.T) {
	checkResult(applyRules("```\nfoo()\n```"), "<pre><code>foo()</code></pre>", t)
}
