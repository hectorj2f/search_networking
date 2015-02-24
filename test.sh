#!/bin/bash -e

echo "Starting golynt..."

source ./build

TESTABLE_AND_FORMATTABLE=""
FORMATTABLE="$TESTABLE_AND_FORMATTABLE api database networking resources server types version"


TEST=$TESTABLE_AND_FORMATTABLE
FMT=$FORMATTABLE


# split TEST into an array and prepend REPO_PATH to each local package
split=(${TEST// / })
TEST=${split[@]/#/${REPO_PATH}/}

echo "Running tests..."
go test -timeout 60s $@ ${TEST} --race

echo "Checking gofmt..."
fmtRes=$(gofmt -l $FMT)
if [ -n "${fmtRes}" ]; then
  echo -e "gofmt checking failed:\n${fmtRes}"
  exit 255
fi

echo "Checking govet..."
vetRes=$(go vet $TEST)
if [ -n "${vetRes}" ]; then
  echo -e "govet checking failed:\n${vetRes}"
  exit 255
fi

echo "Success :D"
