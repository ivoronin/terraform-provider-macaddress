package macaddress

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"macaddress": resourceAddress(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
