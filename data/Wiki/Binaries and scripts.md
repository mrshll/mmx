## Binaries

### date
_print or set the system date and time_
* get dates for e.g. file names : `date -Iminutes`

### maim 
_maim (make image) takes screenshots of your desktop. It has options to take only a region, and relies on another program called slop to query the user for regions using the graphical interface._

- full sceenshot : `maim destination.png`
- current window screenshot : `maim -i (xdotool getactivewindow) destination.png`
- selected area screenshot : `maim -s destination.png`
- name output based on time : `maim (date -Iminutes).png`
- get color of selected pixel : `maim -st 0 | convert - -resize 1x1\! -format '%!!![]()()(pixel:p{0,0})' info:-`

### xclip
_command line interface to X selections (clipboard)_

- pipe screenshot to clipboard : `maim -s | xclip -selection clipboard --target image/png`
- pipe program output to clipboard : `head file.txt | xclip -selection clipbaord`

### tac
_reverse input line-wise_

### rev
_reverse input character-wise_

### didder
_didder is an extensive, fast, and accurate command-line image dithering tool._
[github](https://github.com/makeworld-the-better-one/didder)

### ffmpeg

Create animated gif from screen area
`ffmpeg -framerate 25 -f x11grab (slop -f '-video_size %wx%h -i :0.0+%x,%y' | string split ' ') ~/Videos/capture_(date -Iminutes).mkv`


## Scripts

### Find and replace using sed across all occurrances in a folder

    find . -type f -exec sed -i.txt "s/foo/bar/g" {} \;

### Rename file extension

    !#/bin/bash
    for f in *.JPG;
      do mv -- "$f" "${f%.JPG}.jpg";
    done

## Get paste-able filename list (used for {log} entries)

    ls -1 2021-08-*

### Migrating from .mmx to .mx
Images

    find . -type f -name "*.md" -exec gsed --regexp-extended -i.bak "s/\!!![](](([^]((\S*)()]*))?\)/!![)(\3)(\1)/g" {} \;

Links

    find . -type f -name "*.md" -exec gsed --regexp-extended -i.bak "s/\{(\S*), (!!![](]()(^\})*)\}/[\2)(\1)/g" {} \;
