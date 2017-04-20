#!/bin/bash

set -o pipefail

tags=$1
error=0

# # treat non-app and non gui packages
for pkg in $(retool do govendor list -p -no-status +local | tail -n +2 | grep -v '/gui/\?*'); do
	retool do gometalinter --config=.gometalinter.json $pkg 2>&1 | sed \$d
	test $? -eq 0 || error=1
done

# treat app and gui packages
for pkg in $(retool do govendor list -p -no-status +local | head -n 1 ; retool do govendor list -p -no-status +local | grep '/gui/\?*'); do
	retool do gometalinter --config=.gometalinter.json --disable="aligncheck" --disable="errcheck" \
	--disable="gosimple" --disable="interfacer" --disable="staticcheck" --disable="structcheck" \
	--disable="unconvert" --disable="unused" --disable="varcheck" $pkg 2>&1 | sed \$d
	test $? -eq 0 || error=1
done

exit $error
