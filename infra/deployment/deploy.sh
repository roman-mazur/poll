#!/usr/bin/env bash

set -e

cd out

terraform apply deploy-plan
