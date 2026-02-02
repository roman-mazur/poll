#!/usr/bin/env bash

set -e

mkdir -p ./out 2>/dev/null

poll_cert=$(cat ../../.creds/cert.pem)
poll_key=$(cat ../../.creds/pkey.pem)
admin_secret=$(cat ../../.creds/admin)

cue export -t cert="$poll_cert" -t pkey="$poll_key" -t admin="$admin_secret" -e terraform > out/infra.tf.json

export CLOUDFLARE_API_TOKEN=$(cat ../../.creds/cloudflare)

cd out

tofu init --upgrade
tofu plan -out deploy-plan
