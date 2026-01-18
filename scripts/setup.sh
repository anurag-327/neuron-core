#!/bin/bash

# 1. Build Docker Images for all supported runtimes
# These images are required for the container pool to work.
echo "Building Docker Images..."
docker build -t neuron-cpp runtime/images/cpp
docker build -t neuron-python runtime/images/python
docker build -t neuron-node runtime/images/node
docker build -t neuron-java runtime/images/java


