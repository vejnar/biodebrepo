#!/bin/bash

set -e

if [ ! -d "repo" ]; then
    echo "ERROR: repo not found"
    exit 1
fi

current_tag=$(gh release list --json tagName | jq -r '.[0].tagName')
next_tag=$(( current_tag + 1 ))

echo "Creating release $next_tag in 5 seconds..."
sleep 5

gh release create $next_tag --notes "Release $next_tag" ./repo/*
