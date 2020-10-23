---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_s3_key"
sidebar_current: "docs-profitbricks-resource-s3-key"
description: |-
  Creates and manages Profitbricks S3 keys.
---

# profitbricks_s3_key

Manages an S3 Key on ProfitBricks.

## Example Usage

```hcl
resource "profitbricks_s3_key" "demo" {
  user_id    = <user-uuid>
  active     = true
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required)[string] The UUID of the user owning the S3 Key.
- `active` - (Required)[boolean] Whether the S3 is active / enabled or not - Please keep in mind this is only required on create.

## Import

An S3 Unit resource can be imported using its user id as well as its `resource id`, e.g.

```shell
terraform import profitbricks_s3_key.demo {userId}/{s3KeyId}
```

This can be helpful when you want to import S3 Keys which you have already created manually or using other means, outside of terraform.
