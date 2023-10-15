#!/bin/bash
if ! test -f lem-in; then 
go build -o lem-in .
fi

if ! test -f visualizer; then 
go build -o visualizer .
fi

if [ -z "$1" ]
then
    file="example01.txt"
else
    file="$1"
fi

./lem-in $file | ./visualizer