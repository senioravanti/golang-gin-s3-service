#!/usr/bin/env bash
set -eu

gen_url_safe_password() {
	local PASSWORD_LENGTH="$1"
	openssl rand -base64 "$PASSWORD_LENGTH" | tr '+/' '-_' | tr -d '='
}

CONFIG_DIR=./configs
CONFIG_FILES=(
  "${CONFIG_DIR}/s3-service.yaml"
)

S3_HOST=senioravanti.ru
S3_PORT=9090

S3_SERVICE_HOST=senioravanti.ru
S3_SERVICE_PORT=60221

TLS_CERT_DIR="/etc/letsencrypt/live/$S3_SERVICE_HOST"

while [ $# -gt 0 ]; do
  case "$1" in
    --s3-host) S3_HOST="$2"; shift 2;;
    --s3-port) S3_PORT="$2"; shift 2;;
    -h|--host) S3_SERVICE_HOST; shift 2;;
    -p|--port) S3_SERVICE_PORT="$2"; shift 2;;
    --tls-cert-dir) TLS_CERT_DIR="$2"; shift 2;; 
    *) echo 'unknown argument'; exit 1;;
  esac
done

create_config_file() {
  case "$1" in 
    "${CONFIG_FILES[0]}")
    cat <<-EOL > "${CONFIG_FILES[0]}"
		server:
		  read-timeout: 7s
		  write-timeout: 7s
		  port: "$S3_SERVICE_PORT"
		  tls:
		    cert: "${TLS_CERT_DIR}/fullchain.pem"
		    key: "${TLS_CERT_DIR}/privkey.pem"

		s3:
		  access-key: "$S3_ACCESS_KEY"
		  secret-key: "$S3_SECRET_KEY"
		  url: "https://${S3_HOST}:$S3_PORT"

		app:
		  log-level: DEBUG
EOL
    ;;
    *) 'unknown file name' ;;
  esac
}

if [ ! -d "$CONFIG_DIR" ]; then
  mkdir "$CONFIG_DIR"
fi

clear
for CONFIG_FILE in "${CONFIG_FILES[@]}"; do
  if [ ! -f "$CONFIG_FILE" ]; then
    echo "creating \`${CONFIG_FILE}\` ..."
    create_config_file "$CONFIG_FILE"
  fi
done
