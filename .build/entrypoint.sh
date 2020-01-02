#!/usr/bin/env bash
if [ "$1" == "" ]
then
   mailAssistant run
else
   mailAssitant $@
fi