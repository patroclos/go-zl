#!/bin/bash
tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT
cat > "$tmp"
bpmn-to-image "$tmp:/tmp/bpmn.svg" >/dev/null
cat /tmp/bpmn.svg
