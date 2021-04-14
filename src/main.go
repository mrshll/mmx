package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"bufio"
	"log"
)

type JrnlRecord struct {
	Date string
	ImgPath string
	Description string
	Parent *Entry
}

type Entry struct {
	MmxDoc
	Parent        *Entry
	EmbedInParent bool
	Children      []*Entry
	FirstImageSrc string
	Inbound       []*Entry
	Outbound      []*Entry
	JrnlRecords []JrnlRecord
}

type TemplateContent struct {
	Entry         Entry
	Children      []Entry
	NavHTMLString string
}

type MakeIndexOptions struct {
	DatesOnly bool
	ShowBref  bool
}

func getEntryFilename(e Entry) string {
	var filename string
	if e.Parent != nil && e.EmbedInParent {
		filename = fmt.Sprintf("%s.html#%s", e.Parent.Slug, e.Slug)
	} else {
		filename = fmt.Sprintf("%s.html", e.Slug)
	}

	return filename
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func formatDate(d time.Time) string {
	return d.Format("2006-01-02")
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func findEntry(entries []Entry, name string) *Entry {
	for i := range entries {
		if strings.ToLower(entries[i].Name) == strings.ToLower(name) {
			return &entries[i]
		}
	}
	panic(fmt.Sprintf("No parent found with name %s", name))
}

func appendEntryIfMissing(entries []*Entry, entryToAppend *Entry) []*Entry {
	for _, e := range entries {
		if e == entryToAppend {
			return entries
		}
	}
	return append(entries, entryToAppend)
}

func linkEntries(entries []Entry) {
	for i := range entries {
		// link Host
		parentPtr := findEntry(entries[:], entries[i].Host)
		entries[i].Parent = parentPtr
		(*parentPtr).Children = append((*parentPtr).Children, &(entries[i]))

		// find incoming
		aTagPattern := regexp.MustCompile(`<a.*?href='([\S ]+?)'>.*?<\/a>`)
		matches := aTagPattern.FindAllStringSubmatch(entries[i].Body, -1)
		for _, match := range matches {
			aTag := match[0]
			outboundHref := match[1]

			if strings.HasPrefix(outboundHref, "http") || strings.HasPrefix(outboundHref, "#") {
				// external
				continue
			} else if outboundHref[0] == '^' {
				// module; embed handled in mmxup
				continue
			} else {

				outboundEntry := findEntry(entries, outboundHref)

				if outboundEntry == nil {
					panic(fmt.Sprintf("No entry found with name '%s' in body of '%s'", outboundHref, entries[i].Body))
				}

				newATag := strings.Replace(aTag, outboundHref, getEntryFilename(*outboundEntry), 1)
				entries[i].Body = strings.Replace(entries[i].Body, aTag, newATag, 1)

				entries[i].Outbound = append(entries[i].Outbound, outboundEntry)
				outboundEntry.Inbound = appendEntryIfMissing(outboundEntry.Inbound, &entries[i])
			}
		}
	}
}

func sortEntries(entries []*Entry) []*Entry {
	sorted := make([]*Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		a := sorted[i]
		b := sorted[j]
		if a.Date.IsZero() && b.Date.IsZero() {
			return a.Name < b.Name
		} else if a.Date.IsZero() {
			return true
		} else if b.Date.IsZero() {
			return false
		}

		return a.Date.After(b.Date)
	})
	return sorted
}

func makeSubNav(e Entry, target Entry) string {
	subnav := "<ul>"
	max := 6

	sortedChildren := sortEntries(e.Children)
	for i, cPtr := range sortedChildren {
		child := *cPtr
		if i >= max {
			if i == max {
				subnav += fmt.Sprintf("<li>and %d more</li>", len(e.Children)-max)
			}

			continue
		}

		if child.Name == e.Name {
			continue // this occurs in the case of root node, i.e. Now
		}

		display := child.Name
		if len(child.Children) > 0 {
			display += "/"
		}

		if child.Name == target.Name {
			subnav += fmt.Sprintf("<li><mark><a href='%s'>%s</a><mark></li>", getEntryFilename(child), display)
		} else {
			subnav += fmt.Sprintf("<li><a href='%s'>%s</a><mark></li>", getEntryFilename(child), display)
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
	for _, cPtr := range e.Children {
		children = append(children, *cPtr)
		if cPtr.EmbedInParent {
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

func makeIndex(indexEntry Entry, entries []*Entry, options MakeIndexOptions) string {
	sortedEntries := sortEntries(entries)

	readingIcon := "<span style='margin-right:10px'>📖</span>"
	projectsIcon := "<span style='margin-right:10px'>🖌</span>"
	musicIcon := "<span style='margin-right:10px'>📻</span>"
	elseIcon := "<span style='margin-right:10px'>🗒️</span>"

	body := indexEntry.Body

	y, _, _ := time.Now().Date()
	y++ // increment y so that the first year is less than current and we write it

	for _, e := range sortedEntries {
		if e.Date.IsZero() && options.DatesOnly {
			continue
		}
		if !e.Date.IsZero() && e.Date.Year() < y {
			y = e.Date.Year()
			body += fmt.Sprintf("<div style='font-size:12px;font-weight:bold;margin-top:20px'>%v</div>", y)
		}

		icon := elseIcon
		crumb := ""
		if e.Parent.Slug == "reading" {
			icon = readingIcon
		} else if e.Parent.Slug == "projects" {
			icon = projectsIcon
		} else if e.Slug == "music" || e.Parent.Slug == "music" {
			icon = musicIcon
		}

		if e.Parent.Slug != indexEntry.Slug {
			crumb = fmt.Sprintf("<a href='%s'>%s</a> > ", getEntryFilename(*e.Parent), e.Parent.Name)
		}

		margin := 0
		if options.ShowBref {
			margin = 10
		}

		body += fmt.Sprintf("<div style='margin-bottom:%dpx'>%s %s<a href='%s'>%s</a>", margin, icon, crumb, getEntryFilename(*e), e.Name)
		if !e.Date.IsZero() {
			body += fmt.Sprintf("<em style='color:lightgrey'> %s</em>", formatDate(e.Date))
		}
		if options.ShowBref {
			body += fmt.Sprintf("<br/><span style='margin-left: 32px'>%s</span>", e.Bref)
		}
		body += "</div>"
	}

	return body
}

func _linkJrnl (jrnlRecord JrnlRecord, entryPtr *Entry) {
	(*entryPtr).JrnlRecords = append((*entryPtr).JrnlRecords, jrnlRecord)
	if (entryPtr.Slug != "index") {
		_linkJrnl(jrnlRecord, entryPtr.Parent)
	}
}

func linkJrnl (entries []Entry) {
	file, err := os.Open("../data/mmx.jrnl")
	check(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Split(line, ", ")
		// default link to now if no parent specified
		entryPtr := findEntry(entries[:], "now")
		if len(args) == 3 {
			entryPtr = findEntry(entries[:], args[2])
		}
		record := JrnlRecord{
			ImgPath: args[0],
			Date: formatDate(parseDate(args[0][:10])),
			Description: args[1],
			Parent: entryPtr,
		}
		_linkJrnl(record, entryPtr)
	}
}

func main() {
	var entries []Entry

	// load .mmx
	matches, _ := filepath.Glob("../data/*.mmx")
	for _, match := range matches {
		file, err := os.Open(match)
		check(err)
		defer file.Close()

		docs := parseFile(file)
		for _, doc := range docs {
			entry := Entry{MmxDoc: doc}
			entries = append(entries, entry)
		}
	}

	linkEntries(entries[:])
	linkJrnl(entries[:])

	for i := range entries {
		if entries[i].IsIndex {
			entries[i].Body = makeIndex(entries[i], entries[i].Children, MakeIndexOptions{ShowBref: true})
		}

		if len(entries[i].JrnlRecords) > 0 {
			entries[i].FirstImageSrc = fmt.Sprintf("img/%s", entries[i].JrnlRecords[0].ImgPath)
		} else {
			imgRegex := regexp.MustCompile(`<img\s.*?src=(?:'|")(?P<src>[^'">]+)(?:'|")`)
			imgMatches := imgRegex.FindAllStringSubmatch(entries[i].Body, 1)
			if len(imgMatches) > 0 {
				entries[i].FirstImageSrc = imgMatches[0][1]
			}
		}
	}

	for i, entry := range entries {
		filepath := "../docs/" + entry.Slug + ".html"
		f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		check(err)

		fmt.Println(i, filepath)
		var htmlStr string
		if entry.Slug == "index" {
			// special case to render timeline
			var entryPtrs []*Entry
			for i := range entries {
				entryPtrs = append(entryPtrs, &entries[i])
			}
			htmlStr = makeIndex(entry, entryPtrs, MakeIndexOptions{DatesOnly: true})
			entry.Body = htmlStr
		}
		htmlStr = renderEntryHTML(entry)
		f.WriteString(htmlStr)
	}

	fmt.Println("---")
}
