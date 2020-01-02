#!/bin/bash
if [[ "$1" == "" ]]; then
   /usr/bin/mailAssistant run
else
   /usr/bin/mailAssistant $@
fi
exit $?