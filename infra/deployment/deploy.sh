#!/usr/bin/env bash

set -e

cd out

if [ "$1" == "destroy" ]; then
  terraform destroy --auto-approve && exit 0
fi

terraform apply deploy-plan
