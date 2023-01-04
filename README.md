# MMX

For information about mmx, see it in action at https://mrshll.com/mmx.html

Previous version with different feature set, written in Go now lives [here](https://github.com/mrshll/mmx-go) for archival purposes.

## Dependencies

- Lua
- GNU `date` for generating RSS item dates
- `find`
- _Optional:_ `inotifywatch` via `inotify-tools` for live site recompiling.

## Running

There are two scripts:
1. `build.sh` which compiles the site
1. `run.sh` which recompiles the site when the src or data changes
