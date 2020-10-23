---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_private_crossconnect"
sidebar_current: "docs-profitbricks-resource-private-crossconnect"
description: |-
  Creates and manages Private Cross Connections between virtual datacenters.
---

# profitbricks_private_crossconnect

Manages a Private Cross Connect on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_private_crossconnect" "example" {
  name = "example"
  description = "example pcc"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the private cross-connection.
- `description` - (Optional)[string] A short description for the private cross-connection.

## Import

A Private Cross Connect resource can be imported using its `resource id`, e.g.

```shell
terraform import profitbricks_private_crossconnect.demo {profitbricks_private_crossconnect_uuid}
```

This can be helpful when you want to import private cross-connects which you have already created manually or using other means, outside of terraform.
