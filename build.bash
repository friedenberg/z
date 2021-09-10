#! /bin/bash -e

indexSum=$(sha1sum lib/index.go | cut -d' ' -f1)

go build -ldflags "-X lib.thisFileSha1=$indexSum" ./.
