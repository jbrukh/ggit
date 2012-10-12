#!/usr/bin/env bash

set -e
cd "$1"

git init

for ITER in {1..10}; do
	echo $ITER >> $ITER
	git add $ITER
	git commit -a -m "Commit $ITER"
done