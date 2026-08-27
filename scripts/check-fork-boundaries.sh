#!/bin/sh
set -eu

fail() {
    echo "fork-boundary check failed: $*" >&2
    exit 1
}

grep -Fq 'owner: thanhnam2811' .goreleaser.yml \
    || fail '.goreleaser.yml must publish to thanhnam2811'

grep -Fq 'name: ainovel-cli' .goreleaser.yml \
    || fail '.goreleaser.yml release repository name changed unexpectedly'

grep -Fq 'REPO="thanhnam2811/ainovel-cli"' scripts/install.sh \
    || fail 'installer must download releases from thanhnam2811/ainovel-cli'

grep -Fq 'updateRepo = "thanhnam2811/ainovel-cli"' cmd/ainovel-cli/main.go \
    || fail 'self-update must query thanhnam2811/ainovel-cli'

if grep -Fq 'ghcr.io/voocel/ainovel-cli' docker-compose.yml; then
    fail 'docker-compose must not pull the upstream image'
fi

if grep -Fq 'owner: voocel' .goreleaser.yml; then
    fail 'GoReleaser still points at upstream owner'
fi

echo 'fork-boundary checks passed'
