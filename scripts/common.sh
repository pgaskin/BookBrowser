#!/bin/bash
set -e

indent() {
  while read l; do echo "    $l"; done
}

# Change to project dir
cd ..