#!/usr/bin/env bash
set -euo pipefail

xargs -n1 cue fmt -s <<< "${@}"
