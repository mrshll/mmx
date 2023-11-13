local markdown = require "markdown"
local utils = require "utils"

if not arg[1] or arg[1] == '' then
    error('no argument for the site was provided')
end

local SITE_NAME = 'mrshll.com'
local INDEX_NAME = "index"
local MIN_DATE = "0000-00-00"
local MEDIA_DIR_NAME = 'media'
local DATA_DIR = "../" .. arg[1] .. "/data"
local SITE_DIR = "../" .. arg[1] .. "/site"
local DATA_EXT = ".md"

-- store all image paths to do a check to make sure they exist
local media_paths = {}

local function sub_entry_fields(str, entry)
    return str:gsub("{{(%w+)}}", {
        ["EntryBodyHtml"] = entry.body_html,
        ["EntryName"] = entry.name,
        ["EntryDate"] = entry.date,
        ["EntryDestFileName"] = entry.dest_file_name
    })
end

local function render_head(entry)
    return sub_entry_fields(utils.read_file("templates/head.tpl.html"), entry)
end

local function render_footer()
    return utils.read_file("templates/footer.tpl.html")
end

local function make_sorted_entry_iterator(entries)
    return utils.spairs(entries, function(es, name_a, name_b)
        local entry_a = es[name_a]
        local entry_b = es[name_b]
        if entry_a.date or entry_b.date then
            return (entry_a.date or MIN_DATE) > (entry_b.date or MIN_DATE)
        else
            return name_a < name_b
        end
    end)
end

local function render_nav_section(entry, siblings)
    local acc = "<ul>"
    local i = 0
    for _, e in make_sorted_entry_iterator(siblings) do
        if i == 6 then
            acc = acc .. "<details><summary></summary>"
        end

        i = i + 1
        local is_selected = entry ~= nil and e.name == entry.name
        acc = acc .. "<li>"

        if is_selected then
            acc = acc .. "<mark>"
        end
        acc = acc .. "<a href=\"" .. e.dest_file_name .. "\">" .. e.name .. "</a>"
        if is_selected then
            acc = acc .. "</mark>"
        end
        acc = acc .. "</li>"
    end
    if i >= 6 then
        acc = acc .. "</details>"
    end
    return acc .. "</ul>"
end

local function render_nav(entry, entries)
    local acc = ""

    local children = {}
    for _, e in make_sorted_entry_iterator(entries) do
        if e.parent_name == entry.name and e.name ~= SITE_NAME then
            table.insert(children, e)
        end
    end
    acc = acc .. render_nav_section(nil, children)

    -- go until we are at root
    while entry.name ~= SITE_NAME do
        local siblings = {}
        local parent = nil
        for _, e in make_sorted_entry_iterator(entries) do
            if entry.parent_name == e.name then
                parent = e
            elseif e.parent_name == entry.parent_name then
                table.insert(siblings, e)
            end
        end
        acc = render_nav_section(entry, siblings) .. acc
        if parent == nil then
            error(entry.name .. "(" .. entry.src_path .. ") has no parent")
        end

        entry = parent
    end

    return "<nav>" .. render_nav_section(entry, {entry}) .. acc .. "</nav>"
end

local function process_images(str)
    return str:gsub("<img [^>]+>", function(img_tag)
        local alt = img_tag:match("alt=\"([^\"]*)\"")
        local src_pattern = "src=\"([^\"]+)\""
        local src = img_tag:match(src_pattern)

        local processed_img_tag = img_tag:gsub(src_pattern, function(s)
            -- we don't compress/process other image formats
            if utils.starts_with(s, "http") then
                return "src=\"" .. s .. "\""
            elseif (not (utils.ends_with(s, "jpg") or utils.ends_with(s, "jpeg") or utils.ends_with(s, "png"))) then
                return "src=\"/" .. MEDIA_DIR_NAME .. "/" .. s .. "\""
            else
                -- we store attachments in this directory
                s = MEDIA_DIR_NAME .. "/" .. s
            end

            local parts = utils.split(s, ".")

            -- if we've hard-specified one of our resolutions, don't append it
            if not (utils.ends_with(parts[1], "360") or utils.ends_with(parts[1], "720")) then
                s = parts[1] .. "-720." .. parts[2]
            end

            table.insert(media_paths, s)

            return "loading=\"lazy\" src=\"" .. s .. "\""
        end)
        return
            "<figure><a href=\"" .. MEDIA_DIR_NAME .. "/" .. src .. "\">" .. processed_img_tag .. "</a><figcaption>" ..
                alt .. "</figcaption></figure>"
    end)
