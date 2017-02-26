generated/jobber: #api/swagger.yml
		recipes/generate-server.sh
			touch $@

generated/jobber-client: api/swagger.yml
		recipes/generate-client.sh
			touch $@

generate: generated/jobber generated/jobber-client 

clean:
	rm -rf generated

test-repo:
	go test github.com/ritchida/jobber/pkg/repository

test-api:
	./recipes/start-web-svc.sh
	go test github.com/ritchida/jobber/pkg/repository
	./recipes/stop-web-svc.sh

test: test-repo test-api
