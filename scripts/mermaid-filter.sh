#!/bin/bash

tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT

cat > "$tmp"

/usr/local/bin/mmdc -b white -H 512 -w 512 -t neutral -q -i "$tmp"

cat "$tmp.svg"
