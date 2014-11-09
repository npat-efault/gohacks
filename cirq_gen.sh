#!/bin/sh

if [ $# -ne 4 -a $# -ne 5 ]; then
    echo "Usage is: $(basename $0) <pkg> <qtype> <new> <el> [nil]" 1>&2
    echo "    <pkg> is the package name for the generated file" 1>&2
    echo "    <qtype> is the name for the queue type" 1>&2
    echo "    <new> is the name for the function returning a new queue" 1>&2
    echo "    <el> is the type for queue elements" 1>&2
    echo "    nil: if the element type has a nil zero value" 1>&2
    exit 1
fi 

package=$1
qtype=$2
new=$3
eltype=$4
if [ $# -eq 5 ]; then
   donil='s/__NIL/nil/g'
else
   donil='/__NIL/d'
fi

cat <<EOF
// Auto-generated. !! DO NOT EDIT !!

EOF
sed \
    -e "s/__PACKAGE/$package/g" \
    -e "s/__Q/$qtype/g" \
    -e "s/__NewQ/$new/g" \
    -e "s/__ELTYPE/$eltype/g" \
    -e "$donil" \
    cirq.gox | gofmt 
