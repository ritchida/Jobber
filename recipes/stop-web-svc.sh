#!/bin/bash

set -e

echo "Stopping web service..."
ps aux | grep "command-line-arguments/_obj/exe/main" | grep -v grep | awk '{print $2}' | xargs kill -9
sleep 5
echo "Web service stopped."
