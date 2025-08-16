#!/usr/bin/env bash

TAG='0.1'
PLATFORM='linux/amd64'
IS_PUSH=''
IS_NO_CACHE=''

while [ $# -gt 0 ]; do
  case "$1" in
    --tag) TAG="$2"; shift 2;;
    --platform) PLATFORM="$2"; shift 2;;
    --no-cache) IS_NO_CACHE='--no-cache'; shift;;
    --push) IS_PUSH='--push'; shift;;

    *) exitWithMsg 'build_docker.sh' 'unknown argument' ;;
  esac
done

clear
set -a; . "$ENV_FILE"; set +a
docker build . \
  -f ./build/package/docker/Dockerfile \
  --platform "$PLATFORM" \
  -t "stradiavanti/s3-service:$TAG" \
  $IS_NO_CACHE $IS_PUSH
