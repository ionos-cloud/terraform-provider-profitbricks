---
layout: "profitbricks"
page_title: "ProfitBricks: ipfailover"
sidebar_current: "docs-profitbricks-resource-ipfailover"
description: |-
  Creates and manages ipfailover objects.
---

# profitbricks\_ipfailover

Manages IP Failover groups on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_ipfailover" "failovertest" {
  datacenter_id = "datacenterId"
  lan_id="lanId"
  ip ="reserved IP"
  nicuuid= "nicId"
}
```

##Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `ip` - (Required)[string] The reserved IP address to be used in the IP failover group.
* `lan_id` - (Required)[string] The ID of a LAN.
* `nicuuid` - (Required)[string] The ID of a NIC.
