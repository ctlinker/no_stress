#!/usr/bin/bash

if ! [[ -s "$(ls ./build/server)" ]] || [[ "$@" == "-f" ]]; then
    moon run build
fi

env $(grep -v '^#' .env | xargs) ./build/server
