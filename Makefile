generated/jobber: #api/swagger.yml
		recipes/generate-server.sh
			touch $@

generated/jobber-client: api/swagger.yml
		recipes/generate-client.sh
			touch $@

generate: generated/jobber generated/jobber-client 

clean-generated:
	rm -rf generated

clean: clean-generated

test-repo:
	go test github.com/ritchida/jobber/pkg/repository

test-api:
	./recipes/start-web-svc.sh
	go test github.com/ritchida/jobber/cmd/jobber/restapi
	./recipes/stop-web-svc.sh

test: test-repo test-api
