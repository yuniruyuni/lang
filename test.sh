#/bin/bash

TMPDIR=".tmp"
OUTPUT="${TMPDIR}/tmp.ll"
TARGET="bin/lang"

test() {
    args="$1"
    want="$2"

    mkdir -p "${TMPDIR}"
    echo "$args" | $TARGET > "${OUTPUT}"
    got=`lli ${OUTPUT}`

    if [ "$got" == "$want" ]; then
        echo "[SUCCEED] $args => $got"
    else
        echo "[FAILED] $args => want: $want, got: $got"
    fi
}

test "" "test"
