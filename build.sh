#!/usr/bin/env bash

FILE="lamb_oracle.tar.gz"
if PWD=$(cd "$(dirname "$0")"; pwd); then
  echo "Directory: $PWD"
else
  echo "Error: Failed to get PWD." >&2
fi

rm -fr "$PWD"/_build

if [ "clean" == "$1" ]; then
  exit 0
fi

mkdir -p "$PWD"/_build/lamb_oracle/config

printf "# lns_lamb_oracle\r\nLAMB price oracle application of Lambda Name Service." >>"$PWD"/_build/lamb_oracle/README.md

GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -tags=jsoniter -o "$PWD"/_build/lamb_oracle/lamb_oracle "$PWD"

cp -f "$PWD"/config/config.yaml "$PWD"/_build/lamb_oracle/config

cd "$PWD"/_build && tar cvzf "$FILE" lamb_oracle

exit 0
