#!/bin/bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo -e "${OK_COLOR}==> Calling go generate${NO_COLOR}"
# dependency ordering problem. would be best if we could just call the makefile here.
go generate -x ${ALLPKGS}
