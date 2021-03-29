package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type MmxDoc struct {
	Name string
	Slug string
	Host string
	Bref string
	Body string
	Date time.Time
}

type Rule struct {
	pattern   *regexp.Regexp
	processor func([]string, string) string
}

const DOC_DELIM = "===="
const HEAD_LEN = 6

func parseDate(s string) time.Time {
	d, err := time.Parse("2006-01-02", s)
	check(err)

	return d
}

func parseFile(file *os.File) []MmxDoc {
	var docs []MmxDoc

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	pageAcc := ""
	for scanner.Scan() {
		line := scanner.Text()

		if line == DOC_DELIM {
			// flush pageAcc if it has content and start new page
			if len(pageAcc) > 0 {
				docs = append(docs, parsePage(pageAcc))
			}
			pageAcc = ""
		} else {
			pageAcc += line + "\n"
		}
	}
	return docs
}

func parsePage(text string) MmxDoc {
	scanner := bufio.NewScanner(strings.NewReader(text))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var doc MmxDoc
	cursor := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NAME: ") {
			doc.Name = line[HEAD_LEN:]
			if len(doc.Slug) == 0 {
				doc.Slug = strings.ToLower(strings.ReplaceAll(doc.Name, " ", "_"))
			}
		} else if strings.HasPrefix(line, "SLUG: ") {
			doc.Slug = line[HEAD_LEN:]
		} else if strings.HasPrefix(line, "DATE: ") {
			doc.Date = parseDate(line[HEAD_LEN:])
		} else if strings.HasPrefix(line, "HOST: ") {
			doc.Host = line[HEAD_LEN:]
		} else if strings.HasPrefix(line, "BREF: ") {
			doc.Bref = line[HEAD_LEN:]
		} else if line == "BODY:" {
			body := text[cursor+HEAD_LEN:]
			doc.Body = applyRules(body)
		}
		cursor += len(line) + 1 // add one for newline
	}

	if doc.Name == "" || doc.Slug == "" || doc.Host == "" {
		panic("Doc is missing a required field")
	}
	return doc
}

func applyRules(body string) string {
	var rules = []Rule{
		// headers
		Rule{
			pattern:   regexp.MustCompile(`(?m)^(#+)(.*)$`),
			processor: createTitle,
		},
		// link
		Rule{
			pattern:   regexp.MustCompile(`\{(.*?)\}`),
			processor: createLink,
		},
		// code fences
		Rule{
			pattern:   regexp.MustCompile(`\x60{3}\n(.*)\n\x60{3}`),
			processor: createCodeBlock,
		},
		// blockquote
		Rule{
			pattern:   regexp.MustCompile(`>{3}\n(.*)\n>{3}`),
			processor: createBlockquote,
		},
		// image
		Rule{
			pattern:   regexp.MustCompile(`\[(.*\.(?:png|jpg|gif), .*)\]`),
			processor: createImage,
		},
		// bold
		Rule{
			pattern:   regexp.MustCompile(`\*(.*)\*`),
			processor: createBold,
		},
		// emphasis
		Rule{
			pattern:   regexp.MustCompile(`\_(.*)\_`),
			processor: createEmphasis,
		},
		// strike
		Rule{
			pattern:   regexp.MustCompile(`\~(.*)\~`),
			processor: createStrike,
		},
		// unordered list
		Rule{
			pattern:   regexp.MustCompile(`(?m)^(-\s.*(\n|$))+`),
			processor: createUnorderedList,
		},
		// ordered list
		Rule{
			pattern:   regexp.MustCompile(`(?m)(\+\s.*(\n|$))+`),
			processor: createOrderedList,
		},
		// definition list
		Rule{
			pattern:   regexp.MustCompile(`(?m)(\*\s.*(\n|$))+`),
			processor: createDefinitionList,
		},
		// horizontal brek
		Rule{
			pattern:   regexp.MustCompile(`(?m)^\-\-\-$`),
			processor: createHr,
		},
	}

	for _, r := range rules {
		matches := r.pattern.FindAllStringSubmatch(body, -1)
		for _, match := range matches {
			body = r.processor(match, body)
		}
	}

	return body

}

func createTitle(match []string, body string) string {
	level := len(match[1])
	title := strings.TrimSpace(match[2])
	html := fmt.Sprintf("<h%d>%s</h%d>", level, title, level)
	return strings.Replace(body, match[0], html, 1)
}

