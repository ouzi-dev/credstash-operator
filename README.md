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
kubectl apply -f deploy/deployment.yaml -n credstash
```

### Deploying via helm chart

#### Without existing credentials secret
```
helm upgrade --install credstash https://github.com/ouzi-dev/credstash-operator/releases/download/${VERSION}/credstash-operator-${VERSION}.tgz \
    -n credstash \
    --set awsCredentials.create=true \
    --set awsCredentials.awsAccessKeyId=access_key \
    --set awsCredentials.awsSecretAccessKey=secret_access_key \
    --set awsCredentials.awsRegion=region
```
Where ${VERSION} is the version you want to install
#### With existing credentials secret
```
helm upgrade --install credstash https://github.com/ouzi-dev/credstash-operator/releases/download/${VERSION}/credstash-operator-${VERSION}.tgz \
    -n credstash \
    --set awsCredentials.secretName=aws-credentials
``` 
Where ${VERSION} is the version you want to install

### Using the operator
Once configured submitting a CredstashSecret custom resource like below will create a secret from the credstash keys specified
```yaml
apiVersion: credstash.ouzi.tech/v1alpha1
kind: CredstashSecret
metadata:
  name: very-secret
  namespace: test
spec:
  # Name of the target secret (Optional. Defaults to the CR name)
  name: example
  # List of secrets from credstash to add to the body of the secret 
  secrets:
      # key: the key in credstash to fetch. (Required)
    - key: test-secret
      # name: the name of the resulting data element in the k8s secret (Optional. Defaults to the credstash key)
      name: renamed-test-secret
      # table: the dynamoDB table that contains the credstash secrets (Optional. Defaults to credential-store)
      table: credential-store
      # version: the version of the secret in credstash for the provided key (Optional.Defaults to the latest version)
      version: 1
```
