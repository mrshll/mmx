package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
)

type Entry struct {
	Name          string
	Filename      string
	Bref          string
	Date          time.Time
	Host          string
	Parent        *Entry
	EmbedInParent bool
	Children      []*Entry
	Body          string
	FirstImageSrc string
	Incoming      []*Entry
	Outgoing      []*Entry
}

type TemplateContent struct {
	Entry         Entry
	Children      []Entry
	NavHTMLString string
}

func getEntryFilename(e Entry) string {
	var filename string
	if e.Parent != nil && e.EmbedInParent {
		filename = fmt.Sprintf("%s.html#%s", e.Parent.Filename, e.Filename)
	} else {
		filename = fmt.Sprintf("%s.html", e.Filename)
	}

	return filename
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
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
			filename := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
			entries = append(entries, Entry{Name: name, Filename: filename})
		} else if depth == 1 && !catchBody {
			catchBody = false
			key, value := parseIndentalLine(line)
			if key == "DATE" && value != "" {
				entries[lastEntryIndex].Date = parseDate(value)
			} else if key == "FILE" {
				entries[lastEntryIndex].Filename = value
			} else if key == "HOST" {
				entries[lastEntryIndex].Host = value
			} else if key == "BREF" {
				entries[lastEntryIndex].Bref = value
			} else if key == "EMBD" && value == "true" {
				entries[lastEntryIndex].EmbedInParent = true
			} else {
				catchBody = key == "BODY"
			}
		} else if depth >= 2 {
			if catchBody {
				entries[lastEntryIndex].Body += line[4:] + "\n"
			}
		} else if line == "" && catchBody {
			entries[lastEntryIndex].Body += "\n"
		}
	}

	return entries
}

func parseFrontLine(line string) (string, string) {
	delimiterIndex := strings.Index(line, ":")

	if delimiterIndex == -1 {
		panic(fmt.Sprintf("No delmiter found in Indental line: %s", line))
	}

	key := strings.TrimSpace(line[0:delimiterIndex])
	val := strings.TrimSpace(line[delimiterIndex+1 : len(line)])
	return key, val
}

func loadMd(file *os.File) Entry {
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	i := 0
	captureFront := false
	e := Entry{}

	for scanner.Scan() {
		i++

		line := scanner.Text()
		if i == 1 {
			captureFront = true
			continue
		}

		if captureFront {
			if line == "---" {
				captureFront = false
				continue
			}

			key, val := parseFrontLine(line)
			if key == "name" {
				e.Name = val
				e.Filename = strings.ToLower(strings.ReplaceAll(val, " ", "_"))
			} else if key == "date" {
				e.Date = parseDate(val)
			} else if key == "host" {
				e.Host = val
			} else if key == "bref" {
				e.Bref = val
			}
		} else {
			e.Body += line + "\n"
		}
	}

	return e
}

func findEntry(entries []Entry, name string) *Entry {
	for i := range entries {
		if strings.ToLower(entries[i].Name) == strings.ToLower(name) {
			return &entries[i]
		}
	}
	panic(fmt.Sprintf("No parent found with name %s", name))
}

func makeHr(i int) string {
	poem := []string{
		"⠠⠞⠓⠑ ⠺⠕⠕⠙ ⠞⠓⠗⠥⠎⠓⠂ ⠊⠞ ⠊⠎⠖ ⠠⠝⠕⠺ ⠠⠊ ⠅⠝⠕⠺",
		"⠺⠓⠕ ⠎⠊⠝⠛⠎ ⠞⠓⠁⠞ ⠉⠇⠑⠁⠗ ⠁⠗⠏⠑⠛⠛⠊⠕⠂",
		"⠞⠓⠗⠑⠑ ⠋⠁⠗ ⠝⠕⠞⠑⠎ ⠺⠑⠁⠧⠊⠝⠛",
		"⠊⠝⠞⠕ ⠞⠓⠑ ⠑⠧⠑⠝⠊⠝⠛",
		"⠁⠍⠕⠝⠛ ⠇⠑⠁⠧⠑⠎",
		"⠁⠝⠙ ⠎⠓⠁⠙⠕⠺⠆",
		"⠕⠗ ⠁⠞ ⠙⠁⠺⠝ ⠊⠝ ⠞⠓⠑ ⠺⠕⠕⠙⠎⠂ ⠠⠊⠄⠧⠑ ⠓⠑⠁⠗⠙",
		"⠞⠓⠑ ⠎⠺⠑⠑⠞ ⠁⠎⠉⠑⠝⠙⠊⠝⠛ ⠞⠗⠊⠏⠇⠑ ⠺⠕⠗⠙",
		"⠑⠉⠓⠕⠊⠝⠛ ⠕⠧⠑⠗",
		"⠞⠓⠑ ⠎⠊⠇⠑⠝⠞ ⠗⠊⠧⠑⠗ —",
		"⠃⠥⠞ ⠝⠑⠧⠑⠗",
		"⠎⠑⠑⠝ ⠞⠓⠑ ⠃⠊⠗⠙.",
	}

	poemLine := poem[i%len(poem)]
	return fmt.Sprintf("<div style='color: #ccc; margin: 20px 0;'>%s</div>", poemLine)
}

