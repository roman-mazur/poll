#!/usr/bin/env bash

set -e

export CLOUDFLARE_API_TOKEN=$(cat ../../.creds/cloudflare)
cd out

if [ "$1" == "destroy" ]; then
  terraform destroy --auto-approve
else
  terraform apply deploy-plan
fi

# Persist the latest deploy state for monitoring queries.
terraform output -json | cue import -f -o ../state/deploy-outputs.cue -p state -l 'deployData:' json: -
