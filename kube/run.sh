#!/bin/bash

curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.17.2 sh -
cd istio-1.17.2
export PATH=$PWD/bin:$PATH

istioctl manifest apply --set profile=demo -y
istioctl install --set profile=demo --set meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY -y
cd ..

kubectl label namespace default istio-injection=enabled
kubectl apply -f pokemon-service/service.yml -f pokemon-service/configmap/configmap.yml
kubectl apply -f ingress/gateway.yml -f ingress/virtual-service.yml
kubectl apply -f egress/gateway.yml -f egress/virtual-service.yml -f egress/service-entry.yml
kubectl apply -f nginx/service.yml -f nginx/configmap/configmap.yml

minikube tunnel
