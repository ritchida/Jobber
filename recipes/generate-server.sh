#!/bin/bash

echo -e "${OK_COLOR}==> Calling go generate (jobber)${NO_COLOR}"
swagger generate server -A jobber -t generated/jobber -P models.User -f api/swagger.yml
#sed -e 's/restapi/main/' generated/jobber/restapi/doc.go > cmd/jobber/doc.go
sed -e 's/restapi/restapi/' generated/jobber/restapi/embedded_spec.go > cmd/jobber/restapi/embedded_spec.go

