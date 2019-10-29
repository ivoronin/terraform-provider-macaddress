# terraform-provider-macaddress
Generates random locally administered unicast MAC address

# Use case
```hcl
resource "macaddress" "example_address" {
}

// Terraform Mikrotik Provider - https://github.com/ddelnano/terraform-provider-mikrotik
resource "mikrotik_dhcp_lease" "example_lease" {
  address    = "10.0.0.10"
  macaddress = upper(macaddress.example_address.address)
  comment    = "Example DHCP Lease"
}
```
