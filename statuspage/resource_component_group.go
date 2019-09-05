package statuspage

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	sp "github.com/yannh/statuspage-go-sdk"
)

func resourceComponentGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)

	tfComponents := d.Get("components").(*schema.Set).List()
	components := make([]string, len(tfComponents))
	for i, tfComponent := range tfComponents {
		components[i] = tfComponent.(string)
	}

	componentGroup, err := sp.CreateComponentGroup(
		client,
		d.Get("page_id").(string),
		&sp.ComponentGroup{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Components:  components,
		},
	)

	if err != nil {
		log.Printf("[WARN] Statuspage Failed creating component group: %s\n", err)
		return err
	}

	log.Printf("[INFO] Statuspage Created component group: %s\n", componentGroup.ID)
	d.SetId(componentGroup.ID)

	return resourceComponentGroupRead(d, m)
}

func resourceComponentGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)
	componentGroup, err := sp.GetComponentGroup(client, d.Get("page_id").(string), d.Id())
	if err != nil {
		log.Printf("[ERROR] Statuspage could not find component group with ID: %s\n", d.Id())
		return err
	}

	if componentGroup == nil {
		log.Printf("[INFO] Statuspage could not find component group with ID: %s\n", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Statuspage read component group: %s\n", componentGroup.ID)

	d.Set("name", componentGroup.Name)
	d.Set("description", componentGroup.Description)
	d.Set("components", componentGroup.Components)

	return nil
}

func resourceComponentGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)
	componentGroupID := d.Id()

	tfComponents := d.Get("components").(*schema.Set).List()
	components := make([]string, len(tfComponents))
	for i, tfComponent := range tfComponents {
		components[i] = tfComponent.(string)
	}

	_, err := sp.UpdateComponentGroup(
		client,
		d.Get("page_id").(string),
		componentGroupID,
		&sp.ComponentGroup{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Components:  components,
		},
	)
	if err != nil {
		log.Printf("[WARN] Statuspage Failed creating component group: %s\n", err)
		return err
	}

	d.SetId(componentGroupID)

	return resourceComponentGroupRead(d, m)
}

func resourceComponentGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)

	return sp.DeleteComponentGroup(client, d.Get("page_id").(string), d.Id())
}

func resourceComponentGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	resourceId, err := parsePageResourceId(d.Id())
	if err != nil {
		return nil, errors.New("id is not formatted properly; id should be '$page_id/$component_id', but: " + d.Id())
	}
	client := m.(*sp.Client)

	componentGroup, err := sp.GetComponentGroup(client, resourceId.pageId, resourceId.resourceId)
	if err != nil {
		log.Printf("[ERROR] Statuspage could not find component group with ID: %s\n", d.Id())
		return nil, err
	}

	if componentGroup == nil {
		log.Printf("[ERROR] Statuspage could returns null component group with ID: %s\n", d.Id())
		return nil, errors.New("Statuspage could not find component with ID: " + d.Id())
	}

	d.SetId(componentGroup.ID)
	d.Set("page_id", componentGroup.PageID)
	d.Set("name", componentGroup.Name)
	d.Set("description", componentGroup.Description)
	d.Set("components", componentGroup.Components)

	log.Printf("[INFO] Statuspage imported componentGroup: %s\n", componentGroup.ID)
	return []*schema.ResourceData{d}, nil
}

func resourceComponentGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceComponentGroupCreate,
		Read:   resourceComponentGroupRead,
		Update: resourceComponentGroupUpdate,
		Delete: resourceComponentGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComponentGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"page_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "the ID of the page this component group belongs to",
				Required:    true,
			},
			"components": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "An array with the IDs of the components in this group",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Required:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Display name for this component group",
				Required:    true,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "More detailed description for this component group",
				Optional:    true,
			},
		},
	}
}
