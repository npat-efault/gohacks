#!/bin/bash

if [ $# -ne 5 ]; then
    echo "Usage is: $(basename $0) <outfile> <pkg> <qtype> <new> <el>" 1>&2
    echo "    <outfile> is the name of the generated file" 1>&2
    echo "    <pkg> is the package name for the generated file" 1>&2
    echo "    <qtype> is the name for the queue type" 1>&2
    echo "    <new> is the name for the function returning a new queue" 1>&2
    echo "    <el> is the type for queue elements" 1>&2
    exit 1
fi

outfile="$1"
package=$2
qtype=$3
new=$4
eltype=$5

gox="cirq.gox"
src="./$gox"
pkg="github.com/npat-efault/gohacks/cirq"
if [ ! -f "$src" ]; then
    src=$(go list -f '{{ .Dir }}' "$pkg")/"$gox"
    if [ ! -f "$src" ]; then
        echo "Cannot find \"$gox\" in . or in package \"$pkg\"" 1>&2
        exit 1
    fi
fi

set -e
trap 'rm -f "$outfile".tmp' EXIT

cat > "$outfile".tmp <<EOF
// Auto-generated. !! DO NOT EDIT !!

EOF
sed \
    -e "s/__PACKAGE/$package/g" \
    -e "s/__Q/$qtype/g" \
    -e "s/__NewQ/$new/g" \
    -e "s/__ELTYPE/$eltype/g" \
    "$src" | gofmt >> "$outfile".tmp
mv "$outfile".tmp "$outfile"
chmod -w "$outfile"

