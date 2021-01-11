package profitbricks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceK8sNodePool_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_k8s_node_pool.test_ds_k8s_node_pool", "name", "test_nodepool"),
					resource.TestCheckResourceAttr("data.profitbricks_k8s_node_pool.test_ds_k8s_node_pool", "k8s_version", "1.18.12"),
				),
			},
		},
	})
}

func TestAccDataSourceK8sNodePool_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksK8sNodePoolMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_k8s_node_pool.test_ds_k8s_node_pool", "name", "test_nodepool"),
					resource.TestCheckResourceAttr("data.profitbricks_k8s_node_pool.test_ds_k8s_node_pool", "k8s_version", "1.18.12"),
					resource.TestCheckResourceAttrSet("data.profitbricks_k8s_node_pool.test_ds_k8s_node_pool", "id"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksK8sNodePoolCreateResources = `
resource "profitbricks_datacenter" "test_ds_k8s_datacenter" {
	name              = "test_datacenter"
	location          = "de/fra"
	description       = "test datacenter"
}

resource "profitbricks_k8s_cluster" "test_ds_k8s_cluster" {
	name              = "test_cluster"
}

resource "profitbricks_k8s_node_pool" "test_ds_k8s_node_pool" {
	depends_on 				= [profitbricks_datacenter.test_ds_k8s_datacenter, profitbricks_k8s_cluster.test_ds_k8s_cluster]
	name					= "test_nodepool"
	datacenter_id			= profitbricks_datacenter.test_ds_k8s_datacenter.id
	k8s_cluster_id			= profitbricks_k8s_cluster.test_ds_k8s_cluster.id
	node_count				= 1
	cpu_family				= "INTEL_XEON"
	cores_count				= 1
	ram_size				= 2048
	availability_zone 		= "AUTO"
	storage_type			= "HDD"
	storage_size			= 10
	k8s_version				= "1.18.12"
}
`

const testAccDataSourceProfitBricksK8sNodePoolMatchId = `
resource "profitbricks_datacenter" "test_ds_k8s_datacenter" {
  name              = "test_datacenter"
  location          = "de/fra"
  description       = "test datacenter"
}

resource "profitbricks_k8s_cluster" "test_ds_k8s_cluster" {
  name              = "test_cluster"
}

resource "profitbricks_k8s_node_pool" "test_ds_k8s_node_pool" {
	depends_on 				= [profitbricks_datacenter.test_ds_k8s_datacenter, profitbricks_k8s_cluster.test_ds_k8s_cluster]
  name							= "test_nodepool"
	datacenter_id			= profitbricks_datacenter.test_ds_k8s_datacenter.id
	k8s_cluster_id		= profitbricks_k8s_cluster.test_ds_k8s_cluster.id
	node_count				= 1
	cpu_family				= "INTEL_XEON"
	cores_count				= 1
	ram_size					= 2048
	availability_zone = "AUTO"
	storage_type			= "HDD"
	storage_size			= 10
	k8s_version				= "1.18.12"
  #	public_ips				= [ "85.184.250.211", "157.97.107.226", "157.97.107.242" ]
  # public_ips				= [ ]
  #   public_ips        = [ ]
}

data "profitbricks_k8s_node_pool" "test_ds_k8s_node_pool" {
	k8s_cluster_id 	= profitbricks_k8s_cluster.test_ds_k8s_cluster.id
	id				= profitbricks_k8s_node_pool.test_ds_k8s_node_pool.id
}
`

const testAccDataSourceProfitBricksK8sNodePoolMatchName = `
resource "profitbricks_datacenter" "test_ds_k8s_datacenter" {
  name              = "test_datacenter"
  location          = "de/fra"
  description       = "test datacenter"
}

resource "profitbricks_k8s_cluster" "test_ds_k8s_cluster" {
  name              = "test_cluster"
}

resource "profitbricks_k8s_node_pool" "test_ds_k8s_node_pool" {
	depends_on 				= [profitbricks_datacenter.test_ds_k8s_datacenter, profitbricks_k8s_cluster.test_ds_k8s_cluster]
  name							= "test_nodepool"
	datacenter_id			= profitbricks_datacenter.test_ds_k8s_datacenter.id
	k8s_cluster_id		= profitbricks_k8s_cluster.test_ds_k8s_cluster.id
	node_count				= 1
	cpu_family				= "INTEL_XEON"
	cores_count				= 1
	ram_size					= 2048
	availability_zone = "AUTO"
	storage_type			= "HDD"
	storage_size			= 10
	k8s_version				= "1.18.12"
  #	public_ips				= [ "85.184.250.211", "157.97.107.226", "157.97.107.242" ]
  # public_ips				= [ ]
  #   public_ips        = [ ]
}

data "profitbricks_k8s_node_pool" "test_ds_k8s_node_pool" {
	k8s_cluster_id 	= profitbricks_k8s_cluster.test_ds_k8s_cluster.id
	name			= "_nodepool"
}
`
