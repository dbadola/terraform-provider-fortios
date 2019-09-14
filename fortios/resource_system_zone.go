package fortios

import (
	"fmt"
	"log"

	forticlient "github.com/fgtdev/fortios-sdk-go/sdkcore"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSystemZone() *schema.Resource {

	return &schema.Resource{
		Create: resourceSystemZoneCreateUpdate,
		Read:   resourceSystemZoneRead,
		Update: resourceSystemZoneCreateUpdate,
		Delete: resourceSystemZoneDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"intrazone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "allow",
			},
			"interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

}
func resourceSystemZoneCreateUpdate(d *schema.ResourceData, m interface{}) error {
	mkey := d.Id()

	c := m.(*FortiClient).Client
	c.Retries = 1
	//Get Params from d
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	intrazone := d.Get("intrazone").(string)
	interface1 := d.Get("interface").([]interface{})

	var interfaces []forticlient.ZoneMultiValue

	for _, v := range interface1 {
		if v == nil {
			return fmt.Errorf("null value")
		}
		interfaces = append(interfaces,
			forticlient.ZoneMultiValue{
				InterfaceName: v.(string),
			})
	}
	//Build input data by sdk
	i := &forticlient.JSONSystemZone{
		Name:        name,
		Description: description,
		Intrazone:   intrazone,
		Interface:   interfaces,
	}

	//Call process by sdk
	_, err := c.UpdateSystemZone(i, mkey)
	if err != nil {
		return fmt.Errorf("Error updating System Setting Zone: %s", err)
	}

	d.SetId(name)

	return resourceSystemZoneRead(d, m)

}

func resourceSystemZoneDelete(d *schema.ResourceData, m interface{}) error {
	// no API for this
	return nil
}
func resourceSystemZoneRead(d *schema.ResourceData, m interface{}) error {
	mkey := d.Id()

	c := m.(*FortiClient).Client
	c.Retries = 1

	//Call process by sdk
	o, err := c.ReadSystemZone(mkey)
	if err != nil {
		return fmt.Errorf("Error reading System Setting Zone: %s", err)
	}

	if o == nil {
		log.Printf("[WARN] resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	//Refresh property
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("intrazone", o.Intrazone)
	interface1 := forticlient.ExpandZone(o.Interface)
	if err := d.Set("interface", interface1); err != nil {
		log.Printf("[WARN] Error setting System Interface for (%s): %s", d.Id(), err)
	}

	return nil
}
