#!/usr/bin/env bash

if [ -z "$ENV_FILE" ]; then
  ENV_FILE="$(pwd)/scripts/environment/.env"
fi

clear
set -a; . "$ENV_FILE"; set +a
go run "$(pwd)/cmd"
