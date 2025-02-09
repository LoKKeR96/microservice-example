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
kubectl apply -f kubernetes/volumes/db-data-persistentvolumeclaim.yaml

kubectl apply -f kubernetes/secrets/db-password-secret.yaml

kubectl apply -f kubernetes/deployments/postgres-db-deployment.yaml
kubectl apply -f kubernetes/deployments/server-deployment.yaml

kubectl apply -f kubernetes/services/postgres-db-service.yaml
kubectl apply -f kubernetes/services/server-service.yaml

# Expose the services to host machine
minikube tunnel 