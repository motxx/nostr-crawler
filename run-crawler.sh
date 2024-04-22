#!/bin/bash

if [ -z "$ROOT" ]; then
  echo "Error: ROOT environment variable is not set."
  exit 1
fi

BIN="${ROOT}/nostr-crawler"

$BIN --from wss://relay.snort.social --to ws://localhost:7447
