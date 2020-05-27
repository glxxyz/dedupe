#!/bin/zsh

if [ $VENDOR = apple ]; then
    DEDUPE="./bin/macos-amd64/dedupe"
elif [ $VENDOR = debian ]; then
    DEDUPE="./bin/linux-amd64/dedupe"
else
    # TODO distinguish between Windows and other Linux distros, for now assume Linux
    DEDUPE="./bin/linux-amd64/dedupe"
fi
echo "Testing binary $DEDUPE"

TMPDIR=$(mktemp -d 2>/dev/null || mktemp -d -t "mytmpdir")
echo "Running tests in $TMPDIR"

SUCCESS=1
TEST="$TMPDIR/test"
TRASH="$TMPDIR/trash"
mkdir "$TRASH"
cp -rp ./test "$TEST"

assertRunSucceeded () {
    if [ $? -ne 0 ]; then
        echo "Run command failed"
        SUCCESS=0
    fi
}

assertDirExists() {
    if [ ! -d "$1" ]; then
        echo "Assertion failed: directory $1 doesn't exist"
        SUCCESS=0
    fi
}

assertFileExists() {
    if [ ! -f "$1" ]; then
        echo "Assertion failed: file $1 doesn't exist"
        SUCCESS=0
    fi
}

assertFileDoesntExist() {
    if [ -f "$1" ]; then
        echo "Assertion failed: $1 exists"
        SUCCESS=0
    fi
}

assertFileLines() {
    if [ $(cat "$1" | wc -l) -ne "$2" ]; then
        echo "Assertion failed: File $1 has $(cat $1 | wc -l) lines, expected $2 lines"
        SUCCESS=0
    fi
}

assertFileContains() {
    grep -qE "$2" "$1"
    if [ $? -ne 0 ]; then
        echo "Assertion failed: File $1 doesn't contain '$2'"
        SUCCESS=0
    fi
}

echo "Verify test preconditions"
assertDirExists "$TEST/baz/fooDirLink"
assertFileExists "$TEST/baz/LongFileA.txt"
assertFileExists "$TEST/baz/LongFileB.txt"
assertFileExists "$TEST/baz/ShortFileA.txt"
assertFileExists "$TEST/baz/ShortFileB.txt"
assertFileExists "$TEST/foo/bar/LongFileB.txt"
assertFileExists "$TEST/foo/bar/ShortFileALink"
assertFileExists "$TEST/foo/bar/ShortFileB.txt"
assertFileExists "$TEST/foo/LongFileA differentName.txt"
assertFileExists "$TEST/foo/ShortFileA.txt"
assertFileExists "$TEST/foo/ShortFileB.txt"

OUTPUT="$TMPDIR/test1.txt"
echo "Running test 1, compare all, with output $OUTPUT"
$DEDUPE "$TEST/foo" "$TEST/baz" > "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 10
# TODO: add asserts here

OUTPUT="$TMPDIR/test2.txt"
echo "Running test 2, min size, with output $OUTPUT"
$DEDUPE --min-size=100 "$TEST/foo" "$TEST/baz" > "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 4
# TODO: add asserts here

OUTPUT="$TMPDIR/test3.txt"
echo "Running test 3, min size compare name, with output $OUTPUT"
$DEDUPE --min-size=100 --compare-name "$TEST/foo" "$TEST/baz" > "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 2
# TODO: add asserts here

OUTPUT="$TMPDIR/test4.txt"
echo "Running test 4, no dupes, with output $OUTPUT"
$DEDUPE "$TEST/baz" >  "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 0

OUTPUT="$TMPDIR/test5.txt"
echo "Running test 5, don't follow directory symlink, with output $OUTPUT"
$DEDUPE "$TEST/baz" >  "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 0

OUTPUT="$TMPDIR/test6.txt"
echo "Running test 6, follow directory symlink, with output $OUTPUT"
$DEDUPE --follow-symlinks "$TEST/baz" > "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 12
# TODO: add asserts here

OUTPUT="$TMPDIR/test7.txt"
echo "Running test 7, don't follow file symlink, with output $OUTPUT"
$DEDUPE "$TEST/foo" > "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 2
assertFileContains "$OUTPUT" "Dupe:\t$TEST/foo/ShortFileB.txt\t$TEST/foo/bar/ShortFileB.txt"
assertFileContains "$OUTPUT" "Move:\t$TEST/foo/bar/ShortFileB.txt"

