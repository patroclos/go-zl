#!/bin/bash

tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT

cat > "$tmp"

/usr/local/bin/mmdc -b white -H 512 -w 512 -t neutral -q -i "$tmp" -p <(echo '{"executablePath": "/usr/bin/chromium-browser", "args":["--no-sandbox"]}')

cat "$tmp.svg"
