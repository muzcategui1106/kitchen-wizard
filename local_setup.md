# Overview

the following outlines basic setup that needs to be done as a one time thing for your application to work correctly when deployed in local mode

* create an AWS public hostzed zone
* create entries for the following inside the public hosted_zone you created before. all entries should be A records mapping back to 127.0.0.1
    * api.kitchen-wizard.local
    * collector.observability.local
    * dex.dex.local
* adjust URLs the [helm values](./deploy/k8s/values-local.yaml)  to match your domain name
* ensure you have valid AWS credentials inside ~/.aws/credentials that can do actions describe in https://cert-manager.io/docs/configuration/acme/dns01/route53/
* adjust helm value to ensure the AWS_ACCESS_KEY_ID value is correct
