generated/jobber: #api/swagger.yml
		recipes/generate-server.sh
			touch $@

generated/jobber-client: api/swagger.yml
		recipes/generate-client.sh
			touch $@

generate: generated/jobber generated/jobber-client 

clean:
	rm -rf generated
