#!/bin/sh

minikube start --driver=docker

# Set the Docker environment to use Minikube’s Docker daemon
# This allows you to build and manage Docker images directly in Minikube’s context
eval $(minikube docker-env)

# Load local images from docker (make sure Docker Desktop is running)
minikube image load microservice-example-server:latest
minikube image load postgres:latest

# Setup namespace
# Right now default will be used for simplicity
# kubectl create namespace microservice-example

# Load the architecture
kubectl apply -f volumes/db-data-persistentvolumeclaim.yaml

kubectl apply -f secrets/db-password-secret.yaml

kubectl apply -f deployments/postgres-db-deployment.yaml
kubectl apply -f deployments/server-deployment.yaml

kubectl apply -f services/postgres-db-service.yaml
kubectl apply -f sservices/server-service.yaml

# Set the Docker environment to user Docker daemon again if needed
# eval $(minikube docker-env --unset)

# Expose the services to host machine
minikube tunnel 