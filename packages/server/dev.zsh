#!/usr/bin/bash

if ! [[ -s "$(ls ./server)" ]]; then
    moon run build
fi

env $(grep -v '^#' .env | xargs) ./server
