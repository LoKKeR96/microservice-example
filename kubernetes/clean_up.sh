#!/bin/sh
kubectl delete service server
kubectl delete service postgres-db

kubectl delete deployment postgres-db
kubectl delete deployment server

# Delete Permanent Volumes
kubectl delete pvc --all