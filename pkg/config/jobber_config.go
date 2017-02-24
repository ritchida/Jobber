package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	jobberHost string = "127.0.0.1"
	jobberPort int32  = 9090
)

// JobberConfig hlods the config parameters for the Jobber service
type JobberConfig struct {
	Host      string
	Port      int32
	Cassandra CassandraConfig
}

// CassandraConfig holds the coinfig values needed
// to interact with a Cassandra cluster
type CassandraConfig struct {
	ClusterNodeIPs string
}

// ReadConfig reads in the jobber service configuration
func ReadConfig() (JobberConfig, []error) {
	return getJobberEnvConfig()
}

func getJobberEnvConfig() (JobberConfig, []error) {
	errors := []error{}
	config := JobberConfig{
		Host:      jobberHost,
		Port:      jobberPort,
		Cassandra: getCassandraEnvConfig(),
	}
	envVars := os.Environ()
	for _, envVar := range envVars {
		parts := strings.Split(envVar, "=")
		switch parts[0] {
		case "JOBBER_HOST":
			config.Host = parts[1]
		case "JOBBER_PORT":
			value, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				errors = append(errors, err)
			} else {
				config.Port = int32(value)
			}
		default:
		}
	}
	return config, errors
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
