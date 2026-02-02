-- MMX Unit Tests
-- Run with: lua test.lua

local markdown = require "markdown"
local utils = require "utils"

-- Test framework
local tests_run = 0
local tests_passed = 0
local tests_failed = 0

local function test(name, fn)
    tests_run = tests_run + 1
    local success, err = pcall(fn)
    if success then
        tests_passed = tests_passed + 1
        print("  PASS: " .. name)
    else
        tests_failed = tests_failed + 1
        print("  FAIL: " .. name)
        print("        " .. tostring(err))
    end
end

local function assert_eq(expected, actual, msg)
    if expected ~= actual then
        error((msg or "assertion failed") ..
              "\n        expected: " .. tostring(expected) ..
              "\n        actual:   " .. tostring(actual))
    end
end

local function assert_true(value, msg)
    if not value then
        error(msg or "expected true, got false")
    end
end

local function assert_false(value, msg)
    if value then
        error(msg or "expected false, got true")
    end
end

local function assert_nil(value, msg)
    if value ~= nil then
        error((msg or "expected nil") .. ", got: " .. tostring(value))
    end
end

local function assert_match(pattern, str, msg)
    if not str:match(pattern) then
        error((msg or "pattern not found") ..
              "\n        pattern: " .. pattern ..
              "\n        string:  " .. str)
    end
end

-- =============================================================================
-- Date Extraction Tests
-- =============================================================================
print("\n=== Date Extraction ===")

-- Helper that mimics the date extraction logic from mmx.lua
local function extract_date(body, filename)
    local date
    local date_start, date_end = body:find("%d%d%d%d%-%d%d%-%d%d")
    if date_start == 1 then
        date = body:sub(date_start, date_end)
        body = body:sub(date_end + 1)
    elseif filename:match("^%d%d%d%d%-%d%d%-%d%d$") then
        date = filename
    end
    return date, body
end

test("date from first line of body", function()
    local date, body = extract_date("2023-05-15\nSome content here", "my-post")
    assert_eq("2023-05-15", date)
    assert_eq("\nSome content here", body)
end)

test("date from filename when body has no date", function()
    local date, body = extract_date("Some content without a date", "2023-05-15")
    assert_eq("2023-05-15", date)
    assert_eq("Some content without a date", body)
end)

test("body date takes precedence over filename date", function()
    local date, body = extract_date("2023-01-01\nContent", "2023-12-31")
    assert_eq("2023-01-01", date, "body date should take precedence")
end)

test("no date when neither body nor filename has date", function()
    local date, body = extract_date("Just some content", "my-regular-post")
    assert_nil(date)
    assert_eq("Just some content", body)
end)

test("date not extracted from middle of body", function()
    local date, body = extract_date("This happened on 2023-05-15 at noon", "post")
    assert_nil(date, "date in middle of body should not be extracted")
end)

test("filename with extra characters is not treated as date", function()
    local date, _ = extract_date("Content", "2023-05-15-extra")
    assert_nil(date, "filename with suffix should not match")

    date, _ = extract_date("Content", "prefix-2023-05-15")
    assert_nil(date, "filename with prefix should not match")
end)

test("date stripped from body correctly", function()
    local date, body = extract_date("2024-01-01\n\nFirst paragraph", "post")
    assert_eq("2024-01-01", date)
    assert_eq("\n\nFirst paragraph", body, "rest of body should be preserved")
end)

-- =============================================================================
-- Utils: String Functions
-- =============================================================================
print("\n=== Utils: String Functions ===")

test("starts_with matches prefix", function()
    assert_true(utils.starts_with("hello world", "hello"))
    assert_true(utils.starts_with("hello", "hello"))
    assert_true(utils.starts_with("hello", ""))
end)

test("starts_with rejects non-prefix", function()
    assert_false(utils.starts_with("hello world", "world"))
    assert_false(utils.starts_with("hello", "hello world"))
end)

test("ends_with matches suffix", function()
    assert_true(utils.ends_with("hello.md", ".md"))
    assert_true(utils.ends_with("test", "test"))
    assert_true(utils.ends_with("anything", ""))
end)

