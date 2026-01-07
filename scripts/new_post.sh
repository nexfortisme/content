#!/bin/sh

POSTS_DIR="posts"

if [ -z "$1" ]; then
  echo "Usage: $0 <filename.md>"
  exit 1
fi

FILENAME="$1"
OUTPUT_PATH="$POSTS_DIR/$FILENAME"

# Ensure posts directory exists
mkdir -p "$POSTS_DIR" || {
  echo "Failed to create directory: $POSTS_DIR"
  exit 1
}

cat <<'EOF' > "$OUTPUT_PATH"
----
title: <POST TITLE>
description: <POST DESCRIPTION>
descriptionImage: <POST DESCRIPTION IMAGE> | null
tags: [<TAGS>]
----
EOF

echo "Created $OUTPUT_PATH"
