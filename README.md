# credstash-operator 

## Overview

credstash-operator is a Kubernetes operator that creates Kubernetes secrets from credstash secrets

## Deployment
### Prerequisites

The controller requires AWS credentials to be set before deploying it. This is accomplished by creating a secret with name `aws-credentials` in the controller namespace with the following keys:
* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY
* AWS_REGION

For example running the following will create an appropriate secret in the `credstash` namespace:
```
kubectl create secret generic aws-credentials --from-literal=AWS_ACCESS_KEY_ID=access_key --from-literal=AWS_SECRET_ACCESS_KEY=secret_access_key --from-literal=AWS_REGION=us-west-2 --namespace=credstash
```

### Deploy the operator

Deploy the operator dependencies:
```
kubectl apply -f deploy/crds/credstash.ouzi.tech_credstashsecrets_crd.yaml
kubectl apply -f deploy/service_account.yaml -n credstash
kubectl apply -f deploy/role.yaml -n credstash
kubectl apply -f deploy/role_binding.yaml -n credstash
```

Deploy the operator:
```
kubectl apply -f deploy/operator.yaml -n credstash
```

### Deploying via helm chart

#### Without existing credentials secret
```
helm upgrade --install credstash deploy/helm/credstash-operator \
    -n credstash \
    --set awsCredentials.create=true \
    --set awsCredentials.awsAccessKeyId=access_key \
    --set awsCredentials.awsSecretAccessKey=secret_access_key \
    --set awsCredentials.awsRegion=region
```
#### With existing credentials secret
```
helm upgrade --install credstash deploy/helm/credstash-operator \
    -n credstash \
    --set awsCredentials.secretName=aws-credentials
```
