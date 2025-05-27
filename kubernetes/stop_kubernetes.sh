#!/bin/sh

sh ./clean_up.sh

minikube stop

minikube delete --all
