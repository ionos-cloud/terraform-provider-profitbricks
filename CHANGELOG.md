## 1.1.0 (January 31, 2018)

* resource/profitbricks_loadbalancer: Removed `nic_id` parameter entirely ([#21](https://github.com/terraform-providers/terraform-provider-profitbricks/issues/21))

## 1.0.1 (January 05, 2018)

* Updated dependency profitbricks-sdk-go to 4.0.4. ([#18](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/18))

## 1.0.0 (October 04, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Deprecated nic_id parameter in profitbricks_loadbalancer and replaced it with nic_ids ([#15](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/15))

IMPROVEMENTS: 

* Fix IPFailover failing test ([#13](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/13))
* Fix issue with failover test failing ([#13](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/13))
* Deprecated nic_id parameter in profitbricks_loadbalancer and replaced it with nic_ids ([#15](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/15))


## 0.1.2 (August 23, 2017)

FEATURES:

* **New Data Source:** `profitbricks_resource` ([#8](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/8))
* **New Data Source:** `profitbricks_snapshot` ([#8](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/8))
* **New Resource:** `profitbricks_group` ([#11](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/11))
* **New Resource:** `profitbricks_share` ([#11](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/11))
* **New Resource:** `profitbricks_user` ([#11](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/11))
* **New Resource:** `profitbricks_snapshot` ([#11](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/11))
* **New Resource:** `profitbricks_ipfailover` ([#11](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/11))

IMPROVEMENTS: 

* Update `profitbricks_datacenter` with getImagealias method ([#8](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/8))
* Update `profitbricks_server` to use getImagealias method ([#8](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/8))
* Update `resource_profitbricks_lan` to read IP failover groups ([#11](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/11))
* Update `resource_profitbricks_volume` with imageAlias feature ([#11](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/11))

## 0.1.1 (July 31, 2017)
IMPROVEMENTS: 

* Acceptance test fix ([#1](https://github.com/terraform-providers/terraform-provider-profitbricks/issues/1))
* Added ability to pass snapshot id or name when creating a volume or a server ([#6](https://github.com/terraform-providers/terraform-provider-profitbricks/issues/6))

BUG FIXES:

* resource/profitbricks_server - Fix how primary_nic is updated ([#5](https://github.com/terraform-providers/terraform-provider-profitbricks/issues/5))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