func createCodeBlock(match []string, body string) string {
	code := strings.TrimSpace(match[1])
	html := fmt.Sprintf("<pre><code>%s</code></pre>", code)
	return strings.Replace(body, match[0], html, 1)
}

func createBlockquote(match []string, body string) string {
	quote := strings.TrimSpace(match[1])
	html := fmt.Sprintf("<blockquote>%s</blockquote>", quote)
	return strings.Replace(body, match[0], html, 1)
}

func createLink(match []string, body string) string {
	args := strings.Split(match[1], ",")
	href := args[0]
	text := href

	// local, but needs html
	if !strings.HasPrefix(href, "http") && !strings.HasSuffix(href, ".html") {
		href += ".html"
	}

	if len(args) > 1 {
		text = strings.TrimSpace(args[1])
	}

	html := fmt.Sprintf("<a href='%s'>%s</a>", href, text)
	return strings.Replace(body, match[0], html, 1)
}

func createImage(match []string, body string) string {
	args := strings.Split(match[1], ",")
	src := args[0]
	alt := strings.TrimSpace(args[1])
	html := fmt.Sprintf("<img src='%s' alt='%s'/>", src, alt)
	return strings.Replace(body, match[0], html, 1)
}

func createBold(match []string, body string) string {
	text := strings.TrimSpace(match[1])
	html := fmt.Sprintf("<strong>%s</strong>", text)
	return strings.Replace(body, match[0], html, 1)
}

func createEmphasis(match []string, body string) string {
	text := strings.TrimSpace(match[1])
	html := fmt.Sprintf("<em>%s</em>", text)
	return strings.Replace(body, match[0], html, 1)
}

func createStrike(match []string, body string) string {
	text := strings.TrimSpace(match[1])
	html := fmt.Sprintf("<del>%s</del>", text)
	return strings.Replace(body, match[0], html, 1)
}

func createUnorderedList(match []string, body string) string {
	text := strings.TrimSpace(match[0])
	lis := strings.Split(text, "\n")
	cleanedLis := ""
	for _, li := range lis {
		cleanedLis += fmt.Sprintf("<li>%s</li>", li[2:])
	}
	html := fmt.Sprintf("<ul>%s</ul>", cleanedLis)
	return strings.Replace(body, match[0], html, 1)
}

func createOrderedList(match []string, body string) string {
	text := strings.TrimSpace(match[0])
	lis := strings.Split(text, "\n")
	cleanedLis := ""
	for _, li := range lis {
		cleanedLis += fmt.Sprintf("<li>%s</li>", li[2:])
	}
	html := fmt.Sprintf("<ol>%s</ol>", cleanedLis)
	return strings.Replace(body, match[0], html, 1)
}

func createDefinitionList(match []string, body string) string {
	text := strings.TrimSpace(match[0])
	termDefs := strings.Split(text, "\n")
	cleanedTermDefs := ""
	for _, termDef := range termDefs {
		defParts := strings.Split(termDef, ":")
		term := strings.TrimSpace(defParts[0])
		def := strings.TrimSpace(defParts[1])
		cleanedTermDefs += fmt.Sprintf("<dt>%s</dt><dd>%s</dd>", term[2:], def)
	}
	html := fmt.Sprintf("<dl>%s</dl>", cleanedTermDefs)
	return strings.Replace(body, match[0], html, 1)
}

func createHr(match []string, body string) string {
	return strings.Replace(body, match[0], "<hr/>", 1)
}

func _test() {
	doc := parsePage(`
NAME: Home
HOST: Home
BREF: The personal wiki of Thomasorus
BODY:

# Thomasorus' garden
` + "```" +
		`
code(foo)
` + "```" +

		`
*bold test*
_emphasis test_ ~strike test~
This is my little space where I store things.

- one
- two
- three

foo bar

+ a
+ b
+ c

- {about.html, "About me"}

>>>
foo
>>>

## Self Tracking

- {time.html, "Time tracker"}
- {tracking.html, "About tracking myself"}

## Shared Knowledge

- {html-tips.html, "HTML Tips and tricks, everything I know about it"}

[img/test.jpg, a test image]

## My Tools

- {tools.html, "Philosophy"}
- {kaku.html, "Kaku 書く", Kaku is a markup language}
- {ronbun.html, "Ronbun 論文", Robun is a static site generator}
- {keyboards.html, "Keyboards", The keyboards I use}

* foo : bar
* zee : twee

`)
	fmt.Printf("\n\n%+v\n", doc)
}
