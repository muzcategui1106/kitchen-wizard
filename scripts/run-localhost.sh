#!/bin/bash -x

SCRIPT_PATH=`readlink -f "$0"`
SCRIPT_DIR=`dirname "$SCRIPT_PATH"`
export POSTGRES_DB_PASSWORD=`kubectl get secret kitchenwizard.acid-minimal-cluster.credentials.postgresql.acid.zalan.do -o 'jsonpath={.data.password}' | base64 -d`
S3_STORE_ENDPOINT="http://store.s3.local.uzcatm-skylab.com"

function main() {
    trap 'kill %1; kill %2' SIGINT # catch SIGINT and use it to terminate both functions
    run_backend &
    #run_frontend &

    wait
}

function run_backend() {
    $SCRIPT_DIR/../bin/api --dex-provider-url "https://dex.dex.local.uzcatm-skylab.com" \
    --oidc-client-id example-app \
    --oidc-client-secret ZXhhbXBsZS1hcHAtc2VjcmV0 \
    --oidc-redirect-url "http://localhost:8443" \
    --postgres-db-hostname localhost \
    --postgres-db-username kitchenwizard \
    --postgres-db-port "6432" \
    --store-s3-endpoint $S3_STORE_ENDPOINT
}

function run_frontend() {
    npm start --prefix $SCRIPT_DIR/../frontend
}

main