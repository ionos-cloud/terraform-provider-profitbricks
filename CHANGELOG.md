## 1.3.4 (Unreleased)
## 1.3.3 (August 09, 2018)

IMPROVEMENTS:

* Handle empty endpoint ([#35](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/35))
* Update to profitbricks-sdk-go  v5.0.1 ([#34](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/34))
* Error handling and rename variables. ([#36](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/36))
* IPBlockUpdate added.  ([#37](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/37))

## 1.3.2 (July 03, 2018)

BUG FIXES:

* Reattaching a volume after tainting a server ([#33](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/33))

## 1.3.1 (May 29, 2018)

IMPROVEMENT

* Icmp type ([#32](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/32))

## 1.3.0 (May 29, 2018)

IMPROVEMENT

* Doc sync ([#31](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/31))

BUG FIXES:

* ICMP type and code had wrong type ([#30](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/30))

## 1.2.1 (May 23, 2018)

IMPROVEMENT

* Adding navigation for data source documentation. ([#29](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/29))

BUG FIXES:

* Retries attribute is marked as Deprecated instead of Remove. ([#28](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/28))

## 1.2.0 (May 15, 2018)

IMPROVEMENTS: 

* Documentation update ([#26](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/26))
* ProfitBricks provider support for terraform timeouts ([#25](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/25))

BUG FIXES:

* Fixes issue with server update affecting LAN assignment ([#24](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/24))
* Inconsistent interpretation of endpoint ([#25](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/25))
* Changing data center location silently ignored ([#25](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/25))

## 1.1.1 (February 27, 2018)

* Removed reboot from nic resource ([#22](https://github.com/terraform-providers/terraform-provider-profitbricks/issues/22))

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
