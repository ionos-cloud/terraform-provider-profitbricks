---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_k8s_cluster"
sidebar_current: "docs-profitbricks-resource-k8s-cluster"
description: |-
  Creates and manages Profitbricks Kubernetes Clusters.
---

# profitbricks_k8s_cluster

Manages a Managed Kubernetes cluster on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_k8s_cluster" "example" {
  name        = "example"
  k8s_version = "1.18.5"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:30:00Z"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the Kubernetes Cluster.
- `k8s_version` - (Optional)[string] The desired Kubernetes Version. for supported values, please check the API documentation.
- `maintenance_window` - (Optional) See the **maintenance_window** section in the example above

## Import

A Kubernetes Cluster resource can be imported using its `resource id`, e.g.

```shell
terraform import profitbricks_k8s_cluster.demo {k8s_cluster uuid}
```

This can be helpful when you want to import kubernetes clusters which you have already created manually or using other means, outside of terraform.
