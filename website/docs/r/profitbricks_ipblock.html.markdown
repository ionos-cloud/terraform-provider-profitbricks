---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_ipblock"
sidebar_current: "docs-profitbricks-resource-ipblock"
description: |-
  Creates and manages IP Block objects.
---

# profitbricks\_ipblock

Manages IP Blocks on ProfitBricks. IP Blocks contain reserved public IP addresses that can be assigned servers or other resources.

## Example Usage

```hcl
resource "profitbricks_ipblock" "example" {
  location = "${profitbricks_datacenter.example.location}"
  size     = 1
}
```

##Argument reference

* `location` - (Required)[string] The regional location for this IP Block: us/las, us/ewr, de/fra, de/fkb.
* `size` - (Required)[integer] The number of IP addresses to reserve for this block.
* `ips` - (Computed)[integer] The list of IP addresses associated with this block.
