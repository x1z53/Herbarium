#!/bin/sh
# Should run from project root dir

touch ./po/unsort-POTFILES

echo "data/ru.ximper.Herbarium.desktop.in.in" >> ./po/unsort-POTFILES
echo "data/ru.ximper.Herbarium.metainfo.xml.in.in" >> ./po/unsort-POTFILES

find ./ -iname "*.go" -type f -exec grep -lrE 'T_\(' {} + | while read file; do echo "${file#./}" >> ./po/unsort-POTFILES; done

cat ./po/unsort-POTFILES | sort | uniq > ./po/POTFILES.in

rm ./po/unsort-POTFILES
