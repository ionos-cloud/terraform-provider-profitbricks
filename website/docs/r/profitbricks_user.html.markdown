---
layout: "profitbricks"
page_title: "ProfitBricks: profitbricks_user"
sidebar_current: "docs-profitbricks-resource-user"
description: |-
  Creates and manages user objects.
---

# profitbricks\_user

Manages users and list users and groups associated.

## Example Usage

```hcl
resource "profitbricks_user" "user" {
  first_name = "terraform"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}
```

##Argument reference

* `administrator` - (Required) [Boolean] The group has permission to edit privileges on this resource.
* `email` - (Required) [string] An e-mail address for the user.
* `first_name` - (Required) [string] A name for the user.
* `force_sec_auth` - (Required) [Boolean] The group has permission to user this resource.
* `last_name` - (Required) [string] A name for the user.
* `password` - (Required) [string] A password for the user.