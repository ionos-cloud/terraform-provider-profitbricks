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
  datacenter_id = "${profitbricks_datacenter.example.id}"
  public        = true
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `name` - (Optional)[string] The name of the LAN.
* `public` - (Optional)[Boolean] Indicates if the LAN faces the public Internet (true) or not (false).
