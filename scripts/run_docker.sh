#!/usr/bin/env bash

clear
set -a; . "$ENV_FILE"; set +a
docker run -d -p "${SERVER_PORT}:8080" \
  --name 'credgen' --restart 'on-failure:3' \
  "stradiavanti/s3-service:$S3_SERVICE_TAG"