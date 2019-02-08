#!/bin/sh

# usage: ./reset.sh DIRS...

for dir in $@
do
  echo "Discard changes of $dir..."
  git -C $dir checkout -- .
done

