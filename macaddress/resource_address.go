package macaddress

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MAC_ADDRESS_LENGTH = 6

func resourceAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAddressCreate,
		ReadContext:   resourceAddressNoop,
		DeleteContext: resourceAddressNoop,
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: resourceAddressImport,
		},
	}
}

func resourceAddressImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	address := d.Id()
	parts := strings.Split(address, ":")
	if len(parts) != 6 {
		return nil, fmt.Errorf("%s is not a valid mac address", address)
	}
	for _, p := range parts {
		_, err := strconv.ParseInt(p, 16, 16)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid mac address", address)
		}
	}
	d.Set("address", address)
	return []*schema.ResourceData{d}, nil
}

func resourceAddressCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var groups []string
	buf := make([]byte, MAC_ADDRESS_LENGTH)

	_, err := rand.Read(buf)
	if err != nil {
		return diag.FromErr(err)
	}

	// Locally administered
	buf[0] |= 0x02

	// Unicast
	buf[0] &= 0xfe

	prefix := d.Get("prefix").([]interface{})

	if len(prefix) > MAC_ADDRESS_LENGTH {
		return diag.FromErr(errors.New("error generating random mac address: prefix is too large"))
	}

	for index, val := range prefix {
		if val.(int) > 255 {
			return diag.FromErr(errors.New("error generating random mac address: prefix segment must be in the range [0,256)"))
		}
		buf[index] = byte(val.(int))
	}

	for _, i := range buf {
		groups = append(groups, fmt.Sprintf("%02x", i))
	}

	address := strings.Join(groups, ":")

	d.SetId(address)
	d.Set("address", address)

	return nil
}

func resourceAddressNoop(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