test("ends_with rejects non-suffix", function()
    assert_false(utils.ends_with("hello.md", ".txt"))
    assert_false(utils.ends_with("short", "longer"))
end)

test("split by delimiter", function()
    local parts = utils.split("a/b/c", "/")
    assert_eq(3, #parts)
    assert_eq("a", parts[1])
    assert_eq("b", parts[2])
    assert_eq("c", parts[3])
end)

test("split handles single element", function()
    local parts = utils.split("single", "/")
    assert_eq(1, #parts)
    assert_eq("single", parts[1])
end)

test("slugify converts to lowercase with underscores", function()
    assert_eq("hello_world", utils.slugify("Hello World"))
    assert_eq("my_post_title", utils.slugify("My Post Title"))
end)

test("slugify handles punctuation", function()
    assert_eq("what_s_up_", utils.slugify("What's Up?"))
    assert_eq("test_file_name", utils.slugify("test-file-name"))
end)

test("capitalize first letter", function()
    assert_eq("Hello", utils.capitalize("hello"))
    assert_eq("Already", utils.capitalize("Already"))
end)

test("title_case capitalizes words", function()
    assert_eq("Hello World", utils.title_case("hello world"))
end)

-- =============================================================================
-- Utils: Table Functions
-- =============================================================================
print("\n=== Utils: Table Functions ===")

test("has_keys returns true when all keys present", function()
    local t = {a = 1, b = 2, c = 3}
    assert_true(utils.has_keys(t, {"a", "b"}))
    assert_true(utils.has_keys(t, {"a", "b", "c"}))
end)

test("has_keys returns false when key missing", function()
    local t = {a = 1, b = 2}
    assert_false(utils.has_keys(t, {"a", "c"}))
    assert_false(utils.has_keys(t, {"d"}))
end)

test("get_key_case_insensitive finds variations", function()
    local t = {Hello = "value1", world = "value2", ["Title Case"] = "value3"}
    assert_eq("value1", utils.get_key_case_insensitive(t, "Hello"))
    assert_eq("value1", utils.get_key_case_insensitive(t, "hello"))
    assert_eq("value2", utils.get_key_case_insensitive(t, "world"))
    assert_eq("value2", utils.get_key_case_insensitive(t, "WORLD"))
end)

-- =============================================================================
-- Markdown Processing
-- =============================================================================
print("\n=== Markdown Processing ===")

test("converts headers", function()
    assert_match("<h1>", markdown("# Header"))
    assert_match("<h2>", markdown("## Header"))
    assert_match("<h3>", markdown("### Header"))
end)

test("converts bold and italic", function()
    assert_match("<strong>", markdown("**bold**"))
    assert_match("<em>", markdown("*italic*"))
    assert_match("<em>", markdown("_italic_"))
end)

test("converts unordered lists", function()
    local result = markdown("- item 1\n- item 2")
    assert_match("<ul>", result)
    assert_match("<li>", result)
end)

test("converts links", function()
    local result = markdown("[text](http://example.com)")
    assert_match('<a href="http://example.com"', result)
end)

test("converts images", function()
    local result = markdown("![alt text](image.png)")
    assert_match('<img', result)
    assert_match('src="image.png"', result)
    assert_match('alt="alt text"', result)
end)

test("converts code blocks", function()
    local result = markdown("```\ncode here\n```")
    assert_match("<pre>", result)
    assert_match("<code>", result)
end)

test("converts inline code", function()
    local result = markdown("use `code` here")
    assert_match("<code>", result)
end)

-- =============================================================================
-- Internal Link Syntax
-- =============================================================================
print("\n=== Internal Link Syntax ===")

-- Helper that mimics process_internal_links from mmx.lua
local function process_internal_links(str, entry, entries)
    return str:gsub("[^!]%[%[[^%]]+%]%]", function(match)
        local parts = utils.split(match:sub(4, -3), "|")
        local linked_name = parts[1]
        local display_name = parts[2] or linked_name
        local e = utils.get_key_case_insensitive(entries, linked_name)
        if e == nil then
            return "{" .. display_name .. "}"
        end
        return match:sub(1, 1) .. "<a href=\"" .. e.dest_file_name .. "\">{" .. display_name .. "}</a>"
    end)
end

local mock_entries = {
    ["About"] = {name = "About", dest_file_name = "about.html"},
    ["Contact"] = {name = "Contact", dest_file_name = "contact.html"},
}
local mock_entry = {name = "Test"}

test("internal link basic syntax", function()
    local result = process_internal_links("See [[About]] for more", mock_entry, mock_entries)
    assert_match('href="about.html"', result)
    assert_match("{About}", result)
end)

test("internal link with display text", function()
    local result = process_internal_links("See [[About|my about page]]", mock_entry, mock_entries)
    assert_match('href="about.html"', result)
    assert_match("{my about page}", result)
end)

test("internal link case insensitive", function()
    local result = process_internal_links("See [[about]] page", mock_entry, mock_entries)
    assert_match('href="about.html"', result)
end)

test("internal link to non-existent entry", function()
    local result = process_internal_links("See [[NonExistent]] page", mock_entry, mock_entries)
    assert_match("{NonExistent}", result)
    -- Should not contain href for missing entry
    assert_false(result:match('href="nonexistent.html"'))
end)

test("embed syntax not processed as link", function()
    -- The regex [^!] ensures ![[...]] is not matched
    local result = process_internal_links("Embed: ![[About]]", mock_entry, mock_entries)
    -- Should remain unchanged (not converted to link)
    assert_match("!%[%[About%]%]", result)
end)

-- =============================================================================
-- Embed Syntax
-- =============================================================================
print("\n=== Embed Syntax ===")

-- Helper that mimics process_embeds from mmx.lua
local function process_embeds(body_html, entries)
    return body_html:gsub("!%[%[([^%]]+)%]%]", function(embedded_entry_name)
        for _, entry in pairs(entries) do
            if entry.name == embedded_entry_name then
                return "<article><h2>" ..
                    entry.name ..
                    "</h2><a href=" .. entry.dest_file_name .. ">#</a>" ..
                    markdown(entry.body_raw) .. "</article>"
            end
        end
        return "![[" .. embedded_entry_name .. "]]"  -- Return unchanged if not found
    end)
end

local embed_entries = {
    ["Quote"] = {
        name = "Quote",
        dest_file_name = "quote.html",
        body_raw = "A wise saying"
    },
}

test("embed inserts entry content", function()
    local result = process_embeds("Before ![[Quote]] After", embed_entries)
    assert_match("<article>", result)
    assert_match("<h2>Quote</h2>", result)
    assert_match("A wise saying", result)
end)

test("embed includes link to original", function()
    local result = process_embeds("![[Quote]]", embed_entries)
    assert_match('href=quote.html', result)
    assert_match(">#</a>", result)
end)

-- =============================================================================
-- Image Processing
-- =============================================================================
print("\n=== Image Processing ===")

local MEDIA_DIR_NAME = 'media'

-- Simplified version of process_images for testing
local function process_images(str)
    return str:gsub("<img [^>]+>", function(img_tag)
        local alt = img_tag:match("alt=\"([^\"]*)\"") or ""
        local src_pattern = "src=\"([^\"]+)\""
        local src = img_tag:match(src_pattern)

        -- Handle video files
        local video_exts = {"webm", "mp4", "mov", "ogv"}
        for _, ext in ipairs(video_exts) do
            if utils.ends_with(src, ext) then
                local video_src = MEDIA_DIR_NAME .. "/" .. src
                return "<figure><video controls><source src=\"" .. video_src .. "\"></video><figcaption>" .. alt .. "</figcaption></figure>"
            end
        end

        local processed_img_tag = img_tag:gsub(src_pattern, function(s)
            if utils.starts_with(s, "http") then
                return "src=\"" .. s .. "\""
            elseif (not (utils.ends_with(s, "jpg") or utils.ends_with(s, "jpeg") or utils.ends_with(s, "png"))) then
                return "src=\"/" .. MEDIA_DIR_NAME .. "/" .. s .. "\""
            else
                s = MEDIA_DIR_NAME .. "/" .. s
            end

            local parts = utils.split(s, ".")
            if not (utils.ends_with(parts[1], "360") or utils.ends_with(parts[1], "720")) then
                s = parts[1] .. "-720." .. parts[2]
            end

            return "loading=\"lazy\" src=\"" .. s .. "\""
        end)

        return "<figure><a href=\"" .. MEDIA_DIR_NAME .. "/" .. src .. "\">" ..
               processed_img_tag .. "</a><figcaption>" .. alt .. "</figcaption></figure>"
    end)
end

test("image wrapped in figure", function()
    local result = process_images('<img src="photo.jpg" alt="A photo">')
    assert_match("<figure>", result)
    assert_match("</figure>", result)
    assert_match("<figcaption>A photo</figcaption>", result)
end)

test("image gets lazy loading", function()
    local result = process_images('<img src="photo.jpg" alt="">')
    assert_match('loading="lazy"', result)
end)

test("image src rewritten to 720px variant", function()
    local result = process_images('<img src="photo.jpg" alt="">')
    assert_match('src="media/photo%-720.jpg"', result)
end)

test("already sized image not modified", function()
    local result = process_images('<img src="photo-720.jpg" alt="">')
    assert_match('src="media/photo%-720.jpg"', result)
    -- Should not become photo-720-720.jpg
    assert_false(result:match("720%-720"))

    local result360 = process_images('<img src="photo-360.png" alt="">')
    assert_match('src="media/photo%-360.png"', result360)
end)

test("external URL src attribute preserved", function()
    local result = process_images('<img src="https://example.com/photo.jpg" alt="">')
    assert_match('src="https://example.com/photo.jpg"', result)
    -- Note: the figure anchor href still gets media/ prefix (known limitation)
    -- External images don't get lazy loading added
    assert_false(result:match('loading="lazy"'))
end)

test("video files converted to video tag", function()
    local result = process_images('<img src="clip.mp4" alt="Video">')
    assert_match("<video controls>", result)
    assert_match("<source src", result)
    assert_match("</video>", result)
end)

test("non-jpg/png images get media path without resize", function()
    local result = process_images('<img src="diagram.svg" alt="">')
    assert_match('src="/media/diagram.svg"', result)
end)

-- =============================================================================
-- Entry Hierarchy
-- =============================================================================
print("\n=== Entry Hierarchy ===")

-- Helper to derive parent from path parts (mimics mmx.lua logic)
local function derive_parent(file_path, site_name)
    local parts = utils.split(file_path, "/")
    local parent_name = parts[#parts - 1] or site_name
    return parent_name
end

test("root level entry has site as parent", function()
    local parent = derive_parent("/about.md", "mysite.com")
    assert_eq("mysite.com", parent)
end)

test("nested entry has directory as parent", function()
    local parent = derive_parent("/Writing/my-post.md", "mysite.com")
    assert_eq("Writing", parent)
end)

test("deeply nested entry has immediate parent", function()
    local parent = derive_parent("/Projects/Web/app.md", "mysite.com")
    assert_eq("Web", parent)
end)

-- Helper to check if path should be ignored (underscore prefix)
local function should_ignore_path(file_path)
    local parts = utils.split(file_path, "/")
    for _, part in pairs(parts) do
        if utils.starts_with(part, '_') then
            return true
        end
    end
    return false
end

test("underscore prefix file ignored", function()
    assert_true(should_ignore_path("/_drafts/post.md"))
    assert_true(should_ignore_path("/Writing/_hidden.md"))
end)

test("underscore in middle of name not ignored", function()
    assert_false(should_ignore_path("/my_post.md"))
    assert_false(should_ignore_path("/Writing/post_title.md"))
end)

-- =============================================================================
-- Template Substitution
-- =============================================================================
print("\n=== Template Substitution ===")

local function sub_entry_fields(str, entry)
    return str:gsub("{{(%w+)}}", {
        ["EntryBodyHtml"] = entry.body_html,
        ["EntryName"] = entry.name,
        ["EntryDate"] = entry.date,
        ["EntryDestFileName"] = entry.dest_file_name
    })
end

test("substitutes entry name", function()
    local entry = {name = "Test Post", body_html = "", date = "", dest_file_name = ""}
    local result = sub_entry_fields("<h1>{{EntryName}}</h1>", entry)
    assert_eq("<h1>Test Post</h1>", result)
end)

test("substitutes entry date", function()
    local entry = {name = "", body_html = "", date = "2024-01-15", dest_file_name = ""}
    local result = sub_entry_fields("Posted: {{EntryDate}}", entry)
    assert_eq("Posted: 2024-01-15", result)
end)

test("substitutes dest filename", function()
    local entry = {name = "", body_html = "", date = "", dest_file_name = "my_post.html"}
    local result = sub_entry_fields('<a href="{{EntryDestFileName}}">', entry)
    assert_eq('<a href="my_post.html">', result)
end)

test("substitutes body html", function()
    local entry = {name = "", body_html = "<p>Content</p>", date = "", dest_file_name = ""}
    local result = sub_entry_fields("<main>{{EntryBodyHtml}}</main>", entry)
    assert_eq("<main><p>Content</p></main>", result)
end)

test("unknown placeholders remain unchanged", function()
    local entry = {name = "Test", body_html = "", date = "", dest_file_name = ""}
    local result = sub_entry_fields("{{Unknown}} and {{EntryName}}", entry)
    assert_eq("{{Unknown}} and Test", result)
end)

-- =============================================================================
-- Sorting Logic
-- =============================================================================
print("\n=== Sorting Logic ===")

local MIN_DATE = "0000-00-00"

local function sort_entries(entries)
    local sorted = {}
    for _, e in pairs(entries) do
        table.insert(sorted, e)
    end
    table.sort(sorted, function(a, b)
        if a.date or b.date then
            return (a.date or MIN_DATE) > (b.date or MIN_DATE)
        else
            return a.name < b.name
        end
    end)
    return sorted
end

test("entries sorted by date descending", function()
    local entries = {
        {name = "Old", date = "2023-01-01"},
        {name = "New", date = "2024-01-01"},
        {name = "Mid", date = "2023-06-01"},
    }
    local sorted = sort_entries(entries)
    assert_eq("New", sorted[1].name)
    assert_eq("Mid", sorted[2].name)
    assert_eq("Old", sorted[3].name)
end)

test("entries without dates sorted by name ascending", function()
    local entries = {
        {name = "Zebra"},
        {name = "Apple"},
        {name = "Mango"},
    }
    local sorted = sort_entries(entries)
    assert_eq("Apple", sorted[1].name)
    assert_eq("Mango", sorted[2].name)
    assert_eq("Zebra", sorted[3].name)
end)

test("dated entries come before undated", function()
    local entries = {
        {name = "Undated"},
        {name = "Dated", date = "2020-01-01"},
    }
    local sorted = sort_entries(entries)
    assert_eq("Dated", sorted[1].name)
    assert_eq("Undated", sorted[2].name)
end)

-- =============================================================================
-- RSS Entry Filtering
-- =============================================================================
print("\n=== RSS Entry Filtering ===")

local function is_rss_entry(entry)
    return (entry.parent_name == "Writing" or entry.parent_name == "Log") and entry.date ~= nil
end

test("Writing entries with date included in RSS", function()
    assert_true(is_rss_entry({parent_name = "Writing", date = "2024-01-01"}))
end)

test("Log entries with date included in RSS", function()
    assert_true(is_rss_entry({parent_name = "Log", date = "2024-01-01"}))
end)

test("entries without date excluded from RSS", function()
    assert_false(is_rss_entry({parent_name = "Writing", date = nil}))
    assert_false(is_rss_entry({parent_name = "Log"}))
end)

test("other parent entries excluded from RSS", function()
    assert_false(is_rss_entry({parent_name = "Projects", date = "2024-01-01"}))
    assert_false(is_rss_entry({parent_name = "About", date = "2024-01-01"}))
end)

-- =============================================================================
-- Summary
-- =============================================================================
print("\n" .. string.rep("=", 50))
print(string.format("Tests: %d passed, %d failed, %d total",
    tests_passed, tests_failed, tests_run))
print(string.rep("=", 50))

if tests_failed > 0 then
    os.exit(1)
end
