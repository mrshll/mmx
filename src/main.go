package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
)

type Entry struct {
	Name          string
	Filename      string
	Bref          string
	Host          string
	Parent        *Entry
	EmbedChildren bool
	Children      []*Entry
	Body          string
	Incoming      []*Entry
	Outgoing      []*Entry
}

type TemplateContent struct {
	Entry         Entry
	Children      []Entry
	NavHTMLString string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func spad(line string, c rune) int {
	i := 0
	for _, runeValue := range line {
		if runeValue == c {
			i++
		} else {
			break
		}
	}
	return i
}

func parseDate(s string) time.Time {
	d, err := time.Parse("2006-01-02", s)
	check(err)

	return d
}

func formatDate(d time.Time) string {
	return d.Format("2006-01-02")
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func parseIndentalLine(line string) (string, string) {
	delimiterIndex := strings.Index(line, ":")

	if delimiterIndex == -1 {
		panic(fmt.Sprintf("No delmiter found in Indental line: %s", line))
	}

	key := line[2 : delimiterIndex-1]
	var value string

	if delimiterIndex < len(line)-1 {
		value = line[delimiterIndex+2 : len(line)]
	}

	return key, value
}

func loadIndental(file *os.File) []Entry {

	var entries []Entry
	var catchBody bool

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		depth := spad(line, ' ') / 2
		lastEntryIndex := len(entries) - 1

		if depth == 0 && line != "" {
			catchBody = false
			name := line
			entries = append(entries, Entry{Name: name, Filename: strings.ToLower(strings.ReplaceAll(name, " ", "_"))})
		} else if depth == 1 && !catchBody {
			catchBody = false
			key, value := parseIndentalLine(line)
			// if key == "DATE" {
			// entries[lastEntryIndex].Date = parseDate(value)
			/*} else*/
			if key == "HOST" {
				entries[lastEntryIndex].Host = value
			} else if key == "BREF" {
				entries[lastEntryIndex].Bref = value
			} else if key == "EMBC" && value == "true" {
				entries[lastEntryIndex].EmbedChildren = true
			} else {
				catchBody = key == "BODY"
			}
		} else if depth >= 2 {
			if catchBody {
				fmt.Println(line)
				entries[lastEntryIndex].Body += line[4:] + "\n"
			}
		} else if line == "" && catchBody {
			entries[lastEntryIndex].Body += "\n"
		}
	}

	return entries
}

func loadJournal() []Entry {
	file, err := os.Open("../data/journal.ndtl")
	check(err)
	defer file.Close()

	return loadIndental(file)
}

func loadLex() []Entry {
	file, err := os.Open("../data/lex.ndtl")
	check(err)
	defer file.Close()

	return loadIndental(file)
}

func findEntry(entries []Entry, name string) *Entry {
	for i := range entries {
		if strings.ToLower(entries[i].Name) == strings.ToLower(name) {
			return &entries[i]
		}
	}
	panic(fmt.Sprintf("No parent found with name %s", name))
}

func processBody(e Entry, entries []Entry) string {
	refRegex := regexp.MustCompile(`{[^{}]*}`)
	b := e.Body

	matches := refRegex.FindAllString(b, -1)
	for _, match := range matches {
		cleanMatch := match[1 : len(match)-1]
		matchParts := strings.Split(cleanMatch, "|")

		isExternal := strings.Contains(cleanMatch, "http")

		var display string
		if len(matchParts) > 1 {
			display = strings.Join(matchParts[1:], " ")
		} else {
			display = matchParts[0]
		}

		var link string
		if isExternal {
			// external link
			link = fmt.Sprintf("<a href='%s' target='_blank'>[%s]</a>", matchParts[0], display)
		} else {
			refEntry := findEntry(entries, matchParts[0])
			if refEntry == nil {
				panic(fmt.Sprintf("No entry found with name %s", matchParts[0]))
			}
			e.Outgoing = append(e.Outgoing, refEntry)
			refEntry.Incoming = append(refEntry.Incoming, &e)

			link = fmt.Sprintf("<a href='./%s.html'>{%s}</a>", refEntry.Filename, display)
		}

		fmt.Println(link, display)
		b = strings.Replace(b, match, link, 1)
	}

	// convert markdown
	output := markdown.ToHTML([]byte(b), nil, nil)
	b = string(output)

	return b
}

func linkEntries(entries []Entry) {
	for i := range entries {
		parentPtr := findEntry(entries[:], entries[i].Host)
		entries[i].Body = processBody(entries[i], entries)
		entries[i].Parent = parentPtr
		(*parentPtr).Children = append((*parentPtr).Children, &(entries[i]))
	}
}

func makeSubNav(e Entry, target Entry) string {
	subnav := "<ul>"
	max := 8
	for i, cPtr := range e.Children {
		child := *cPtr
		if i >= max {
			if i == max {
				subnav += fmt.Sprintf("<li>and %d more</li>", len(e.Children)-max)
			}

			continue
		}

		if child.Name == e.Name {
			continue // this occurs in the case of root node, i.e. Home
		}

		if e.EmbedChildren {
			subnav += fmt.Sprintf("<li><a href='%s.html#%s'>%s/</a><mark></li>", e.Filename, child.Filename, child.Name)
		} else if child.Name == target.Name {
			subnav += fmt.Sprintf("<li><mark><a href='%s.html'>%s/</a><mark></li>", child.Filename, child.Name)
		} else {
			subnav += fmt.Sprintf("<li><a href='%s.html'>%s</a></li>", child.Filename, child.Name)
		}
	}

	subnav += "</ul>"
	return subnav
}

func makeNav(e Entry) string {
	nav := "<nav>"
	if e.Parent == nil {
		panic(fmt.Sprintf("No parent found with name %s", e.Name))
	}
	if e.Parent.Parent == nil {
		panic(fmt.Sprintf("No parent found with name %s", e.Parent.Name))
	}

	// this happens for our root node
	if e.Parent.Parent.Name == e.Parent.Name {
		nav += makeSubNav(*e.Parent.Parent, e)
	} else {
		nav += makeSubNav(*e.Parent.Parent, *e.Parent)
	}

	if e.Parent.Parent.Name != e.Parent.Name {
		nav += makeSubNav(*e.Parent, e)
	}

	if e.Parent.Name != e.Name {
		nav += makeSubNav(e, e)
	}

	nav += "</nav>"
	return nav
}

func renderEntryHTML(e Entry) string {
	templateFuncs := template.FuncMap{
		"noescape": noescape,
	}
	tmpl := template.Must(template.New("entry.html").Funcs(templateFuncs).ParseGlob("./templates/*.html"))
	embededChildTmpl := template.Must(template.New("embeddedChild.html").Funcs(templateFuncs).ParseFiles("./templates/embeddedChild.html", "./templates/incoming.html"))

	var children []Entry
	var embeddedHTMLStr string
	for _, cPtr := range e.Children {
		children = append(children, *cPtr)
		if e.EmbedChildren {
			var tpl bytes.Buffer
			tmplContent := TemplateContent{Entry: *cPtr}
			err := embededChildTmpl.Execute(&tpl, tmplContent)
			check(err)

			embeddedHTMLStr += tpl.String()
		}
	}

	e.Body += embeddedHTMLStr

	var tpl bytes.Buffer
	tmplContent := TemplateContent{Entry: e, Children: children, NavHTMLString: makeNav(e)}
	err := tmpl.Execute(&tpl, tmplContent)
	check(err)

	htmlStr := tpl.String()
	return htmlStr
}

func main() {
	var entries []Entry
	matches, _ := filepath.Glob("../data/*.ndtl")
	for _, match := range matches {
		file, err := os.Open(match)
		check(err)
		defer file.Close()

		entries = append(entries, loadIndental(file)...)
	}

	linkEntries(entries[:])

	for i, entry := range entries {
		if entry.Parent.EmbedChildren {
			continue
		}

		filepath := "../docs/" + entry.Filename + ".html"
		f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		check(err)

		fmt.Println(i, filepath)
		htmlStr := renderEntryHTML(entry)
		f.WriteString(htmlStr)
	}

	fmt.Println("---")

	// noteFiles, err := ioutil.ReadDir("./notes")
	// check(err)

	// for _, f := range noteFiles {
	// 	fmt.Println(f.Name())
	// }
}
