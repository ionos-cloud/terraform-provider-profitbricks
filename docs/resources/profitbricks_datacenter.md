---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_datacenter"
sidebar_current: "docs-profitbricks-resource-datacenter"
description: |-
  Creates and manages Profitbricks Virtual Data Center.
---

# profitbricks\_datacenter

Manages a Virtual Data Center on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_datacenter" "example" {
  name        = "datacenter name"
  location    = "us/las"
  description = "datacenter description"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)[string] The name of the Virtual Data Center.
* `location` - (Required)[string] The regional location where the Virtual Data Center will be created.
* `description` - (Optional)[string] Description for the Virtual Data Center.

## Import

Resource Datacenter can be imported using the `resource id`, e.g.

```shell
terraform import profitbricks_datacenter.mydc {datacenter uuid}
```
