#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [[ $OSTYPE == 'darwin'* ]]; then
    sudo docker info
    if [ $? -ne 0 ]; then
        echo "please start your docker engine before proceeding"
        exit 1
    fi

else
    systemctl is-active docker --quiet
    if [ $? -ne 0 ]; then
        echo "starting docker daemon using systemd"
        systemctl start docker
    fi
fi

echo "creating kind cluster"
sudo kind create cluster --config $SCRIPT_DIR/kind-config.yaml

echo "applying contour ingress controllers"
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl patch daemonsets -n projectcontour envoy -p '{"spec":{"template":{"spec":{"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/control-plane","operator":"Equal","effect":"NoSchedule"},{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'

echo "installing cert-manager"
kubectl create namespace cert-manager
AWS_SECRET=`cat ~/.aws/credentials | grep secret | awk -F= '{print $NF}' | base64`
kubectl apply -n cert-manager -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: aws-creds
data:
  secret-key: $AWS_SECRET
EOF

sudo kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
echo "waiting 10 seconds for cert-manager to startup properly"
sleep 10 


echo "applying observability layer"
kubectl create namespace observability
kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.42.0/jaeger-operator.yaml -n observability
kubectl apply -n observability -f - <<EOF
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: collector
spec:
  ingress:
    annotations:
      kubernetes.io/ingress.class: contour
    hosts:
    - collector.observability.local.uzcatm-skylab.com
EOF

##### postgres operator and instance #########
# add repo for postgres-operator
helm repo add postgres-operator-charts https://opensource.zalando.com/postgres-operator/charts/postgres-operator

# install the postgres-operator
helm install postgres-operator postgres-operator-charts/postgres-operator
################################################