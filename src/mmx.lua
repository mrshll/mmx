local markdown = require "markdown"
local utils = require "utils"

local INDEX_NAME = "Now"
local MIN_DATE = "0000-00-00"

local function sub_entry_fields(str, entry)
  return str:gsub("{{(%w+)}}", {
    ["EntryBodyHtml"] = entry.body_html,
    ["EntryName"] = entry.name,
    ["EntryDate"] = entry.date,
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
    i = i + 1
    local is_selected = entry ~= nil and e.name == entry.name
    acc = acc .. "<li>"

    if is_selected then acc = acc .. "<mark>" end
    acc = acc .. "<a href=\"" .. e.dest_path .. "\">" .. e.name .. "</a>"
    if is_selected then acc = acc .. "</mark>" end
    acc = acc .. "</li>"
    if i == 6 then
      acc = acc .. "<details><summary>See all</summary>"
    end
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
    if e.parent_name == entry.name and e.name ~= INDEX_NAME then
      table.insert(children, e)
    end
  end
  acc = acc .. render_nav_section(nil, children)

  while entry.name ~= INDEX_NAME do
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
    if parent == nil then error(entry.name .. " has no parent") end

    entry = parent
  end

  return "<nav>" .. render_nav_section(entry, { entry }) .. acc .. "</nav>"
end

local function render_body(entry, entries)
  return sub_entry_fields(
    "<h1>{{EntryName}}</h1>" ..
    (entry.date ~= nil and "<div style='color:#ccc'>last updated {{EntryDate}}</div>" or "") ..
    "{{EntryBodyHtml}}"
    , entry):gsub("%[%[([^%]]+)%]%]", function(match)
    local parts = utils.split(match, "|")
    local linked_name = parts[1]
    local display_name = parts[2] or linked_name
    local e = utils.get_key_case_insensitive(entries, linked_name)
    if e == nil then
      print("Warning, linked entry \"" .. linked_name .. "\" not found when rendering \"" .. entry.name .. "\"")
      return "{" .. display_name .. "}"
    end

    return "<a href=\"" .. e.dest_path .. "\">{" .. display_name .. "}</a>"
  end)
end

local function render_entry(entry, entries)
  local html = string.format("<html>%s<body><div class=\"content\">%s%s<p style=\"color:#ccc\"><em>Compiled %s</em></p>%s</body></html>"
    ,
    render_head(entry), render_nav(entry, entries), render_body(entry, entries), utils.today(), render_footer())
  utils.write_file(entry.dest_path, html)
end

local DATA_DIR = "../data"
local DATA_EXT = '.md'
local entries = {}
local file_paths = utils.list_files(DATA_DIR, DATA_EXT)
for _, file_path in pairs(file_paths) do

  -- remove the DATA_DIR and split on slash to get path fragments
  local parts = utils.split(file_path:sub(#DATA_DIR + 1), "/")

  -- directory scheme is /entry-name or [/parent's parent]/parent/entry-name (recursive)
  local parent_name = parts[#parts - 1] or INDEX_NAME
  -- remove the file extension
  local name = parts[#parts]:sub(0, -1 * #DATA_EXT - 1)
  local dest_path = '../docs/' .. name .. '.html'

  if name == parent_name then
    if name == INDEX_NAME then
      -- special case for the root node
      dest_path = '../docs/index.html'
    else
      -- special case for folder indeces
      parent_name = parts[#parts - 2] or INDEX_NAME
    end
  end

  local body = utils.read_file(file_path);
  local date
  local date_start, date_end = body:find('%d%d%d%d%-%d%d%-%d%d')
  if date_start == 1 then
    date = body:sub(date_start, date_end)
    body = body:sub(date_end + 1)
  end


  local html = markdown(body)
  entries[name] = {
    name = name,
    src_path = file_path,
    dest_path = dest_path,
    parent_name = parent_name,
    body_raw = body,
    body_html = html,
    date = date
  }

  -- local entry_acc = {}
  -- for line in lines(fileContent) do
  --   if starts_with(line, "name:") then
  --     entry_acc["name"] = line:sub(6)
  --   elseif starts_with(line, "host:") then
  --     entry_acc["host"] = line:sub(6)
  --   end
  --   if has_keys(entry_acc, { "name", "host" }) then
  --     table.insert(entries, entry_acc)
  --     entry_acc = {}
  --   end
  -- end
end

local i = 0
for _, entry in pairs(entries) do
  render_entry(entry, entries)
  i = i + 1
end

print("Rendered " .. i .. " entries.")
