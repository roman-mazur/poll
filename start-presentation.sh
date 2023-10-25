#!/usr/bin/env bash

title=$1

encoded_title=$(echo "$title" | jq -sRr @uri)

talk_id=$(curl -f -X POST -H "Authorization: $(cat .creds/admin)" "https://poll.rmazur.io/config/new/$encoded_title" 2>/dev/null)

link="https://poll.rmazur.io?id=$(echo "$talk_id" | jq -sRr @uri)"

qrencode -o present/poll/qr.png "$link"

cd present && present --base=.
