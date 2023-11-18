#!/bin/bash
# This script uses imagemagick to convert images in a folder (and it's
# subfolders) into smaller sized variants.

# This script is a remixed version of Clemens Scott's
# https://git.sr.ht/~rostiger/batchResize/ used on site https://nchrs.xyz

# Use with care and backup your images!
# -----------------------------------------------------------------------------
CONTENT_DIR=$1
SITE_DIR=$2

if [ -z $CONTENT_DIR ] || [ -z $SITE_DIR ]; then
  echo "Usage: ./build.sh path/to/content path/to/site"
  exit 1
fi

# path where the original images are located
SRC="$CONTENT_DIR/media"
# path where the images will be stored
DST="$SITE_DIR/media"
mkdir $DST

# sizes to convert to
SIZES=(360 720)
MAXWIDTH=2400
#dithering
COLORS=16

function resize() {

  # Security check to prevent an endless loop when
  # $DST is inside $SRC (don't do that!)
  [[ $file == *"${DST}/*"* ]] && continue

  for file in $1; do
    # skip file if it doesn't exist
    [ -f "$file" ] || continue

    path=$(dirname $file)  # just/the/path
    name=$(basename $file) # filename.ext
    fileBase="${name%%.*}" # filename
    fileExt="${name#*.}"   # ext

    if ! grep -q -r "$name" $DST/../; then
      echo "${file} not used, skipping"
      echo "${file}" >>unused_media.txt
      continue
    fi

    # substitute source path with destination path
    # ${firstString/pattern/secondString}
    dst="${path/$SRC/$DST}"

    # existing images are skipped (delete images if they were updated)
    # create the output path (and parents) if it doesn't exist
    if [[ ! -d "$dst" ]]; then
      mkdir -p $dst
    fi

    # copy the file as is if it doesn't have the right extension
    if [[ "$fileExt" != "jpg" && "$fileExt" != "jpeg" && "$fileExt" != "png" && "$fileExt" != "MP.jpg" ]]; then
      cp -r $file $dst
      echo "Copied ${file}"
      continue
    fi

    # create smaller sizes for responsive image selection
    echo $file
    # get the width of the image
    width=$(identify -format "%w" "$file") >/dev/null
    for size in "${SIZES[@]}"; do
      # define output path and file
      output="$dst/$fileBase-${size}.${fileExt}"
      if [[ ! -f $output ]]; then
        # resize only  if original image is greater than or equal to (ge) the current size
        if [[ $width -ge $size ]]; then
          echo -n "| ${size} "
          convert $file -strip -auto-orient -resize $size -dither FloydSteinberg -colors $COLORS $output
        else
          #dither only
          echo -n "| ${width} "
          convert $file -strip -auto-orient -dither FloydSteinberg -colors $COLORS $output
        fi
      else
        echo -n "| ----- "
      fi
    done

    # Finally also strip the original image of it's EXIF data
    # and resize it to a max width of 1200
    output="$dst/$name"
    if [[ ! -f $output ]]; then
      if [[ $width -gt $MAXWIDTH ]]; then
        convert $file -strip -auto-orient -resize $MAXWIDTH $output
        echo -n "| ${MAXWIDTH} "
      else
        convert $file -strip -auto-orient $output
        echo -n "| ${width} "
      fi
    else
      echo -n "| ----- "
    fi
    echo -en "|\n"
  done
}

echo "" >unused_media.txt
# find all file in the source folder and run resize() on each
find $SRC | while read file; do resize "${file}"; done
