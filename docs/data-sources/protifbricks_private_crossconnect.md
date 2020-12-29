---
layout: "profitbricks"
page_title: "ProfitBricks : profitbrick_private_crossconnect"
sidebar_current: "docs-profitbricks-datasource-private-crossconnect"
description: |-
Get information on a ProfitBricks Private Crossconnects
---

# profitbricks\_private_crossconnect

The private crossconnect data source can be used to search for and return existing private crossconnects.

## Example Usage

```hcl
data "profitbricks_private_crossconnect" "pcc_example" {
  name     = "My PCC"
}
```

## Argument Reference

* `name` - (Optional) Name or part of the name of an existing private crossconnect that you want to search for.
* `id` - (Optional) ID of the private crossconnect you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id`
* `name`
* `description`
* `peers` - list of
    * `lan_id`
    * `lan_name`
    * `datacenter_id`
    * `datacenter_name`
    * `location`
* `connectable_datacenters` - list of
    * `id`
    * `name`
    * `location`
