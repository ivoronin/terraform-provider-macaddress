package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccMacAddressConfig = `resource "macaddress" "address_test" {}`
)

var testAccMacAddressPatten = regexp.MustCompile(`^([0-9a-f]{2}):([0-9a-f]{2}):([0-9a-f]{2}):([0-9a-f]{2}):([0-9a-f]{2}):([0-9a-f]{2})`)

func TestAccMacAddress(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMacAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccMacAddressCheck("macaddress.address_test"),
				),
			},
		},
	})
}

func testAccMacAddressCheck(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var i uint8

		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		address := rs.Primary.Attributes["address"]

		matches := testAccMacAddressPatten.FindAllString(address, -1)
		if matches == nil {
			return fmt.Errorf("result address (%s) format is incorrect", address)
		}

		fmt.Sscanf(matches[0], "%x", &i)

		mcastBit := hasBit(i, 0)
		localBit := hasBit(i, 1)

		if mcastBit {
			return fmt.Errorf("result address (%s) is incorrect - mcast bit is set", address)
		}

		if !localBit {
			return fmt.Errorf("result address (%s) is incorrect - local bit is not set", address)
		}

		return nil
	}
}

func hasBit(n uint8, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func testAccPreCheck(t *testing.T) {
}
