package profitbricks

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	profitbricks "github.com/profitbricks/profitbricks-sdk-go/v5"
)

func TestAccProfitBricksk8sNodepool_Basic(t *testing.T) {
	var k8sNodepool profitbricks.KubernetesNodePool
	k8sNodepoolName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDProfitBricksk8sNodepoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckProfitBricksk8sNodepoolConfigBasic, k8sNodepoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksk8sNodepoolExists("profitbricks_k8s_node_pool.example", &k8sNodepool),
					resource.TestCheckResourceAttr("profitbricks_k8s_node_pool.example", "name", k8sNodepoolName),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckProfitBricksk8sNodepoolConfigUpdate, k8sNodepoolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfitBricksk8sNodepoolExists("profitbricks_k8s_node_pool.example", &k8sNodepool),
					resource.TestCheckResourceAttr("profitbricks_k8s_node_pool.example", "name", k8sNodepoolName),
					resource.TestCheckResourceAttr("profitbricks_k8s_node_pool.example", "maintenance_window.0.day_of_the_week", "Tuesday"),
					resource.TestCheckResourceAttr("profitbricks_k8s_node_pool.example", "maintenance_window.0.time", "11:00:00Z"),
				),
			},
		},
	})
}

func testAccCheckDProfitBricksk8sNodepoolDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*profitbricks.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profitbricks_k8s_node_pool" {
			continue
		}

		_, err := client.GetKubernetesNodePool(rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID)

		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("K8s node pool still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch k8s node pool %s %s", rs.Primary.ID, err)
		}

		_, ddcErr := client.GetDatacenter(rs.Primary.Attributes["datacenter_id"])

		if apiError, ok := ddcErr.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("Data center for node pool still exists %s %s", rs.Primary.Attributes["k8s_cluster_id"], apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch data center for node pool %s %s", rs.Primary.Attributes["k8s_cluster_id"], err)
		}

		_, dkErr := client.GetKubernetesCluster(rs.Primary.Attributes["k8s_cluster_id"])

		if apiError, ok := dkErr.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() != 404 {
				return fmt.Errorf("K8s cluster for node pool still exists %s %s", rs.Primary.Attributes["k8s_cluster_id"], apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch k8s cluster for node pool %s %s", rs.Primary.Attributes["k8s_cluster_id"], err)
		}

	}

	return nil
}

func testAccCheckProfitBricksk8sNodepoolExists(n string, k8sNodepool *profitbricks.KubernetesNodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*profitbricks.Client)
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		foundK8sNodepool, err := client.GetKubernetesNodePool(rs.Primary.Attributes["k8s_cluster_id"], rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error occured while fetching k8s node pool: %s", rs.Primary.ID)
		}
		if foundK8sNodepool.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		k8sNodepool = foundK8sNodepool

		return nil
	}
}

const testAccCheckProfitBricksk8sNodepoolConfigBasic = `
resource "profitbricks_datacenter" "example" {
  name        = "example"
  location    = "de/fra"
  description = "Datacenter created through terraform"
}

resource "profitbricks_k8s_cluster" "example" {
  name        = "example"
  k8s_version = "1.17.5"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}

resource "profitbricks_k8s_node_pool" "example" {
  name        = "%s"
  k8s_version = "${profitbricks_k8s_cluster.example.k8s_version}"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
  datacenter_id     = "${profitbricks_datacenter.example.id}"
  k8s_cluster_id    = "${profitbricks_k8s_cluster.example.id}"
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
}`

const testAccCheckProfitBricksk8sNodepoolConfigUpdate = `
resource "profitbricks_datacenter" "example" {
  name        = "example"
  location    = "de/fra"
  description = "Datacenter created through terraform"
}

resource "profitbricks_k8s_cluster" "example" {
  name        = "example"
  k8s_version = "1.17.5"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00Z"
  }
}

resource "profitbricks_k8s_node_pool" "example" {
  name        = "%s"
  k8s_version = "${profitbricks_k8s_cluster.example.k8s_version}"
  maintenance_window {
    day_of_the_week = "Tuesday"
    time            = "11:00:00Z"
  }
  datacenter_id     = "${profitbricks_datacenter.example.id}"
  k8s_cluster_id    = "${profitbricks_k8s_cluster.example.id}"
  cpu_family        = "INTEL_XEON"
  availability_zone = "AUTO"
  storage_type      = "SSD"
  node_count        = 1
  cores_count       = 2
  ram_size          = 2048
  storage_size      = 40
}`
