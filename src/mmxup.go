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
	Index 	string
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
	replacer := strings.NewReplacer(" ", "_", "'", "_", "\"", "_")
	return strings.ToLower(replacer.Replace(name))

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
			doc.Index = line[HEAD_LEN:]
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
			pattern:   regexp.MustCompile(`(?m)^(#+)[^!](.*)$`),
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
			pattern:   regexp.MustCompile(`(?m)\x60{3}\n([\s\S])+?\n\x60{3}`),
			processor: createCodeBlock,
		},
		// inline code
		// matches inline text surrounded by ``
		Rule{
			pattern:   regexp.MustCompile(`\x60(.*?)\x60`),
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
			pattern:   regexp.MustCompile(`(?m)(?:\W|^)_(.*?)_(?:\W|$)`),
			processor: createEmphasis,
		},
		// strike
		Rule{
			pattern:   regexp.MustCompile(`\~(.*)\~`),
			processor: createStrike,
		},
		// make ordered or unordered list
		Rule{
			pattern:   regexp.MustCompile(`(?m)^(\s*[-+]\s.*(\n|$))+`),
			processor: createList,
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
		// is important that this come last, so that it doesn't wrap
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
	for _, skipTag := range []string{"<hr/>", "<div>", "<pre>", "<blockquote>", "<ul>", "<ol>", "<dl>", "<h"} {
		if strings.HasPrefix(text, skipTag) {
		return body
	}
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
	code := match[0][4 : len(match[0])-4]
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
	quote = strings.Replace(quote, "\n>>>", "</p>", 1)
	quote = strings.Replace(quote, "\n", "</p><p>", numPs-1)

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
	} else if module == "twitter" {
		id := strings.TrimSpace(args[1])
		embedHtml = fmt.Sprintf("<blockquote class='twitter-tweet'><a href='https://twitter.com/x/status/%s'></a></blockquote> <script async src='https://platform.twitter.com/widgets.js' charset='utf-8'></script>", id)
	} else {
		panic(fmt.Sprintf("Unsupported module '%s' in embed", module))
	}
	return strings.Replace(body, match[0], embedHtml, 1)
}

func createLink(match []string, body string) string {
	args := strings.Split(match[1], ",")
	href := args[0]
	text := href

	// we use the character codes for {} so that subsequent find-and-replaces do not collide
	template := "<a href='%s'>&#123;%s&#125;</a>"
	if strings.HasPrefix(href, "http") || strings.HasPrefix(href, "#") {
		template = "<a href='%s' target='_blank'>&#123;^%s&#125;</a>"
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

	// positional args
	// 0. src
	// 1. alt
	// 3. style

	alt := ""
	if len(args) > 1 {
		alt = strings.TrimSpace(args[1])
	}

	style := ""
	if len(args) > 2 {
		style = strings.TrimSpace(args[2])
	}

	html := fmt.Sprintf("<img src='%s' alt='%s' style='%s'/>", src, alt, style)
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
	html := fmt.Sprintf("<s>%s</s>", text)
	return strings.Replace(body, match[0], html, 1)
}

func getListType (c string) string {
	if (c == "-") {
		return "ul"
	} else if (c == "+") {
		return "ol"
	}
	return "";
}

func createList(match []string, body string) string {
	text := strings.TrimSpace(match[0])
	items := strings.Split(text, "\n")
	html := ""
	listType := ""
	level := 0
	outerListType := "ul"
	if text[0] == '+' {
		outerListType = "ol"
	}

	for _, item := range items {
		if item == "" {
			continue
		}

		newLevel := len(item) - len(strings.TrimLeft(item, " "))
		newListType := getListType(item[newLevel:newLevel + 1])

		for i := 0; i < abs(newLevel - level) / 2; i++ {
			// for each two spaces of difference, open or close sublists
			if newLevel > level {
				// open the NEW list type
				html += fmt.Sprintf("<li class='sublist-container'><%s>", newListType)
			} else {
				// close the PREVIOUS list type
				html += fmt.Sprintf("</%s></li>", listType)
			}
		}
		listType = newListType
		level = newLevel

		html += fmt.Sprintf("<li>%s</li>", item[level + 2:])
	}

	if level != 0 {
		// we ended on an indented level, so close
		html += fmt.Sprintf("</%s></li>", listType)
	}

	html = fmt.Sprintf("<%s>%s</%s>", outerListType, html, outerListType)
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
