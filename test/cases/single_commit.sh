#!/usr/bin/env bash

set -e
cd "$1"

git init
touch TEST
git add .
git commit -a -m "First commit."
