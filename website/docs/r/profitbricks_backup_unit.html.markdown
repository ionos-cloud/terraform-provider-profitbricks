---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_backup_unit"
sidebar_current: "docs-profitbricks-resource-backup-unit"
description: |-
  Creates and manages Profitbricks Backup Units.
---

# profitbricks_backup_unit

Manages a Backup Unit on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_backup_unit" "example" {
  name        = "example"
  password    = "<example-password>"
  email       = "example@example-domain.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the Backup Unit.
- `password` - (Required)[string] The desired password for the Backup Unit.
- `email` - (Required)[string] The email address assigned to the backup unit

## Import

A Backup Unit resource can be imported using its `resource id`, e.g.

```shell
terraform import profitbricks_backup_unit.demo {backup_unit_uuid}
```

This can be helpful when you want to import backup units which you have already created manually or using other means, outside of terraform. Please note that you need to manually specify the password when first declaring the resource in terraform, as there is no way to retrieve the password from the Cloud API.

## Important Notes

- Please note that at the moment, Backup Units cannot be renamed
- Please note that the password attribute is write-only, and it cannot be retrieved from the API when importing a profitbricks_backup_unit. Basically, the only way to keep track of it in terraform is to specify it.