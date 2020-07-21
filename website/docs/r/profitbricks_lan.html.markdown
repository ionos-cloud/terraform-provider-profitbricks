---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_lan"
sidebar_current: "docs-profitbricks-resource-lan"
description: |-
  Creates and manages LAN objects.
---

# profitbricks\_lan

Manages a LAN on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_lan" "example" {
  datacenter_id = profitbricks_datacenter.example.id
  public        = true
  pcc           = profitbricks_private_crossconnect.example.id
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `name` - (Optional)[string] The name of the LAN.
* `public` - (Optional)[Boolean] Indicates if the LAN faces the public Internet (true) or not (false).
* `pcc` - (Optional)[String] The unique id of a `profitbricks_private_crossconnect` resource, in order

## Import

Resource Lan can be imported using the `resource id`, e.g.

```shell
terraform import profitbricks_lan.mylan {datacenter uuid}/{lan id}
```

## Important Notes

- Please note that only LANS datacenters found in the same physical location can be connected through a private cross-connect
- A LAN cannot be a part of two private cross-connects