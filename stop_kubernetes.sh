#!/bin/sh

sh ./kubernetes/clean_up.sh

minikube stop

minikube delete --all
