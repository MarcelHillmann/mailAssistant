#!/bin/bash
branch=$1

if [[ "$branch" = "develop" ]]
then
    go clean -cache -testcache
fi
exit 0