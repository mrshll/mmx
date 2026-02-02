local markdown = require "markdown"
local utils = require "utils"

if not arg[1] or arg[1] == '' or not arg[2] or arg[2] == '' then
    error('usage: lua mmx.lua path/to/content path/to/site')
end

local SITE_NAME = 'mrshll.com'
local INDEX_NAME = "index"
local MIN_DATE = "0000-00-00"
local MEDIA_DIR_NAME = 'media'
local DATA_DIR = "../" .. arg[1]
local SITE_DIR = "../" .. arg[2]
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

local function render_nav(current_entry, entries)
    -- Build set of ancestor names for current entry
    local ancestors = {}
    local temp = current_entry
    while temp do
        ancestors[temp.name] = true
        if temp.name == SITE_NAME then
            break  -- Stop at root
        end
        local parent = nil
        for _, e in pairs(entries) do
            if e.name == temp.parent_name then
                parent = e
                break
            end
        end
        temp = parent
    end

    -- Get sorted children of a parent
    local function get_children(parent_name)
        local children = {}
        for _, e in pairs(entries) do
            if e.parent_name == parent_name and e.name ~= SITE_NAME then
                table.insert(children, e)
            end
        end
        -- Sort by date desc, then name asc
        table.sort(children, function(a, b)
            if a.date or b.date then
                return (a.date or MIN_DATE) > (b.date or MIN_DATE)
            else
                return a.name < b.name
            end
        end)
        return children
    end

    local items = {}

    -- Recursive function to collect tree items
    local function collect_items(parent_name, prefix)
        local children = get_children(parent_name)

        for i, child in ipairs(children) do
            local is_last = (i == #children)
            local is_current = (child.name == current_entry.name)
            local is_ancestor = ancestors[child.name]

            -- Choose tree characters
            local branch = is_last and "└── " or "├── "
            local child_prefix = prefix .. (is_last and "    " or "│   ")

            table.insert(items, {
                entry = child,
                prefix = prefix .. branch,
                is_current = is_current
            })

            -- Recursively collect children if this is an ancestor or current
            if is_ancestor or is_current then
                collect_items(child.name, child_prefix)
            end
        end
    end

    -- Start from root
    collect_items(SITE_NAME, "")

    local root_entry = entries[SITE_NAME]
    local html = "<ul>"
    local is_root_current = (current_entry.name == SITE_NAME)
    html = html .. "<li>"
    if is_root_current then
        html = html .. "<mark>"
    end
    html = html .. "<a href=\"" .. root_entry.dest_file_name .. "\">" .. SITE_NAME .. "</a>"
    if is_root_current then
        html = html .. "</mark>"
    end
    html = html .. "</li>"

    for _, item in ipairs(items) do
        html = html .. "<li>"
        html = html .. "<span class=\"tree-prefix\">" .. item.prefix .. "</span>"
        if item.is_current then
            html = html .. "<mark>"
        end
        html = html .. "<a href=\"" .. item.entry.dest_file_name .. "\">" .. item.entry.name .. "</a>"
        if item.is_current then
            html = html .. "</mark>"
        end
        html = html .. "</li>"
    end

    html = html .. "</ul>"

    return "<nav>" .. html .. "</nav>"
end

local function process_images(str)
    return str:gsub("<img [^>]+>", function(img_tag)
        local alt = img_tag:match("alt=\"([^\"]*)\"")
        local src_pattern = "src=\"([^\"]+)\""
        local src = img_tag:match(src_pattern)

        -- Handle video files
        local video_exts = {"webm", "mp4", "mov", "ogv"}
        for _, ext in ipairs(video_exts) do
            if utils.ends_with(src, ext) then
                local video_src =  MEDIA_DIR_NAME .. "/" .. src
                return "<figure><video controls><source src=\"" .. video_src .. "\"></video><figcaption>" .. alt .. "</figcaption></figure>"
            end
        end

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
    return string:gsub("[^!]%[%[[^%]]+%]%]", function(match)
        local parts = utils.split(match:sub(4, -3), "|")
        local linked_name = parts[1]
        local display_name = parts[2] or linked_name
        local e = utils.get_key_case_insensitive(entries, linked_name)
        if e == nil then
            print("Warning, linked entry \"" .. linked_name .. "\" not found when rendering \"" .. entry.name .. "\"")
            return "{" .. display_name .. "}"
        end

        return match:sub(1, 1) .. "<a href=\"" .. e.dest_file_name .. "\">{" .. display_name .. "}</a>"
    end)
end

local function render_body(entry, entries)
    local title_html = entry.name == SITE_NAME and "" or "<h1>{{EntryName}}</h1>"
    return process_images(process_internal_links(sub_entry_fields(title_html ..
            (entry.date ~= nil and
                "<div style='color:#ccc'>last updated {{EntryDate}}</div>" or
                "") .. "{{EntryBodyHtml}}", entry), entry,
        entries))
end

local function render_entry(entry, entries)
    local html = string.format(
        "<!doctype html><html>%s<body><div class=\"layout\">%s<div class=\"content\"><main id=\"entry-body\">%s</main><p style=\"color:#ccc\"><em>Compiled %s</em></p></div></div>%s</body></html>",
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

-- expects an html body and will embed files with the ![[name]] syntax
local function process_embeds(body_html, entries)
    return body_html:gsub("!%[%[([^%]]+)%]%]", function(embedded_entry_name)
        for _, entry in pairs(entries) do
            if entry.name == embedded_entry_name then
                return "<article><h2>" ..
                    entry.name ..
                    "</h2><a href=" .. entry.dest_file_name .. ">#</a>" .. markdown(entry.body_raw) .. "</article>"
            end
        end
    end)
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

    local is_index = name == INDEX_NAME
    if is_index then
        name = parent_name
        if parent_name == SITE_NAME then
            -- root node
            dest_file_name = "index.html"
        else
            -- subtree root nodes
            dest_file_name = utils.slugify(name) .. ".html"
            parent_name = parts[#parts - 2] or SITE_NAME
        end
    end

    local body = utils.read_file(DATA_DIR .. file_path)

    local date
    local date_start, date_end = body:find("%d%d%d%d%-%d%d%-%d%d")
    if date_start == 1 then
        -- Date found at start of file content
        date = body:sub(date_start, date_end)
        body = body:sub(date_end + 1)
    elseif name:match("^%d%d%d%d%-%d%d%-%d%d$") then
        -- Filename itself is a date (e.g., 2023-01-05.md)
        date = name
    end

    local html = markdown(body)
    entries[name] = {
        name = name,
        src_path = file_path,
        dest_file_name = dest_file_name,
        parent_name = parent_name,
        body_raw = body,
        body_html = html,
        date = date,
        is_index = is_index
    }

    ::continue::
end

local i = 0
for _, entry in pairs(entries) do
    entry.body_html = process_embeds(entry.body_html, entries)
    render_entry(entry, entries)
    i = i + 1
end

local rss_entries = {}
for _, e in pairs(entries) do
    if (e.parent_name == "Writing" or e.parent_name == 'Log') and e.date ~= nil then
        table.insert(rss_entries, e)
    end
end
render_rss(rss_entries)

for _, media_path in pairs(media_paths) do
    local path = SITE_DIR .. "/" .. media_path
    if not utils.file_exists(path) then
        print("referenced media not found", path, " ... attempting to process")
        utils.process_image(DATA_DIR, SITE_DIR, path:gsub("-720", ""):gsub("-360", ""):gsub(SITE_DIR, DATA_DIR))
    end
end

print("Rendered " .. i .. " entries.")
