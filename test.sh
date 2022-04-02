#/bin/bash

TMPDIR=".tmp"
OUTPUT="${TMPDIR}/tmp.ll"
TARGET="bin/lang"

test_with() {
    file="$1"
    want="$2"

    cat "$file"

    mkdir -p "${TMPDIR}"
    cat "$file" | $TARGET > "${OUTPUT}"
    got=`lli ${OUTPUT}`

    if [ "$got" == "$want" ]; then
        echo "[SUCCEED] $args => $got"
    else
        echo "[FAILED] $args => want: $want, got: $got"
    fi
}

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

fail() {
    args="$1"
    want="$2"

    mkdir -p "${TMPDIR}"
    echo "$args" | $TARGET 2> "${OUTPUT}"
    got=`cat ${OUTPUT}`

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
test '2 * 3 * 4' '24'
test '2 < 4' '1'
test '2 < 2' '0'
test '2 == 2' '1'
test '2 == 4' '0'
test 'if 1 { 10 } else { 20 }' '10'
test 'if 0 { 10 } else { 20 }' '20'
test 'if 1 < 2 { 10 } else { 20 }' '10'
test 'if 2 < 1 { 10 } else { 20 }' '20'

test 'if 1 { if 1 { 10 } else { 20 } } else { 30 }' '10'
test 'if 1 { if 0 { 10 } else { 20 } } else { 30 }' '20'
test 'if 0 { if 0 { 10 } else { 20 } } else { 30 }' '30'

test 'if 1 { 10 } else { if 0 { 20 } else { 30 } }' '10'
test 'if 0 { 10 } else { if 1 { 20 } else { 30 } }' '20'
test 'if 0 { 10 } else { if 0 { 20 } else { 30 } }' '30'

test_with 'test/if.yuni' '10'

fail 'if' 'failed to parse code: invalid tokens'
