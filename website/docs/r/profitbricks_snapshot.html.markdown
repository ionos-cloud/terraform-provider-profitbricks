---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_snapshot"
sidebar_current: "docs-profitbricks-resource-snapshot"
description: |-
  Creates and manages snapshot objects.
---

# profitbricks\_snapshot

Manages snapshots on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_snapshot" "test_snapshot" {
  datacenter_id = "datacenterId"
  volume_id = "volumeId"
  name = "my snapshot"
}
```

##Argument reference

* `datacenter_id` - (Required) [Boolean] The ID of the Virtual Data Center.
* `name` - (Required) [Boolean] The name of the snapshot.
* `volume_id` - (Required) [Boolean] The ID of the specific volume to take snapshot from.
