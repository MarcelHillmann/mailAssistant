#!/bin/bash
function reset(){
    echo -e -n "\e[0m"
}

function header() {
    tab=${2:-0}
    echo -e -n "\e[33m $1"

    for (( t=1; t <$tab; t++ ))
    do
      echo -en "\t"
    done
    reset
}

function go_meta_linter() {
    reset
    header "run golangci-lint" 4
    golangci-lint run --out-format checkstyle ./... > .sonarqube/linter-report.xml 2>&1
    showState $?
    reset
}

function go_lint(){
    reset
    header "run golint" 5
    golint ./... > .sonarqube/golint-report.out 2>&1
    showState $?
    reset
}

function go_vet(){
    reset
    header "run go vet" 5
    go vet ./... > .sonarqube/govet-report.out 2>&1
    showState $?
    reset
}

function go_test_json(){
    reset
    header "run tests with result as TEXT" 2
    go test -json -covermode=count -coverprofile=.sonarqube/cover.out ./... > .sonarqube/test-report.json
    showState $? 1
    reset
}

function go_junit() {
    reset
    if [[ -f .sonarqube/test-report.json ]]; then
        header "run go-junit-report" 4
        cat .sonarqube/test-report.json | go-junit-report -package-name "${PKG}" -set-exit-code > .sonarqube/test.xml
        showState $?
    fi
    reset
}

function showState(){
    local rc=$1
    local match=${2:-0}

    if [[ ${rc} = ${match} ]]
    then
        echo -e "\e[32m PASSED\e[0m"
    else
        echo -e "\e[91m Failed\e[0m"
    fi # failed
}

go clean -testcache
go_meta_linter
go_lint
go_vet
go_test_json
go_junit
exit 0