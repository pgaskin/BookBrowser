#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"
cd Content
cd Books
find . -name "*.txt" -type f -delete
find . -name "*.png" -type f -delete
find . -name "*.jpg" -type f -delete
find . -name "*.jpeg" -type f -delete
