---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_loadbalancer"
sidebar_current: "docs-profitbricks-resource-loadbalancer"
description: |-
  Creates and manages Load Balancers
---

# profitbricks\_loadbalancer

Manages a Load Balancer on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_loadbalancer" "example" {
  datacenter_id = "${profitbricks_datacenter.example.id}"
  nic_ids        = ["${profitbricks_nic.example.id}"]
  name          = "load balancer name"
  dhcp          = true
}
```

## Argument reference

* `name` - (Required)[string] The name of the load balancer.
* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `nic_ids` - (Required)[list] A list of NIC IDs that are part of the load balancer.
* `dhcp` - (Optional)[Boolean] Indicates if the load balancer will reserve an IP using DHCP.
* `ip` - (Optional)[string] IPv4 address of the load balancer.