func replaceHr(b string) string {
	hrRegex := regexp.MustCompile(`(?i)<hr ?\/?>`)
	matches := hrRegex.FindAllString(b, -1)
	for i, match := range matches {
		b = strings.Replace(b, match, makeHr(i), 1)
	}

	return b
}

func appendEntryIfMissing(entries []*Entry, entryToAppend *Entry) []*Entry {
	for _, e := range entries {
		if e == entryToAppend {
			return entries
		}
	}
	return append(entries, entryToAppend)
}

func processBody(e Entry, entries []Entry) string {
	refRegex := regexp.MustCompile(`{[^{}]*}`)
	b := e.Body

	offset := 0
	matches := refRegex.FindAllStringIndex(b, -1)
	for _, match := range matches {
		cleanMatch := b[match[0]+1+offset : match[1]-1+offset]
		matchParts := strings.Split(cleanMatch, "|")

		isModule := cleanMatch[0] == '^'
		isExternal := strings.Contains(cleanMatch, "http")

		var display string
		if len(matchParts) > 1 {
			display = strings.Join(matchParts[1:], " ")
		} else {
			display = matchParts[0]
		}

		var link string
		if isModule {
			if strings.Contains(matchParts[0], "^bandcamp") {
				link = fmt.Sprintf("<iframe style='border: 0; width: 400px; height: 300px;' src='https://bandcamp.com/EmbeddedPlayer/album=%s/size=large/bgcol=ffffff/artwork=small/transparent=true/' seamless></iframe>", matchParts[1])
			}
			if strings.Contains(matchParts[0], "^youtube") {
				link = fmt.Sprintf("<iframe width='560' height='315' src='https://www.youtube-nocookie.com/embed/%s' frameborder='0' allow='accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture' allowfullscreen></iframe>", matchParts[1])
			}
		} else if isExternal {
			// external link
			link = fmt.Sprintf("<a href='%s' target='_blank'>[%s]</a>", matchParts[0], display)
		} else {
			refEntry := findEntry(entries, matchParts[0])
			if refEntry == nil {
				panic(fmt.Sprintf("No entry found with name %s", matchParts[0]))
			}
			e.Outgoing = append(e.Outgoing, refEntry)
			refEntry.Incoming = appendEntryIfMissing(refEntry.Incoming, &e)

			link = fmt.Sprintf("<a href='%s'>{%s}</a>", getEntryFilename(*refEntry), display)
		}

		b = b[:match[0]+offset] + link + b[match[1]+offset:]
		offset += len(link) - (match[1] - match[0])
	}

	// convert markdown
	output := markdown.ToHTML([]byte(b), nil, nil)
	b = string(output)

	return b
}

func linkEntries(entries []Entry) {
	for i := range entries {
		parentPtr := findEntry(entries[:], entries[i].Host)
		entries[i].Parent = parentPtr
		(*parentPtr).Children = append((*parentPtr).Children, &(entries[i]))
	}
}

