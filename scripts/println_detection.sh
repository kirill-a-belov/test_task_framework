#!/bin/bash

# find fmt.Println leftovers
if [ $# != 1 ] ; then
  echo "usage: $0 path"
  exit 1
fi

RESULT=$(grep -Rnw $1 --exclude-dir=scripts -e 'fmt.Println')
if [[ $RESULT != "" ]]
then
  echo "fmt.Println detected: $RESULT"
  exit 1
fi
