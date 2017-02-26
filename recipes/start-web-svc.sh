#!/bin/bash

set -e

echo "Starting web service..."
go run cmd/jobber/main.go&
sleep 5
echo "Web service started."
