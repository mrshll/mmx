package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Entry struct {
	Doc           MmxDoc
	Parent        *Entry
	EmbedInParent bool
	Children      []*Entry
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
		filename = fmt.Sprintf("%s.html#%s", e.Parent.Doc.Slug, e.Doc.Slug)
	} else {
		filename = fmt.Sprintf("%s.html", e.Doc.Slug)
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

// func loadIndental(file *os.File) []Entry {

// 	var entries []Entry
// 	var catchBody bool

// 	scanner := bufio.NewScanner(file)
// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		depth := spad(line, ' ') / 2
// 		lastEntryIndex := len(entries) - 1

// 		if depth == 0 && line != "" {
// 			catchBody = false
// 			name := line
// 			filename := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
// 			entries = append(entries, Entry{Name: name, Filename: filename})
// 		} else if depth == 1 && !catchBody {
// 			catchBody = false
// 			key, value := parseIndentalLine(line)
// 			if key == "DATE" && value != "" {
// 				entries[lastEntryIndex].Date = parseDate(value)
// 			} else if key == "FILE" {
// 				entries[lastEntryIndex].Filename = value
// 			} else if key == "HOST" {
// 				entries[lastEntryIndex].Host = value
// 			} else if key == "BREF" {
// 				entries[lastEntryIndex].Bref = value
// 			} else if key == "EMBD" && value == "true" {
// 				entries[lastEntryIndex].EmbedInParent = true
// 			} else {
// 				catchBody = key == "BODY"
// 			}
// 		} else if depth >= 2 {
// 			if catchBody {
// 				entries[lastEntryIndex].Body += line[4:] + "\n"
// 			}
// 		} else if line == "" && catchBody {
// 			entries[lastEntryIndex].Body += "\n"
// 		}
// 	}

// 	return entries
// }

// func parseFrontLine(line string) (string, string) {
// 	delimiterIndex := strings.Index(line, ":")

// 	if delimiterIndex == -1 {
// 		panic(fmt.Sprintf("No delmiter found in Indental line: %s", line))
// 	}

// 	key := strings.TrimSpace(line[0:delimiterIndex])
// 	val := strings.TrimSpace(line[delimiterIndex+1 : len(line)])
// 	return key, val
// }

// func loadMd(file *os.File) Entry {
// 	scanner := bufio.NewScanner(file)
// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	i := 0
// 	captureFront := false
// 	e := Entry{}

// 	for scanner.Scan() {
// 		i++

// 		line := scanner.Text()
// 		if i == 1 {
// 			captureFront = true
// 			continue
// 		}

// 		if captureFront {
// 			if line == "---" {
// 				captureFront = false
// 				continue
// 			}

// 			key, val := parseFrontLine(line)
// 			if key == "name" {
// 				e.Name = val
// 				e.Filename = strings.ToLower(strings.ReplaceAll(val, " ", "_"))
// 			} else if key == "date" {
// 				e.Date = parseDate(val)
// 			} else if key == "host" {
// 				e.Host = val
// 			} else if key == "bref" {
// 				e.Bref = val
// 			}
// 		} else {
// 			e.Body += line + "\n"
// 		}
// 	}

// 	return e
// }

func findEntry(entries []Entry, name string) *Entry {
	for i := range entries {
		if strings.ToLower(entries[i].Doc.Name) == strings.ToLower(name) {
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

// func processBody(e Entry, entries []Entry) string {
// 	refRegex := regexp.MustCompile(`{[^{}]*}`)
// 	b := e.Body

// 	offset := 0
// 	matches := refRegex.FindAllStringIndex(b, -1)
// 	for _, match := range matches {
// 		cleanMatch := b[match[0]+1+offset : match[1]-1+offset]
// 		matchParts := strings.Split(cleanMatch, "|")

// 		isModule := cleanMatch[0] == '^'
// 		isExternal := strings.Contains(cleanMatch, "http")

// 		var display string
// 		if len(matchParts) > 1 {
// 			display = strings.Join(matchParts[1:], " ")
// 		} else {
// 			display = matchParts[0]
// 		}

// 		var link string
// 		if isModule {
// 			if strings.Contains(matchParts[0], "^bandcamp") {
// 				link = fmt.Sprintf("<iframe style='border: 0; width: 400px; height: 300px;' src='https://bandcamp.com/EmbeddedPlayer/album=%s/size=large/bgcol=ffffff/artwork=small/transparent=true/' seamless></iframe>", matchParts[1])
// 			}
// 			if strings.Contains(matchParts[0], "^youtube") {
// 				link = fmt.Sprintf("<iframe width='560' height='315' src='https://www.youtube-nocookie.com/embed/%s' frameborder='0' allow='accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture' allowfullscreen></iframe>", matchParts[1])
// 			}
// 		} else if isExternal {
// 			// external link
// 			link = fmt.Sprintf("<a href='%s' target='_blank'>[%s]</a>", matchParts[0], display)
// 		} else {
// 			refEntry := findEntry(entries, matchParts[0])
// 			if refEntry == nil {
// 				panic(fmt.Sprintf("No entry found with name %s", matchParts[0]))
// 			}
// 			e.Outgoing = append(e.Outgoing, refEntry)
// 			refEntry.Incoming = appendEntryIfMissing(refEntry.Incoming, &e)

// 			link = fmt.Sprintf("<a href='%s'>{%s}</a>", getEntryFilename(*refEntry), display)
// 		}

// 		b = b[:match[0]+offset] + link + b[match[1]+offset:]
// 		offset += len(link) - (match[1] - match[0])
// 	}

// 	// convert markdown
// 	output := markdown.ToHTML([]byte(b), nil, nil)
// 	b = string(output)

// 	return b
// }

func linkEntries(entries []Entry) {
	for i := range entries {
		parentPtr := findEntry(entries[:], entries[i].Doc.Host)
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
		docA := sortedChildren[i].Doc
		docB := sortedChildren[j].Doc
		if docA.Date.IsZero() && docB.Date.IsZero() {
			return docA.Name < docB.Name
		} else if docA.Date.IsZero() {
			return true
		} else if docB.Date.IsZero() {
			return false
		}

		return docA.Date.After(docB.Date)
	})

	for i, cPtr := range sortedChildren {
		child := *cPtr
		if i >= max {
			if i == max {
				subnav += fmt.Sprintf("<li>and %d more</li>", len(e.Children)-max)
			}

			continue
		}

		if child.Doc.Name == e.Doc.Name {
			continue // this occurs in the case of root node, i.e. Now
		}

		if child.Doc.Name == target.Doc.Name {
			subnav += fmt.Sprintf("<li><mark><a href='%s'>%s/</a><mark></li>", getEntryFilename(child), child.Doc.Name)
		} else {
			subnav += fmt.Sprintf("<li><a href='%s'>%s/</a><mark></li>", getEntryFilename(child), child.Doc.Name)
		}
	}

	subnav += "</ul>"
	return subnav
}

func makeNav(e Entry) string {
	nav := "<nav>"
	if e.Parent == nil {
		panic(fmt.Sprintf("No parent found with name %s", e.Doc.Name))
	}
	if e.Parent.Parent == nil {
		panic(fmt.Sprintf("No parent found with name %s (%s)", e.Parent.Doc.Name, e.Doc.Name))
	}

	// this happens for our root node
	if e.Parent.Parent.Doc.Name == e.Parent.Doc.Name {
		nav += makeSubNav(*e.Parent.Parent, e)
	} else {
		nav += makeSubNav(*e.Parent.Parent, *e.Parent)
	}

	if e.Parent.Parent.Doc.Name != e.Parent.Doc.Name {
		nav += makeSubNav(*e.Parent, e)
	}

	if e.Parent.Doc.Name != e.Doc.Name {
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

	e.Doc.Body += embeddedHTMLStr

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
		return sortedEntries[i].Doc.Date.After(sortedEntries[j].Doc.Date)
	})

	// var activeProjectsBody string
	// for _, e := range entries {
	// 	if e.Name == "Active Projects" {
	// 		activeProjectsBody = e.Body
	// 		break
	// 	}
	// }

	readingIcon := "<span style='margin-right:10px'>📖</span>"
	projectsIcon := "<span style='margin-right:10px'>🖌</span>"
	musicIcon := "<span style='margin-right:10px'>📻</span>"
	elseIcon := "<span style='margin-right:10px'>🗒️</span>"

	nowBody := "<h5>Now</h5>"
	// nowBody := "<h5>Active Projects</h5>"
	// nowBody += activeProjectsBody
	nowBody += "<h5>Updates</h5>"
	y, _, _ := time.Now().Date()
	y++ // increment y so that the first date is less than current and we write it

	for _, e := range sortedEntries {
		if e.Doc.Date.IsZero() {
			continue
		}
		if e.Doc.Date.Year() < y {
			y = e.Doc.Date.Year()
			nowBody += fmt.Sprintf("<div style='font-size:12px;font-weight:bold;margin-top:20px'>%v</div>", y)
		}

		icon := elseIcon
		crumb := ""
		if e.Parent.Doc.Slug == "reading" {
			icon = readingIcon
		} else if e.Parent.Doc.Slug == "active_projects" || e.Parent.Doc.Slug == "dormant_projects" {
			icon = projectsIcon
		} else if e.Doc.Slug == "music" || e.Parent.Doc.Slug == "music" {
			icon = musicIcon
		} else if e.Parent.Doc.Name != "Now" {
			crumb = fmt.Sprintf("<a href='%s'>%s</a> > ", getEntryFilename(*e.Parent), e.Parent.Doc.Name)
		}

		nowBody += fmt.Sprintf("<div>%s %s<a href='%s'>%s</a> <em>%s</em></div>", icon, crumb, getEntryFilename(e), e.Doc.Name, formatDate(e.Doc.Date))
	}

	return nowBody
}

func main() {
	rand.Seed(time.Now().Unix())

	var entries []Entry

	// load .mmx
	matches, _ := filepath.Glob("../data/*.mmx")
	for _, match := range matches {
		file, err := os.Open(match)
		check(err)
		defer file.Close()

		docs := parseFile(file)
		for _, doc := range docs {
			entries = append(entries, Entry{
				Doc: doc,
			})
		}
	}

	linkEntries(entries[:])

	// load .ndtl

	// load .ndtl
	// matches, _ := filepath.Glob("../data/*.ndtl")
	// for _, match := range matches {
	// 	file, err := os.Open(match)
	// 	check(err)
	// 	defer file.Close()

	// 	entries = append(entries, loadIndental(file)...)
	// }

	// load .md
	// matches, _ = filepath.Glob("../data/**/*.md")
	// for _, match := range matches {
	// 	file, err := os.Open(match)
	// 	check(err)
	// 	defer file.Close()

	// 	entries = append(entries, loadMd(file))
	// }
	// matches, _ := filepath.Glob("../data/*.ndtl")
	// for _, match := range matches {
	// 	file, err := os.Open(match)
	// 	check(err)
	// 	defer file.Close()

	// 	entries = append(entries, loadIndental(file)...)
	// }

	// load .md
	// matches, _ = filepath.Glob("../data/**/*.md")
	// for _, match := range matches {
	// 	file, err := os.Open(match)
	// 	check(err)
	// 	defer file.Close()

	// 	entries = append(entries, loadMd(file))
	// }
	// for i := range entries {
	// 	entries[i].Body = processBody(entries[i], entries)

	// 	imgRegex := regexp.MustCompile(`<img\s.*?src=(?:'|")(?P<src>[^'">]+)(?:'|")`)
	// 	imgMatches := imgRegex.FindAllStringSubmatch(entries[i].Body, 1)
	// 	if len(imgMatches) > 0 {
	// 		entries[i].FirstImageSrc = imgMatches[0][1]
	// 	}
	// }

	for i, entry := range entries {
		// 	if entry.EmbedInParent {
		// 		continue
		// 	}

		filepath := "../docs/" + entry.Doc.Slug + ".html"
		f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		check(err)

		fmt.Println(i, filepath)
		var htmlStr string
		if entry.Doc.Slug == "index" {
			// special case to render timeline
			htmlStr = makeIndex(entries)
			entry.Doc.Body = htmlStr
		}
		htmlStr = renderEntryHTML(entry)
		f.WriteString(htmlStr)
	}

	fmt.Println("---")
}
