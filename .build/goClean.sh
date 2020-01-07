#!/bin/bash
if [[ "${1}" = "develop" ]]
then
    go clean -cache -testcache
fi
exit 0