# Bazaar Helm chart

## Install the chart

```
helm install --namespace bazaar --atomic  bazaar .
```

## Upgrading the chart

```
helm upgrade --namespace bazaar --atomic bazaar .
```

## Using locally with latest dev image

A `values.local.yml` file is provided which will point the installation to the latest `dev` image.
