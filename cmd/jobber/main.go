package main

import (
	"fmt"
	"log"
	"os"

	spec "github.com/go-swagger/go-swagger/spec"
	flags "github.com/jessevdk/go-flags"
	"github.com/ritchida/jobber/cmd/jobber/restapi"
	swaggerapi "github.com/ritchida/jobber/generated/jobber/restapi"
	"github.com/ritchida/jobber/generated/jobber/restapi/operations"
	"github.com/ritchida/jobber/pkg/config"
	"github.com/ritchida/jobber/pkg/repository"
)

var options struct {
	Host string `long:"host" description:"the IP to listen on" default:"127.0.0.1" env:"HOST"`
	Port int    `long:"port" description:"the port to listen on for insecure connections, defaults to a random value" env:"PORT"`
}

func main() {
	swaggerSpec, err := spec.New(swaggerapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	parser := flags.NewParser(&options, flags.Default)
	parser.ShortDescription = swaggerSpec.Spec().Info.Title
	parser.LongDescription = swaggerSpec.Spec().Info.Description

	api := operations.NewJobberAPI(swaggerSpec)

	config, err := initializeServiceConfig()
	if err != nil {
		fmt.Printf("Unable to start service: %v\n", err)
		os.Exit(1)
	}

	server := restapi.NewServer(api, *config)
	defer server.Shutdown()

	for _, optsGroup := range api.CommandLineOptionsGroups {
		parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
	}

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func initializeServiceConfig() (*config.JobberConfig, error) {
	// load the configuration
	config, configErrors := config.GetJobberConfig()
	if len(configErrors) > 0 {
		for _, configErr := range configErrors {
			fmt.Printf("Configuration error: %v\n", configErr)
		}
		return nil, fmt.Errorf("Unable to initialize service from configuration")
	}

	// initialize the repository session
	_, err := repository.GetCassandraJobberRepository()
	if err != nil {
		return nil, err
	}

	return &config, nil
}
