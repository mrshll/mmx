# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is MMX?

MMX is a minimal static site generator written in Lua. It converts Markdown content files into a complete HTML website with navigation, RSS feed, and processed images.

## Commands

### Build the site
```bash
./build.sh path/to/content path/to/site
```

### Watch for changes and rebuild (development)
```bash
./run.sh path/to/content path/to/site
```
Uses `fswatch` on macOS, `inotifywait` on Linux.

### Process images
```bash
./processImages.sh path/to/content path/to/site [path/to/specific/image]
```
Requires ImageMagick (`convert`, `identify`). Creates 360px and 720px variants with Floyd-Steinberg dithering (16 colors).

### Run Lua tests
```bash
cd src && lua test.lua
```

## Architecture

### Core Pipeline (`src/mmx.lua`)

1. **Content Discovery**: Scans `DATA_DIR` for `.md` files using `find`
2. **Entry Parsing**: Extracts date (if line 1 is `YYYY-MM-DD`), body, and derives parent/child relationships from directory structure
3. **Markdown Conversion**: Uses `markdown.lua` to convert body to HTML
4. **Embed Processing**: Resolves `![[entry-name]]` syntax to embed other entries
5. **Internal Links**: Converts `[[entry-name]]` or `[[entry-name|display text]]` to HTML links
6. **Image Processing**: Wraps images in `<figure>`, adds lazy loading, references 720px variants
7. **Navigation Rendering**: Builds hierarchical nav from parent-child relationships
8. **HTML Generation**: Applies templates and writes to `SITE_DIR`
9. **RSS Generation**: Creates feed for entries under "Writing" or "Log" parents with dates

### Directory-Based Hierarchy

Entry relationships are determined by filesystem structure:
- `/foo.md` → parent is site root
- `/parent/child.md` → parent is "parent"
- `/parent/index.md` → represents the "parent" entry itself (subtree root)

### Key Files

| File | Purpose |
|------|---------|
| `src/mmx.lua` | Main build orchestration, entry parsing, rendering |
| `src/markdown.lua` | Markdown-to-HTML converter |
| `src/utils.lua` | File I/O, string manipulation, date utilities |
| `src/templates/*.tpl.html` | HTML templates with `{{placeholder}}` syntax |
| `src/templates/*.tpl.rss` | RSS feed templates |
| `static/` | Files copied directly to output (CSS, highlight.js) |
| `links/` | Additional CSS files copied to output |

### Template Placeholders

- `{{EntryBodyHtml}}` - Rendered HTML content
- `{{EntryName}}` - Entry title (derived from filename)
- `{{EntryDate}}` - Date if present
- `{{EntryDestFileName}}` - Output HTML filename
- `{{RSSDate}}` - RFC 2822 formatted date (RSS only)
- `{{Items}}` - RSS item list (RSS template only)

### Special Syntax

- `[[entry-name]]` - Internal link to another entry
- `[[entry-name|display text]]` - Internal link with custom text
- `![[entry-name]]` - Embed another entry's content
- Files/folders starting with `_` are ignored

## Dependencies

- Lua
- GNU `date` (as `gdate` on macOS) for RSS date formatting
- `find` command
- ImageMagick for image processing
- `fswatch` (macOS) or `inotify-tools` (Linux) for watch mode
