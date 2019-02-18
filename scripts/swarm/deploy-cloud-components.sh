#!/bin/bash

cp ./tmp/generated-gateway_config.yml ./tmp/openfaas-cloud/gateway_config.yml
cp ./tmp/generated-github.yml ./tmp/openfaas-cloud/github.yml
cp ./tmp/generated-slack.yml ./tmp/openfaas-cloud/slack.yml
cp ./tmp/generated-dashboard_config.yml ./tmp/openfaas-cloud/dashboard/dashboard_config.yml

./tmp/openfaas-cloud/of-builder/deploy_swarm.sh

AUTH_URL=echo

if [ "$ENABLE_OAUTH" = "true" ] ; then
    export TAG=0.4.0
    docker service rm auth
    docker service create --name auth \
    -e oauth_client_secret_path="/run/secrets/of-client-secret" \
    -e client_id="$CLIENT_ID" \
    -e PORT=8080 \
    -p 8085:8080 \
    -e external_redirect_domain="http://auth.system.$DOMAIN/" \
    -e cookie_root_domain=".system.$DOMAIN" \
    -e public_key_path=/run/secrets/jwt-public-key \
    -e private_key_path=/run/secrets/jwt-private-key \
    -e oauth_provider="github" \
    --network func_functions \
    --secret jwt-private-key \
    --secret jwt-public-key \
    --secret of-client-secret \
    openfaas/cloud-auth:$TAG

    AUTH_URL=auth
fi

echo "Deploying router service."

TAG=0.6.0
docker service rm of-router

docker service create --network=func_functions \
--env upstream_url=http://gateway:8080 \
--env auth_url=http://$AUTH_URL:8080 \
--publish 80:8080 \
--name of-router \
-d openfaas/cloud-router:$TAG

cd ./tmp/openfaas-cloud
faas-cli template pull

for i in {1..60};
do
    echo "Checking if OpenFaaS GW is up."

    curl -if 127.0.0.1:8080
    if [ $? == 0 ];
    then
        break
    fi

    sleep 1
done

ADMIN_PASSWORD=$(docker container exec $(docker ps --filter name=gateway -q) cat /var/run/secrets/basic-auth-password)

echo -n $ADMIN_PASSWORD | faas-cli login --username admin --password-stdin

faas-cli deploy

cd ./dashboard
faas-cli template pull https://github.com/openfaas-incubator/node8-express-template
faas-cli deploy
