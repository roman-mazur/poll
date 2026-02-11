#!/usr/bin/env bash

title=$1
key=$2
option=$3

svc_url="https://poll.rmazur.io"
secret=$(cat .creds/admin)
if [ "$option" == "dev" ]; then
  svc_url="http://localhost:17000"
  secret=""
fi

echo "===="
echo "New config for $svc_url"
echo "===="

talk_id=$(curl -f -X POST -H "Authorization: $secret" "$svc_url/config/new?key=$key")
echo "Poll response: $talk_id"

link="$svc_url?id=$(echo "$talk_id" | jq -sRr @uri)&name=$(echo "$title" | jq -sRr @uri)"

echo
echo
echo "Poll link: $link"
echo
echo

qrencode -o present/poll/qr.png "$link" || echo "ERROR: check if qrencode is installed"

go install golang.org/x/tools/cmd/present
cd present && present --base=.
