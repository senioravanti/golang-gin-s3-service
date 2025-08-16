#!/usr/bin/env bash

if [ ! -d "$(pwd)/scripts/environment/" ]; then
  mkdir "$(pwd)/scripts/environment/"
fi

if [ -z "$ENV_FILE" ]; then
  ENV_FILE="$(pwd)/scripts/environment/.env"
fi

if [ -f "$ENV_FILE" ]; then
  echo "file \`${ENV_FILE}\` -> already exists"
  exit 1
fi

echo "creating \`${ENV_FILE}\` ...\n"
cat <<EOL > "$ENV_FILE"
# s3
S3_ACCESS_KEY=''
S3_SECRET_KEY=''
S3_URL=''

# server
SERVER_PORT=8070
EOL