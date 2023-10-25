#!/usr/bin/env bash

title=$1

encoded_title=$(echo "$title" | jq -sRr @uri)

svc_url="https://poll.rmazur.io"
secret=$(cat .creds/admin)
if [ "$2" == "dev" ]; then
  svc_url="http://localhost:17000"
  secret=""
fi

talk_id=$(curl -f -X POST -H "Authorization: $secret" "$svc_url/config/new?name=$encoded_title" 2>/dev/null)

link="$svc_url?id=$(echo "$talk_id" | jq -sRr @uri)"

echo "Pool link: $link"

qrencode -o present/poll/qr.png "$link"

cd present && present --base=.
