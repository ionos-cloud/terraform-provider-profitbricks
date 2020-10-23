---
layout: "profitbricks"
page_title: "ProfitBricks : profitbrick_snapshot"
sidebar_current: "docs-profitbricks-datasource-snapshot"
description: |-
  Get information on a ProfitBricks Snapshots
---

# profitbricks\_snapshot

The snapshots data source can be used to search for and return an existing snapshot which can then be used to provision a server.

## Example Usage

```hcl
data "profitbricks_snapshot" "snapshot_example" {
  name     = "my snapshot"
  size     = "2"
  location = "location_id"
}
```

## Argument Reference

 * `name` - (Required) Name or part of the name of an existing snapshot that you want to search for.
 * `location` - (Optional) Id of the existing snapshot's location.
 * `size` - (Optional) The size of the snapshot to look for.

## Attributes Reference

 * `id` - UUID of the snapshot
