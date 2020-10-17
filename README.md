# terraform-provider-macaddress
Generates random locally administered unicast MAC address

Terraform Registry: https://registry.terraform.io/providers/ivoronin/macaddress/latest

Docs: https://registry.terraform.io/providers/ivoronin/macaddress/latest/docs

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
