#!/bin/bash

# Working around go-swagger bug: https://github.com/go-swagger/go-swagger/issues/230
# # Bug puts some of the generated code in the current directory under client
echo -e "${OK_COLOR}==> Calling go generate (jobber-client)${NO_COLOR}"
mkdir -p generated/jobber-client
cd generated/jobber-client; swagger generate client -A jobber -f ../../api/swagger.yml;

