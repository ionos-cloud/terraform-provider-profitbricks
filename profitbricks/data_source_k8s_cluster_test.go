package profitbricks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceK8sCluster_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksK8sClusterCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksK8sClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_k8s_cluster.test_ds_k8s_cluster", "name", "TEST DS K8S CLUSTER"),
					resource.TestCheckResourceAttr("data.profitbricks_k8s_cluster.test_ds_k8s_cluster", "k8s_version", "1.18.12"),
					resource.TestCheckResourceAttrSet("data.profitbricks_k8s_cluster.test_ds_k8s_cluster", "kube_config"),
				),
			},
		},
	})
}

func TestAccDataSourceK8sCluster_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProfitBricksK8sClusterCreateResources,
			},
			{
				Config: testAccDataSourceProfitBricksK8sClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.profitbricks_k8s_cluster.test_ds_k8s_cluster", "name", "TEST DS K8S CLUSTER"),
					resource.TestCheckResourceAttr("data.profitbricks_k8s_cluster.test_ds_k8s_cluster", "k8s_version", "1.18.12"),
					resource.TestCheckResourceAttrSet("data.profitbricks_k8s_cluster.test_ds_k8s_cluster", "kube_config"),
					resource.TestCheckResourceAttrSet("data.profitbricks_k8s_cluster.test_ds_k8s_cluster", "id"),
				),
			},
		},
	})

}

const testAccDataSourceProfitBricksK8sClusterCreateResources = `
resource "profitbricks_k8s_cluster" "test_ds_k8s_cluster" {
  name         = "TEST DS K8S CLUSTER"
  k8s_version  = "1.18.12"
}
`

const testAccDataSourceProfitBricksK8sClusterMatchId = `
resource "profitbricks_k8s_cluster" "test_ds_k8s_cluster" {
  name         = "TEST DS K8S CLUSTER"
  k8s_version  = "1.18.12"
}

data "profitbricks_k8s_cluster" "test_ds_k8s_cluster" {
  id	= profitbricks_k8s_cluster.test_ds_k8s_cluster.id
}
`

const testAccDataSourceProfitBricksK8sClusterMatchName = `
resource "profitbricks_k8s_cluster" "test_ds_k8s_cluster" {
  name         = "TEST DS K8S CLUSTER"
  k8s_version  = "1.18.12"
}

data "profitbricks_k8s_cluster" "test_ds_k8s_cluster" {
  name	= "DS K8S"
}
`
