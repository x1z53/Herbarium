#!/bin/sh
# Should run from project root dir

sh po/update_potfiles.sh

# Create temporary POT file for Go files
grep '\.go$' ./po/POTFILES.in | xargs xgettext --language=C --keyword=T_ -o po/herbarium-go.pot --from-code=UTF-8 --add-comments --package-name=herbarium

# Create temporary POT file for desktop.in files
grep '\.desktop\.in.in$' ./po/POTFILES.in | xargs xgettext --language=Desktop -o po/herbarium-desktop.pot --from-code=UTF-8 --add-comments --package-name=herbarium

# Merge the POT files
msgcat po/herbarium-go.pot po/herbarium-desktop.pot -o po/herbarium.pot

# Clean up temporary files
rm -f po/herbarium-go.pot po/herbarium-desktop.pot
