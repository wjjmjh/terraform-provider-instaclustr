package test

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/instaclustr/terraform-provider-instaclustr/instaclustr"
	"io/ioutil"
	"os"
	"testing"
)

func TestAccCluster_importBasic(t *testing.T) {
	testAccProvider := instaclustr.Provider()
	testAccProviders := map[string]terraform.ResourceProvider{
		"instaclustr": testAccProvider,
	}
	validConfig, _ := ioutil.ReadFile("data/valid_single_bundle.tf")
	username := os.Getenv("IC_USERNAME")
	apiKey := os.Getenv("IC_API_KEY")
	hostname := getOptionalEnv("IC_API_URL", instaclustr.DefaultApiHostname)
	oriConfig := fmt.Sprintf(string(validConfig), username, apiKey, hostname)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { AccTestEnvVarsCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckResourceDeleted("valid", hostname, username, apiKey),
		Steps: []resource.TestStep{
			{
				Config: oriConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceValid("valid"),
					testCheckResourceCreated("valid", hostname, username, apiKey),
				),
			},
			{
				Config:            oriConfig,
				ResourceName:      "instaclustr_cluster.valid",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEncryptionKey_importBasic(t *testing.T) {
	testAccEBSKeyProvider := instaclustr.Provider()
	testAccEBSKeyProviders := map[string]terraform.ResourceProvider{
		"instaclustr": testAccEBSKeyProvider,
	}
	validConfig, _ := ioutil.ReadFile("data/valid_encryption_key.tf")
	username := os.Getenv("IC_USERNAME")
	apiKey := os.Getenv("IC_API_KEY")
	hostname := getOptionalEnv("IC_API_URL", instaclustr.DefaultApiHostname)
	kmsArn := os.Getenv("KMS_ARN")
	oriConfig := fmt.Sprintf(string(validConfig), username, apiKey, hostname, kmsArn)
	resource.Test(t, resource.TestCase{
		Providers:    testAccEBSKeyProviders,
		PreCheck:     func() { AccTestEnvVarsCheck(t) },
		CheckDestroy: testCheckAccEBSResourceDeleted("valid", hostname, username, apiKey),
		Steps: []resource.TestStep{
			{
				Config: oriConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAccEBSResourceValid("valid"),
					testCheckAccEBSResourceCreated("valid", hostname, username, apiKey),
				),
			},
			{
				Config:            oriConfig,
				ResourceName:      "instaclustr_encryption_key.valid",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFirewallRule_importBasic(t *testing.T) {
	testProviders := map[string]terraform.ResourceProvider{
		"instaclustr": instaclustr.Provider(),
	}
	tfFile, _ := ioutil.ReadFile("data/valid_with_firewall_rule.tf")
	username := os.Getenv("IC_USERNAME")
	apiKey := os.Getenv("IC_API_KEY")
	hostname := getOptionalEnv("IC_API_URL", instaclustr.DefaultApiHostname)
	config := fmt.Sprintf(string(tfFile), username, apiKey, hostname)

	resource.Test(t, resource.TestCase{
		Providers:    testProviders,
		CheckDestroy: checkFirewallRuleDeleted(hostname, username, apiKey),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					checkFirewallRuleState,
					checkFirewallRuleCreated(hostname, username, apiKey),
				),
			},
			{
				Config:            config,
				ResourceName:      "instaclustr_firewall_rule.valid_with_firewall_rule",
				ImportState:       true,
				ImportStateIdFunc: testAccFirewallRuleImportStateIdFunc("instaclustr_firewall_rule.valid_with_firewall_rule"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFirewallRuleImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s&%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["rule_cidr"]), nil
	}
}

func TestKafkaUserResource_importBasic(t *testing.T) {
	testProviders := map[string]terraform.ResourceProvider{
		"instaclustr": instaclustr.Provider(),
	}

	configBytes1, _ := ioutil.ReadFile("data/kafka_user_create_cluster.tf")
	username := os.Getenv("IC_USERNAME")
	apiKey := os.Getenv("IC_API_KEY")
	hostname := getOptionalEnv("IC_API_URL", instaclustr.DefaultApiHostname)

	zookeeperNodeSize := "zk-developer-t3.small-20"

	createClusterConfig := fmt.Sprintf(string(configBytes1), username, apiKey, hostname, zookeeperNodeSize)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: createClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceValidKafka("instaclustr_cluster.kafka_cluster"),
					checkKafkaClusterRunning(hostname, username, apiKey),
				),
			},
			{
				Config:            createClusterConfig,
				ResourceName:      "instaclustr_cluster.kafka_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
