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

test '"test"' 'test'
test '"日本語"' '日本語'
test ' "日本語" ' '日本語'
test '"日本語" ' '日本語'
test ' "日本語"' '日本語'
test '1' '1'
test '1 + 1' '2'
test '10 + 1 + 100' '111'
test '3 - 2' '1'
test '2 * 3' '6'
test '6 / 2' '3'
test '3 * 2 + 1' '7'
test '6 / 2 + 1' '4'
test '1 + 3 * 2' '7'
test '1 + 6 / 2' '4'
test '3 * (2 + 1)' '9'
test '(2 + 1) * 3' '9'
test '6 / (2 + 1)' '2'
