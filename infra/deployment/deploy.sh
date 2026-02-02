#!/usr/bin/env bash

set -e

export CLOUDFLARE_API_TOKEN=$(cat ../../.creds/cloudflare)
cd out

if [ "$1" == "destroy" ]; then
  tofu destroy --auto-approve
else
  tofu apply deploy-plan
fi

# Persist the latest deploy state for monitoring queries.
tofu output -json | cue import -f -o ../state/deploy-outputs.cue -p state -l 'deployData:' json: -
