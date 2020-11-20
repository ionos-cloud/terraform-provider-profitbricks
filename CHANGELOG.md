## 1.6.0 (Unreleased)

## 1.5.9 (November 20th, 2020)

BUG FIXES:
- The missing documentation regarding k8s node pools public IPs was added.

## 1.5.8 (November 20th, 2020)

FEATURES:
- Public IPs support in k8s node pools

## 1.5.7 (September 17, 2020)

BUG FIXES:
- Fixed issue #78, a bug causing the node count for a node pool from being increased

## 1.5.6 (September 17, 2020)

BUG FIXES:
- Fixed SegFault during apply

## 1.5.5 (August 28, 2020)

FEATURES:
- S3 Keys support (CRUD + Import) + documentation
- Backup Units support (CRUD + Import) + documentation

BUG FIXES:
- Fixed an error that was preventing the provider to find an image when its uuid was used instead of its name - Fixes #73 
- Fixed an error causing the **boot_image** attribute of a server   **profitbricks_server** to be changed if it was specified on create.
- the node_count property of a **profitbricks_k8s_node_pool** is now dynamically updated when autoscaling is on
- Ensures the **lans** attribute for a **profitbricks_k8s_node_pool** are now handled upon creation as well
- Some typos within the docs, code are now fixed

ENHANCEMENTS:
- Updated the provider to use golang sdk **github.com/profitbricks/profitbricks-sdk-go/v5** **v5.0.26**

## 1.5.4 (July 23, 2020)

FEATURES:

- Additional **lans** support for the **profitbricks_k8s_node_pool** resource (including import state sync)
- New **profitbricks_private_crossconnect** resource
- New **pcc** property for the **profitbricks_lan** resource so that it works in conjunction with private cross-connects
- Acceptance Tests for **profitbricks_private_crossconnect** (including an import test)
- **profitbricks_k8s_cluster** resource documentation updates

ENHANCEMENTS

- Revamped the **profitbricks_lan** resource implementation
- Documentation website updates for all new functionalities, as well as clarifying dependencies

## 1.5.3 (June 30, 2020)


BUG FIXES:
- Fix missing .Timeout property for resource definitions

ENHACEMENTS:
- Use golang sdk v5.0.16

FEATURES:
- Add autoscaling support

## 1.5.2 (May 27, 2020)

BUG FIXES
- Fixes a typo preventing users from updating a k8s node pool's nodes count
- Fixes #56 
- Fixes #66 
- Updates for all go dependencies

## 1.5.1 (May 25, 2020)

ENHANCEMENTS:
- Golang SDK version **v5.0.14**
- Added Kubernetes cluster resource
- Added Kubernetes node pool resource
- Possibility to import existing Kubernetes clusters, which were not created via Terraform
- Possibility to import existing Kubernetes node pools, which were not created via Terraform
- CRUD acceptance test for the Kubernetes cluster resource
- CRUD acceptance test for the Kubernetes node pool resource
- Acceptance test for the kubernetes cluster import functionality
- Acceptance test for the kubernetes node pool import functionality

BUG FIXES:
- A bugfix to the appearance of the provider documentation page

## 1.5.0 (April 09, 2020)

ENHANCEMENTS:
- Use profitbricks-sdk-go v5.0.9
- Use golang 1.13+
- Uses the new v5 API endpoint
- Terraform v0.12.21

BUG FIXES:
- Persist resource ids in state before syncing [#57](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/57)

BACKWARDS INCOMPATIBILITIES / NOTES:
- This provider version might have some incompatibility issues with older terraform code

## 1.4.5 (Unreleased)

BUG FIXES:
- Persist resource ids in state before syncing [#57](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/57)

## 1.4.4 (April 04, 2019)

BUG FIXES:
- Set cpu_family to computed. ([#54](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/54))

## 1.4.3 (March 14, 2019)

BUG FIXES:

- Fix nic.0.dhcp & nic.0.firewall_active update have no effect  ([#51](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/51))

ENHANCEMENTS:

- Use go modules

## 1.4.2 (February 11, 2019)

BUG FIXES:

* Fix conflicting auth token default value ([#47](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/47))

## 1.4.1 (January 21, 2019)

ENHANCEMENTS:

* Add parameter check for datacenter resource ([45](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/45))

## 1.4.0 (September 11, 2018)

ENHANCEMENTS:

* Add importer to server, nic, lan, dc, ipblock and firewall resources ([#38](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/38))

IMPROVEMENTS:

* Allowing usage of private images when provisioning a server ([#39](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/39))
* Fix for image property when using image alias ([#40](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/40))
* Force recreation of a resource when parents is changed ([#41](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/41))
* Discovering and reattaching detached volumes ([#42](https://github.com/terraform-providers/terraform-provider-profitbricks/pull/42)) 


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
