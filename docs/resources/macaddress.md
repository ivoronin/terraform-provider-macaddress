# Macaddress Resource

Generates random locally administered unicast MAC address or MAC address with specified prefix.

## Example Usage

```hcl
resource "macaddress" "example_local_address" {
}

resource "macaddress" "example_prefix_address" {
    // AA:00:04
    prefix = [170, 0, 4]
}
```

## Argument Reference

* `prefix` - (Optional) Address prefix as a list of decimal integers

## Attribute Reference

* `address` - Generated macaddress