end

local function process_internal_links(string, entry, entries)
    return string:gsub("%[%[([^%]]+)%]%]", function(match)
        local parts = utils.split(match, "|")
        local linked_name = parts[1]
        local display_name = parts[2] or linked_name
        local e = utils.get_key_case_insensitive(entries, linked_name)
        if e == nil then
            print("Warning, linked entry \"" .. linked_name .. "\" not found when rendering \"" .. entry.name .. "\"")
            return "{" .. display_name .. "}"
        end

        return "<a href=\"" .. e.dest_file_name .. "\">{" .. display_name .. "}</a>"
    end)
end

local function render_body(entry, entries)
    return process_images(process_internal_links(sub_entry_fields("<h1>{{EntryName}}</h1>" ..
                                                                      (entry.date ~= nil and
                                                                          "<div style='color:#ccc'>last updated {{EntryDate}}</div>" or
                                                                          "") .. "{{EntryBodyHtml}}", entry), entry,
        entries))
end

local function render_entry(entry, entries)
    local html = string.format(
        "<!doctype html><html>%s<body><div class=\"content\"><header>%s</header><main id=\"entry-body\">%s</main><p style=\"color:#ccc\"><em>Compiled %s</em></p></div>%s</body></html>",
        render_head(entry), render_nav(entry, entries), render_body(entry, entries), utils.today(), render_footer())
    utils.write_file(SITE_DIR .. "/" .. entry.dest_file_name, html)
end

local function render_rss(rss_entries)
    local rss_template = utils.read_file("templates/feed.tpl.rss")
    local item_template = utils.read_file("templates/item.tpl.rss")
    local items_str = ""
    for _, e in make_sorted_entry_iterator(rss_entries) do
        local rss_date = utils.rss_date(e.date)
        items_str = items_str ..
                        sub_entry_fields(item_template, e):gsub("{{RSSDate}}", rss_date)
                :gsub("src=\"img", "src=\"https://mrshll.com/img"):gsub("%%", "%%%%") -- this escapes %, which is lua's escape char, otherwise the final gsub fails
    end
    utils.write_file(SITE_DIR .. "/feed.rss", rss_template:gsub("{{Items}}", items_str))
end

local entries = {}
local file_paths = utils.list_files(DATA_DIR, DATA_EXT)
for _, file_path in pairs(file_paths) do

    -- split on slash to get path fragments
    local parts = utils.split(file_path, "/")

    for _, key in pairs(parts) do
        if utils.starts_with(key, '_') then
            goto continue
        end
    end

    -- directory scheme is /entry-name or [/parent's parent]/parent/entry-name (recursive)
    local parent_name = parts[#parts - 1] or SITE_NAME

    -- remove the file extension
    local name = parts[#parts]:sub(0, -1 * #DATA_EXT - 1)
    local dest_file_name = utils.slugify(name) .. ".html"

    if name == INDEX_NAME then
        name = parent_name
        if parent_name == SITE_NAME then
            -- root node
            dest_file_name = "index.html"
        else
            -- subtree root nodes
            dest_file_name = name .. ".html"
            parent_name = parts[#parts - 2] or SITE_NAME
        end
    end

    local body = utils.read_file(DATA_DIR .. file_path)

    local date
    local date_start, date_end = body:find("%d%d%d%d%-%d%d%-%d%d")
    if date_start == 1 then
        date = body:sub(date_start, date_end)
        body = body:sub(date_end + 1)
    end

    local html = markdown(body)
    entries[name] = {
        name = name,
        src_path = file_path,
        dest_file_name = dest_file_name,
        parent_name = parent_name,
        body_raw = body,
        body_html = html,
        date = date
    }

    ::continue::
end

local i = 0
for _, entry in pairs(entries) do
    render_entry(entry, entries)
    i = i + 1
end

local rss_entries = {}
for _, e in pairs(entries) do
    if e.parent_name == "Writing" and e.date ~= nil then
        table.insert(rss_entries, e)
    end
end
render_rss(rss_entries)

for _, media_path in pairs(media_paths) do
    local path = SITE_DIR .. "/" .. media_path
    if not utils.file_exists(path) then
        print("referenced media not found", path)
    end
end

print("Rendered " .. i .. " entries.")