func makeSubNav(e Entry, target Entry) string {
	subnav := "<ul>"
	max := 10

	sortedChildren := make([]*Entry, len(e.Children))
	copy(sortedChildren, e.Children)
	sort.Slice(sortedChildren, func(i, j int) bool {
		if sortedChildren[i].Date.IsZero() && sortedChildren[j].Date.IsZero() {
			return sortedChildren[i].Name < sortedChildren[j].Name
		} else if sortedChildren[i].Date.IsZero() {
			return true
		} else if sortedChildren[j].Date.IsZero() {
			return false
		}

		return sortedChildren[i].Date.After(sortedChildren[j].Date)
	})

	for i, cPtr := range sortedChildren {
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

		if child.Name == target.Name {
			subnav += fmt.Sprintf("<li><mark><a href='%s'>%s/</a><mark></li>", getEntryFilename(child), child.Name)
		} else {
			subnav += fmt.Sprintf("<li><a href='%s'>%s/</a><mark></li>", getEntryFilename(child), child.Name)
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
		panic(fmt.Sprintf("No parent found with name %s (%s)", e.Parent.Name, e.Name))
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
		"noescape":   noescape,
		"formatDate": formatDate,
	}
	tmpl := template.Must(template.New("entry.html").Funcs(templateFuncs).ParseGlob("./templates/*.html"))
	embededChildTmpl := template.Must(template.New("embeddedChild.html").Funcs(templateFuncs).ParseFiles("./templates/embeddedChild.html", "./templates/incoming.html"))

	var children []Entry
	var embeddedHTMLStr string
	for i, cPtr := range e.Children {
		children = append(children, *cPtr)
		if cPtr.EmbedInParent {
			var tpl bytes.Buffer
			tmplContent := TemplateContent{Entry: *cPtr}
			err := embededChildTmpl.Execute(&tpl, tmplContent)
			check(err)

			embeddedHTMLStr += tpl.String()
			embeddedHTMLStr += makeHr(i)
		}
	}

	e.Body += embeddedHTMLStr

	var tpl bytes.Buffer
	tmplContent := TemplateContent{Entry: e, Children: children, NavHTMLString: makeNav(e)}
	err := tmpl.Execute(&tpl, tmplContent)
	check(err)

	htmlStr := tpl.String()
	htmlStr = replaceHr(htmlStr)

	return htmlStr
}

func makeIndex(entries []Entry) string {
	sortedEntries := make([]Entry, len(entries))
	copy(sortedEntries, entries)
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].Date.After(sortedEntries[j].Date)
	})

	var activeProjectsBody string
	for _, e := range entries {
		if e.Name == "Active Projects" {
			activeProjectsBody = e.Body
			break
		}
	}

	readingIcon := "<span style='margin-right:10px'>📖</span>"
	projectsIcon := "<span style='margin-right:10px'>🖌</span>"
	musicIcon := "<span style='margin-right:10px'>📻</span>"
	elseIcon := "<span style='margin-right:10px'>🗒️</span>"

	homeBody := "<h5>Active Projects</h5>"
	homeBody += activeProjectsBody
	homeBody += "<hr/><h5>Update Timeline</h5>"
	y, _, _ := time.Now().Date()
	y++ // increment y so that the first date is less than current and we write it

	for _, e := range sortedEntries {
		if e.Date.IsZero() {
			continue
		}
		if e.Date.Year() < y {
			y = e.Date.Year()
			homeBody += fmt.Sprintf("<div style='font-size:12px;font-weight:bold;margin-top:20px'>%v</div>", y)
		}

		icon := elseIcon
		crumb := ""
		if e.Parent.Filename == "reading" {
			icon = readingIcon
		} else if e.Parent.Filename == "active_projects" || e.Parent.Filename == "dormant_projects" {
			icon = projectsIcon
		} else if e.Filename == "music" || e.Parent.Filename == "music" {
			icon = musicIcon
		} else if e.Parent.Name != "Home" {
			crumb = fmt.Sprintf("<a href='%s'>%s</a> > ", getEntryFilename(*e.Parent), e.Parent.Name)
		}

		homeBody += fmt.Sprintf("<div>%s %s<a href='%s'>%s</a> <em>%s</em></div>", icon, crumb, getEntryFilename(e), e.Name, formatDate(e.Date))
	}

	return homeBody
}

func main() {
	rand.Seed(time.Now().Unix())

	var entries []Entry

	// load .ndtl
	matches, _ := filepath.Glob("../data/*.ndtl")
	for _, match := range matches {
		file, err := os.Open(match)
		check(err)
		defer file.Close()

		entries = append(entries, loadIndental(file)...)
	}

	// load .md
	matches, _ = filepath.Glob("../data/**/*.md")
	for _, match := range matches {
		file, err := os.Open(match)
		check(err)
		defer file.Close()

		entries = append(entries, loadMd(file))
	}

	linkEntries(entries[:])
	for i := range entries {
		entries[i].Body = processBody(entries[i], entries)

		imgRegex := regexp.MustCompile(`<img\s.*?src=(?:'|")(?P<src>[^'">]+)(?:'|")`)
		imgMatches := imgRegex.FindAllStringSubmatch(entries[i].Body, 1)
		if len(imgMatches) > 0 {
			entries[i].FirstImageSrc = imgMatches[0][1]
		}
	}

	for i, entry := range entries {
		if entry.EmbedInParent {
			continue
		}

		filepath := "../docs/" + entry.Filename + ".html"
		f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		check(err)

		fmt.Println(i, filepath)
		var htmlStr string
		if entry.Filename == "index" {
			// special case to render timeline
			htmlStr = makeIndex(entries)
			entry.Body = htmlStr
		}
		htmlStr = renderEntryHTML(entry)
		f.WriteString(htmlStr)
	}

	fmt.Println("---")
}
