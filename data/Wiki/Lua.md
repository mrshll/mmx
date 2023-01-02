## Tables
### Sorting

Create an iterator

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
	    if keys![](i) then
	      return keys![t[keys[i](i])]
	    end
	  end
	end

Example usage (from [[mmx]]):

	local function make_sorted_entry_iterator(entries)
	  return utils.spairs(entries, function(es, name_a, name_b)
	    local entry_a = es![](name_a)
	    local entry_b = es![](name_b)
	    if entry_a.date or entry_b.date then
	      return (entry_a.date or MIN_DATE) > (entry_b.date or MIN_DATE)
	    else
	      return name_a < name_b
	    end
	  end)
	end

