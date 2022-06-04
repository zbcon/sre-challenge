#!/bin/bash

## Deploy invoice-app & Service
kubectl apply -f invoice-app/deployment.yaml
kubectl apply -f invoice-app/service.yaml

## Deploy payment-provider app & Service
kubectl apply -f payment-provider/deployment.yaml
kubectl apply -f payment-provider/service.yaml