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
	Name    string
	Slug    string
	Host    string
	Bref    string
	Body    string
	Date    time.Time
	IsIndex bool
}

type Rule struct {
	pattern   *regexp.Regexp
	processor func([]string, string) string
}

const DOC_DELIM = "===="
const HEAD_LEN = 6
const H_ADJUSTMENT = 4

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

func makeSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))

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
		if strings.HasPrefix(line, "name: ") {
			doc.Name = line[HEAD_LEN:]
			if len(doc.Slug) == 0 {
				doc.Slug = makeSlug(doc.Name)
			}
		} else if strings.HasPrefix(line, "slug: ") {
			doc.Slug = line[HEAD_LEN:]
		} else if strings.HasPrefix(line, "date: ") {
			doc.Date = parseDate(line[HEAD_LEN:])
		} else if strings.HasPrefix(line, "host: ") {
			doc.Host = line[HEAD_LEN:]
		} else if strings.HasPrefix(line, "bref: ") {
			doc.Bref = line[HEAD_LEN:]
		} else if strings.HasPrefix(line, "indx: ") {
			doc.IsIndex = line[HEAD_LEN:] == "true"
		} else if line == "body:" {
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
		// matches markdown-style headers that take the full line
		Rule{
			pattern:   regexp.MustCompile(`(?m)^(#+)(.*)$`),
			processor: createTitle,
		},
		// embed
		Rule{
			pattern:   regexp.MustCompile(`\{\^(.*?)\}`),
			processor: createEmbed,
		},
		// link
		// matches links in {url, display} format. if url is a doc name, it appends html to it
		Rule{
			pattern:   regexp.MustCompile(`\{(.*?)\}`),
			processor: createLink,
		},
		// code fences
		// matches multiline code blocks surrounded by ```
		Rule{
			pattern:   regexp.MustCompile(`\x60{3}\n(.*)\n\x60{3}`),
			processor: createCodeBlock,
		},
		// inline code
		// matches inline text surrounded by ``
		Rule{
			pattern:   regexp.MustCompile(`\x60(.*)\x60`),
			processor: createCode,
		},
		// blockquote
		// matches multiline qupte blocks surrounded by >>>
		Rule{
			pattern:   regexp.MustCompile(`(?m)>>>\n([\s\S])+?>>>(\s.*)?`),
			processor: createBlockquote,
		},
		// image
		// matches img urls [img.ext, alt]
		Rule{
			pattern:   regexp.MustCompile(`\[(.*\.(?:png|jpg|jpeg|gif)(, .*)?)\]`),
			processor: createImage,
		},
		// bold
		Rule{
			pattern:   regexp.MustCompile(`\*(.*)\*`),
			processor: createBold,
		},
		// emphasis
		Rule{
			pattern:   regexp.MustCompile(`(?m)(?:^| )_(.*)?_(?:$| |\.)`),
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
			pattern:   regexp.MustCompile(`(?m)^(\+\s.*(\n|$))+`),
			processor: createOrderedList,
		},
		// definition list
		Rule{
			pattern:   regexp.MustCompile(`(?m)^(\*\s.*(\n|$))+`),
			processor: createDefinitionList,
		},
		// horizontal break
		Rule{
			pattern:   regexp.MustCompile(`(?m)^\-\-\-$`),
			processor: createHr,
		},
		// paragraphs
		//
		// NOTE: it is important that this come last, so that it doesn't wrap
		// text that should be tranformed by one of the preceeding rules
		Rule{
			pattern:   regexp.MustCompile(`(?s)((?:[^\n][\n]?)+)`),
			processor: createParagraph,
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

func createParagraph(match []string, body string) string {
	text := strings.TrimSpace(match[1])
	if text[0] == '<' {
		return body
	}

	html := fmt.Sprintf("<p>%s</p>", text)
	return strings.Replace(body, match[0], html, 1)
}

func createTitle(match []string, body string) string {
	level := len(match[1]) + H_ADJUSTMENT
	title := strings.TrimSpace(match[2])
	html := fmt.Sprintf("<h%d>%s</h%d>", level, title, level)
	return strings.Replace(body, match[0], html, 1)
}

func createCodeBlock(match []string, body string) string {
	code := strings.TrimSpace(match[1])
	html := fmt.Sprintf("<pre><code>%s</code></pre>", code)
	return strings.Replace(body, match[0], html, 1)
}

func createCode(match []string, body string) string {
	text := strings.TrimSpace(match[1])
	html := fmt.Sprintf("<code>%s</code>", text)
	return strings.Replace(body, match[0], html, 1)
}

func createBlockquote(match []string, body string) string {
	quote := strings.Replace(match[0], ">>>\n", "<p>", 1)
	numPs := strings.Count(quote, "\n")
	quote = strings.Replace(quote, "\n", "</p><p>", numPs-1)
	quote = strings.Replace(quote, "\n>>>", "</p>", 1)

	// clean up empty paragraphs
	quote = strings.Replace(quote, "<p></p>", "", -1)

	if len(match) == 3 && match[2] != "" {
		// we have a citation
		citation := strings.TrimSpace(match[2])
		quote = strings.Replace(quote, citation, "", 1)
		quote = strings.TrimSpace(quote)
		quote += fmt.Sprintf("<cite>%s</cite>", citation)
	} else {
		quote = strings.TrimSpace(quote)
	}

	html := fmt.Sprintf("<blockquote>%s</blockquote>", quote)
	return strings.Replace(body, match[0], html, 1)
}

func createEmbed(match []string, body string) string {
	args := strings.Split(match[1], ",")
	module := args[0]

	var embedHtml string
	if module == "bandcamp" {
		id := strings.TrimSpace(args[1])
		embedHtml = fmt.Sprintf("<iframe style='border: 0; width: 400px; height: 300px;' src='https://bandcamp.com/EmbeddedPlayer/album=%s/size=large/bgcol=ffffff/artwork=small/transparent=true/' seamless></iframe>", id)
	} else if module == "youtube" {
		id := strings.TrimSpace(args[1])
		embedHtml = fmt.Sprintf("<iframe width='560' height='315' src='https://www.youtube-nocookie.com/embed/%s' frameborder='0' allow='accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture' allowfullscreen></iframe>", id)
	} else if module == "buildtime" {
		t := time.Now()
		embedHtml = fmt.Sprintf("<span>%s</span>", t.Format("2006-01-02 3:04PM MST"))
	} else {
		panic(fmt.Sprintf("Unsupported module '%s' in embed", module))
	}
	return strings.Replace(body, match[0], embedHtml, 1)
}

func createLink(match []string, body string) string {
	args := strings.Split(match[1], ",")
	href := args[0]
	text := href

	template := "<a href='%s'>%s</a>"
	if strings.HasPrefix(href, "http") || strings.HasPrefix(href, "#") {
		template = "<a href='%s' target='_blank'>%s</a>"
	}

	if len(args) > 1 {
		text = strings.TrimSpace(args[1])
	}

	html := fmt.Sprintf(template, href, text)
	return strings.Replace(body, match[0], html, 1)
}

func createImage(match []string, body string) string {
	args := strings.Split(match[1], ",")
	src := args[0]
	alt := ""
	if len(args) > 1 {
		alt = strings.TrimSpace(args[1])
	}
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
	return strings.Replace(body, "_"+match[1]+"_", html, 1)
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
		defParts := strings.Split(termDef, " : ")
		term := strings.TrimSpace(defParts[0])
		if term != "" {
			cleanedTermDefs += fmt.Sprintf("<dt>%s</dt>", term[1:])
		}

		if len(defParts) > 1 {
			def := strings.TrimSpace(defParts[1])
			if def != "" {
				cleanedTermDefs += fmt.Sprintf("<dd>%s</dd>", def)
			}
		}
	}
	html := fmt.Sprintf("<dl>%s</dl>", cleanedTermDefs)
	return strings.Replace(body, match[0], html, 1)
}

func createHr(match []string, body string) string {
	return strings.Replace(body, match[0], "<hr/>", 1)
}
