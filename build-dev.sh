#!/bin/bash
DATE=$(date '+%Y%m%d%H%M%S')
COMMIT=$(git rev-parse --short HEAD)
COUNT=$(git rev-list HEAD --count)
VERSION="dev-$DATE.$COUNT"

function build() {
  OS=$1
  ARCH=$2
  NAME="explorer-$VERSION-$COMMIT-$OS-$ARCH"
  LD_FLAG="-s -w"
  if [[ $OS == windows* ]]; then
      NAME="$NAME.exe"
      LD_FLAG="-H windowsgui $LD_FLAG"
  fi
  echo "build $NAME@$COMMIT for $OS"
  GOOS=$OS GOARCH=$ARCH go build -o "$NAME" -ldflags="$LD_FLAG" ./cmd/explorer
  tar -cvzf $NAME.tar.gz $NAME
  echo "build $NAME success"
}

build "linux" "amd64"
build "linux" "386"
build "darwin" "amd64"