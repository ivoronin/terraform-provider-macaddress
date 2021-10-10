package macaddress

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMacAddress(t *testing.T) {
	prefix := []byte{0x10, 0xfe, 0x55}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMacAddressConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacAddressValid("macaddress.address_basic"),
				),
			},
			{
				ResourceName:      "macaddress.address_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMacAddressConfigPrefix(prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacAddressPrefixMatch("macaddress.address_prefix", prefix),
				),
			},
		},
	})
}

func testAccMacAddressConfigBasic() string {
	return fmt.Sprintf(`
	resource "macaddress" "address_basic" {
	}
	`)
}

func testAccMacAddressConfigPrefix(prefix []byte) string {
	prefixStr := make([]string, len(prefix))
	for index, element := range prefix {
		prefixStr[index] = fmt.Sprintf("%d", element)
	}
	prefixListStr := strings.Join(prefixStr, ",")
	return fmt.Sprintf(`
	resource "macaddress" "address_prefix" {
        prefix = [%s]
	}
	`, prefixListStr)
}

func testAccCheckMacAddressValid(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		address := rs.Primary.Attributes["address"]

		octets, err := parseMacAddress(address)
		if err != nil {
			return err
		}

		mcastBit := hasBit(octets[0], 0)
		localBit := hasBit(octets[0], 1)

		if mcastBit {
			return fmt.Errorf("result address (%s) is incorrect - mcast bit is set", address)
		}

		if !localBit {
			return fmt.Errorf("result address (%s) is incorrect - local bit is not set", address)
		}

		return nil
	}
}

func testAccCheckMacAddressPrefixMatch(id string, prefix []byte) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		address := rs.Primary.Attributes["address"]

		octets, err := parseMacAddress(address)
		if err != nil {
			return err
		}

		for index, element := range prefix {
			if element != octets[index] {
				return fmt.Errorf("address (%s) does not match specified prefix", address)
			}
		}

		return nil
	}
}

func parseMacAddress(address string) ([MAC_ADDRESS_LENGTH]byte, error) {
	segments := strings.Split(address, ":")
	var octets [MAC_ADDRESS_LENGTH]byte

	if len(segments) != MAC_ADDRESS_LENGTH {
		return octets, fmt.Errorf("address (%s) is not valid", address)
	}

	for index, element := range segments {
		_, err := fmt.Sscanf(element, "%x", &octets[index])
		if err != nil {
			return octets, fmt.Errorf("element (%s) is not valid", element)
		}
	}

	return octets, nil
}

func hasBit(n byte, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}
