# MMX

For information about mmx, see it in action at https://mrshll.com/mmx.html

Previous version with different feature set, written in Go now lives [here](https://github.com/mrshll/mmx-go) for archival purposes.

## Dependencies

- Lua
- GNU `date` for generating RSS item dates
- `find`
- _Optional:_ `inotifywatch` via `inotify-tools` for live site recompiling.

## Running

There are three main scripts:
1. `./build.sh path/to/content path/to/site` which compiles the site
1. `./run.sh path/to/content path/to/site` which recompiles the site when the src or data changes
1. `./processImages.sh path/to/content path/to/site` which processes the images with resizing, compression and (optional) dithering
