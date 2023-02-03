#!/bin/bash -x

systemctl is-active docker --quiet
if [ $? -ne 0 ]; then
    echo "starting docker daemon"
    systemctl start docker
fi

echo "creating kind cluster"
sudo kind create cluster

echo "applying contour ingress controllers"
sudo kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
sudo kubectl patch daemonsets -n projectcontour envoy -p '{"spec":{"template":{"spec":{"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/control-plane","operator":"Equal","effect":"NoSchedule"},{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'