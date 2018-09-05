---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_nic"
sidebar_current: "docs-profitbricks-resource-nic"
description: |-
  Creates and manages Network Interface objects.
---

# profitbricks\_nic

Manages a NIC on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_nic" "example" {
  datacenter_id = "${profitbricks_datacenter.example.id}"
  server_id     = "${profitbricks_server.example.id}"
  lan           = 2
  dhcp          = true
  ip            = "${profitbricks_ipblock.example.ip}"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `server_id` - (Required)[string] The ID of a server.
* `lan` - (Required)[integer] The LAN ID the NIC will sit on.
* `name` - (Optional)[string] The name of the LAN.
* `dhcp` - (Optional)[Boolean] Indicates if the NIC should get an IP address using DHCP (true) or not (false).
* `ip` - (Optional)[string] IP assigned to the NIC.
* `firewall_active` - (Optional)[Boolean] If this resource is set to true and is nested under a server resource firewall, with open SSH port, resource must be nested under the NIC.
* `nat` - (Optional)[Boolean] Boolean value indicating if the private IP address has outbound access to the public internet.
* `ips` - (Computed) The IP address or addresses assigned to the NIC.

## Import

Resource Nic can be imported using the `resource id`, e.g.

```shell
terraform import profitbricks_nic.mynic {datacenter uuid}/{server uuid}/{nic uuid}
```
