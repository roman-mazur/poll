#!/usr/bin/env bash

set -e

mkdir -p ./out 2>/dev/null

cue export -e terraform > out/infra.tf.json

cd out

terraform init
terraform plan -out deploy-plan
