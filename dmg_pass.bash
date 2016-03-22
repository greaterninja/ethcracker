#!/bin/sh

i=0

while IFS='' read -r line || [[ -n "$line" ]]; do
    
    ((i++))

    echo "$i: $line"
    
    res="$(printf $line | hdiutil attach -quiet -stdinpass $2 )" 

#    echo "RESULT: $?"
    if [ $? -eq 0 ]; then
    
        echo "================================================================"
        echo ""
        echo "            The disk is succesfully mounted!!! "
        echo "            Your password is: $line"
        echo ""
        echo "================================================================"
        exit    
    fi
    
done < "$1"