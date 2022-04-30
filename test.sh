#/bin/bash

TMPDIR=".tmp"
OUTPUT="${TMPDIR}/tmp.ll"
TARGET="bin/lang"

test_with() {
    file="$1"
    want="$2"

    mkdir -p "${TMPDIR}"
    cat "$file" | $TARGET > "${OUTPUT}"
    got=`lli ${OUTPUT}`

    if [ "$got" == "$want" ]; then
        echo "[SUCCEED] $file => $got"
    else
        echo "[FAILED] $file => want: $want, got: $got"
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
        echo "[SUCCEED(for fail test)] $args => $got"
    else
        echo "[FAILED(for fail test)] $args => want: $want, got: $got"
    fi
}

test 'func main(){ printf("%d", 1,) }' '1'
test 'func main(){ printf("%d", 1 + 1,) }' '2'
test 'func main(){ printf("%d", 10 + 1 + 100,) }' '111'
test 'func main(){ printf("%d", 3 - 2,) }' '1'
test 'func main(){ printf("%d", 2 * 3,) }' '6'
test 'func main(){ printf("%d", 6 / 2,) }' '3'
test 'func main(){ printf("%d", 3 * 2 + 1,) }' '7'
test 'func main(){ printf("%d", 6 / 2 + 1,) }' '4'
test 'func main(){ printf("%d", 1 + 3 * 2,) }' '7'
test 'func main(){ printf("%d", 1 + 6 / 2,) }' '4'
test 'func main(){ printf("%d", 3 * (2 + 1),) }' '9'
test 'func main(){ printf("%d", (2 + 1) * 3,) }' '9'
test 'func main(){ printf("%d", 6 / (2 + 1),) }' '2'
test 'func main(){ printf("%d", 2 * 3 * 4,) }' '24'
test 'func main(){ printf("%d", 2 < 4,) }' '1'
test 'func main(){ printf("%d", 2 < 2,) }' '0'
test 'func main(){ printf("%d", 2 == 2,) }' '1'
test 'func main(){ printf("%d", 2 == 4,) }' '0'
test 'func main(){ printf("%d", if 1 { 10 } else { 20 },) }' '10'
test 'func main(){ printf("%d", if 0 { 10 } else { 20 },) }' '20'
test 'func main(){ printf("%d", if 1 < 2 { 10 } else { 20 },) }' '10'
test 'func main(){ printf("%d", if 2 < 1 { 10 } else { 20 },) }' '20'
test 'func main(){ printf("%d", if 1 { if 1 { 10 } else { 20 } } else { 30 },) }' '10'
test 'func main(){ printf("%d", if 1 { if 0 { 10 } else { 20 } } else { 30 },) }' '20'
test 'func main(){ printf("%d", if 0 { if 0 { 10 } else { 20 } } else { 30 },) }' '30'
test 'func main(){ printf("%d", if 1 { 10 } else { if 0 { 20 } else { 30 } },) }' '10'
test 'func main(){ printf("%d", if 0 { 10 } else { if 1 { 20 } else { 30 } },) }' '20'
test 'func main(){ printf("%d", if 0 { 10 } else { if 0 { 20 } else { 30 } },) }' '30'
test 'func main(){ 1; printf("%d",2,) }' '2'
test 'func main(){ 1; printf("%d",if 0 { 10 } else { 20 },) }' '20'
test 'func main(){ printf("%d", if 0 { 10 } else { 20; 30 },) }' '30'
test 'func main(){ printf("%d", if 0; 1 { 10 } else { 20; 30 },) }' '10'
test 'func main(){ let x = 10; printf("%d", x,) }' '10'
test 'func main(){ let x = 10; printf("%d", 20,) }' '20'
test 'func main(){ let x = 10; x = 20; printf("%d", x,) }' '20'
test 'func main(){ let x = 10; x = 20; printf("%d", x * 10,) }' '200'
test 'func main(){ printf("%d", 10,); 100 }' '10'

test_with 'test/if.yuni' '10'
test_with 'test/var-if.yuni' '100'
test_with 'test/large.yuni' '40'
test_with 'test/while.yuni' '45'
test_with 'test/fact.yuni' '362880'
fail 'if' 'failed to parse code: invalid tokens'
