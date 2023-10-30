#!/usr/bin/env bash

rm -rf *_out.cue

cue cmd check && echo "OK"
