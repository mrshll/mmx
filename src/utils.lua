local utils = {}


-- file io

function utils.read_file(path)
  local file = io.open(path, "rb") -- r read mode and b binary mode
  if not file then error("unable to open file for reading at " .. path) end
  local content = file:read "*a" -- *a or *all reads the whole file
  file:close()
  return content
end

function utils.write_file(path, content)
  local file = io.open(path, "w")
  if not file then error("unable to open file for writing at " .. path) end
  file:write(content)
  file:close()
end

function utils.list_files(directory, extension)
  local pfile = io.popen('find ' .. directory .. ' -type f -name "*' .. extension .. '"')
  if not pfile then
    error("unable to list files in directory " .. directory .. " with extension " .. extension)
  end

  local file_names = {}
  for file_name in pfile:lines() do
    table.insert(file_names, file_name)
  end

  pfile:close()
  return file_names
end

-- strings

function utils.starts_with(str, start)
  return str:sub(1, #start) == start
end

function utils.ends_with(str, ending)
  return ending == "" or str:sub(- #ending) == ending
end

function utils.lines(str)
  if str:sub(-1) ~= "\n" then str = str .. "\n" end
  return str:gmatch("(.-)\n")
end

function utils.split(input, sep)
  if sep == nil then
    sep = "%s"
  end
  local t = {}
  for str in string.gmatch(input, "([^" .. sep .. "]+)") do
    table.insert(t, str)
  end
  return t
end

function utils.capitalize(str)
  return (str:gsub("^%l", string.upper))
end

function utils.title_case(str)
  return utils.capitalize(str:gsub(" %l", string.upper))
end

-- tables

function utils.has_keys(table, keys)
  for _, key in pairs(keys) do
    if table[key] == nil then
      return false
    end
  end

  return true
end

function utils.get_key_case_insensitive(table, key)
  return table[key] or table[utils.capitalize(key)] or table[key:lower()] or table[key:upper()] or
      table[utils.title_case(key)]
end

function utils.spairs(t, order)
  -- collect the keys
  local keys = {}
  for k in pairs(t) do keys[#keys + 1] = k end

  -- if order function given, sort by it by passing the table and keys a, b,
  -- otherwise just sort the keys
  if order then
    table.sort(keys, function(a, b) return order(t, a, b) end)
  else
    table.sort(keys)
  end

  -- return the iterator function
  local i = 0
  return function()
    i = i + 1
    if keys[i] then
      return keys[i], t[keys[i]]
    end
  end
end

function utils.dump(o)
  if type(o) == 'table' then
    local s = '{ '
    for k, v in pairs(o) do
      if type(k) ~= 'number' then k = '"' .. k .. '"' end
      s = s .. '[' .. k .. '] = ' .. utils.dump(v) .. ','
    end
    return s .. '} '
  else
    return tostring(o)
  end
end

-- dates
function utils.today()
  local date_table = os.date("*t")
  local year, month, day = date_table.year, date_table.month, date_table.day -- date_table.wday to date_table.day
  return string.format("%d-%d-%d", year, month, day)
end

return utils
