# Jobber

Jobber is a Cassandra-backed job reporting service.

## Running tests

On a clean machine you must first generate the API and related objects:
```bash
make generate
```
To run the tests you must have an instance of Cassandra running and publicly-accessible.
Once you have an instance/cluster running, set the CASSANDRA_CLUSTER_IPS environment variable:
```bash
export CASSANDRA_CLUSTER_IPS=<ip of cassandra node>
```
Then you can run the Cassandra integration tests:
```bash
go test github.com/ritchida/jobber/pkg/repository
```