OUTPUT="$TMPDIR/test8.txt"
echo "Running test 8, follow file symlink, with output $OUTPUT"
$DEDUPE --follow-symlinks  "$TEST/foo" > "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 4
assertFileContains "$OUTPUT" "Dupe:\t$TEST/foo/ShortFileA.txt\t.+$TEST/baz/ShortFileA.txt"
assertFileContains "$OUTPUT" "Dupe:\t$TEST/foo/ShortFileB.txt\t$TEST/foo/bar/ShortFileB.txt"
assertFileContains "$OUTPUT" "Move:\t.+$TEST/baz/ShortFileA.txt"
assertFileContains "$OUTPUT" "Move:\t$TEST/foo/bar/ShortFileB.txt"

OUTPUT="$TMPDIR/test9.txt"
echo "Running test 9, move test, with output $OUTPUT"
$DEDUPE --trash="$TRASH" "$TEST/foo/bar" "$TEST/baz" "$TEST/foo" > "$OUTPUT"
assertRunSucceeded
assertFileLines "$OUTPUT" 10
assertFileContains "$OUTPUT" "Dupe:\t$TEST/foo/bar/LongFileB.txt\t$TEST/baz/LongFileB.txt"
assertFileContains "$OUTPUT" "Dupe:\t$TEST/foo/bar/ShortFileB.txt\t$TEST/.+/ShortFileB.txt"
assertFileContains "$OUTPUT" "Dupe:\t$TEST/.+ShortFileB.txt\t$TEST/foo/ShortFileB.txt"
assertFileContains "$OUTPUT" "Dupe:\t$TEST/baz/ShortFileA.txt\t$TEST/foo/ShortFileA.txt"
assertFileContains "$OUTPUT" "Dupe:\t$TEST/baz/LongFileA.txt\t$TEST/foo/LongFileA.+differentName.txt"
assertFileContains "$OUTPUT" "Move:\t$TEST/baz/ShortFileB.txt\t$TRASH$TEST/baz/ShortFileB.txt"
assertFileContains "$OUTPUT" "Move:\t$TEST/baz/LongFileB.txt\t$TRASH$TEST/baz/LongFileB.txt"
assertFileContains "$OUTPUT" "Move:\t$TEST/foo/LongFileA.+differentName.txt\t$TRASH$TEST/foo/LongFileA.+differentName.txt"
assertFileContains "$OUTPUT" "Move:\t$TEST/foo/ShortFileA.txt\t$TRASH$TEST/foo/ShortFileA.txt"
assertFileContains "$OUTPUT" "Move:\t$TEST/foo/ShortFileB.txt\t$TRASH$TEST/foo/ShortFileB.txt"

assertDirExists "$TEST/baz/fooDirLink"
assertFileExists "$TEST/baz/LongFileA.txt"
assertFileExists "$TEST/baz/ShortFileA.txt"
assertFileExists "$TEST/foo/bar/LongFileB.txt"
assertFileExists "$TEST/foo/bar/ShortFileB.txt"
assertFileExists "$TEST/foo/bar/ShortFileALink"

assertFileDoesntExist "$TEST/baz/LongFileB.txt"
assertFileDoesntExist "$TEST/baz/ShortFileB.txt"
assertFileDoesntExist "$TEST/foo/LongFileA differentName.txt"
assertFileDoesntExist "$TEST/foo/ShortFileA.txt"
assertFileDoesntExist "$TEST/foo/ShortFileB.txt"

assertFileExists "$TRASH/$TEST/baz/LongFileB.txt"
assertFileExists "$TRASH/$TEST/baz/ShortFileB.txt"
assertFileExists "$TRASH/$TEST/foo/LongFileA differentName.txt"
assertFileExists "$TRASH/$TEST/foo/ShortFileA.txt"
assertFileExists "$TRASH/$TEST/foo/ShortFileB.txt"

assertFileDoesntExist "$TRASH/$TEST/baz/fooDirLink"
assertFileDoesntExist "$TRASH/$TEST/baz/LongFileA.txt"
assertFileDoesntExist "$TRASH/$TEST/baz/ShortFileA.txt"
assertFileDoesntExist "$TRASH/$TEST/foo/bar/LongFileB.txt"
assertFileDoesntExist "$TRASH/$TEST/foo/bar/ShortFileB.txt"
assertFileDoesntExist "$TRASH/$TEST/foo/bar/ShortFileALink"

if [ $SUCCESS -eq 1 ]; then
    echo "All tests passed!"
else
    echo "Tests failed, see output for details!"
fi
