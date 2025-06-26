# Configuring complex inputs for Cloud automation for creating Kubernetes Resources

* namespaces (`namespaces`)

## namespaces <a name="namespaces"></a>

The `namespaces` input variable allows you to provide a list of namespaces along with labels and annotations which will be created in the cluster.

* Variable name: `namespaces`.
* Type: An object representing a set of namespaces.
* Default value: []

### Example namespaces

```hcl
[
  {
    name = "dev"
    metadata = {
      labels = {
        environment = "development"
        team        = "backend"
      }
      annotations = {
        "owner" = "dev-team@example.com"
      }
    }
  },
  {
    name = "prod"
    metadata = {
      labels = {
        environment = "production"
        team        = "frontend"
      }
      annotations = {
        "owner" = "prod-team@example.com"
      }
    }
  }
]
```
