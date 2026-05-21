#!/bin/bash

dist="dist"
rm -rf ${dist}
mkdir -p ${dist}

build_single() {
  GOOS=${1} GOARCH=${2} go build -ldflags "-s -w" -o "${dist}/${3}_${1}_${2}${4}" "cmd/${3}/main.go"
}

oa_es=(linux:amd64 linux:arm64 linux:arm linux:mipsle darwin:amd64 darwin:arm64 windows:amd64:.exe)

build_single_4_platforms() {
  for oa in "${oa_es[@]}"
  do
    # shellcheck disable=SC2206
    oa_s=(${oa//:/ })
    echo build_single "${oa_s[0]}" "${oa_s[1]}" "$1" "${oa_s[2]}"
    build_single "${oa_s[0]}" "${oa_s[1]}" "$1" "${oa_s[2]}"
  done
}


build_single_4_platforms notifier

