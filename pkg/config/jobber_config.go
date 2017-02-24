package config

import (
	"os"
	"strings"
)

// JobberConfig hlods the config parameters for the Jobber service
type JobberConfig struct {
	Cassandra CassandraConfig
}

// CassandraConfig holds the coinfig values needed
// to interact with a Cassandra cluster
type CassandraConfig struct {
	ClusterNodeIPs string
}

// ReadConfig reads in the jobber service configuration
func ReadConfig() JobberConfig {
	config := JobberConfig{
		Cassandra: getCassandraEnvConfig(),
	}
	return config
}

func getCassandraEnvConfig() CassandraConfig {
	config := CassandraConfig{}
	envVars := os.Environ()
	for _, envVar := range envVars {
		parts := strings.Split(envVar, "=")
		switch parts[0] {
		case "CASSANDRA_CLUSTER_IPS":
			config.ClusterNodeIPs = parts[1]
		default:
		}
	}
	return config
}
