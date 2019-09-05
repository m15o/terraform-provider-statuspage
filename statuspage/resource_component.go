package statuspage

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	sp "github.com/yannh/statuspage-go-sdk"
)

func resourceComponentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)
	component, err := sp.CreateComponent(
		client, d.Get("page_id").(string),
		&sp.Component{
			Name:               d.Get("name").(string),
			Description:        d.Get("description").(string),
			OnlyShowIfDegraded: d.Get("only_show_if_degraded").(bool),
			Status:             d.Get("status").(string),
			Showcase:           d.Get("showcase").(bool),
		},
	)
	if err != nil {
		log.Printf("[WARN] Statuspage Failed creating component: %s\n", err)
		return err
	}

	log.Printf("[INFO] Statuspage Created: %s\n", component.ID)
	d.SetId(component.ID)

	return resourceComponentRead(d, m)
}

func resourceComponentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)
	component, err := sp.GetComponent(client, d.Get("page_id").(string), d.Id())
	if err != nil {
		log.Printf("[ERROR] Statuspage could not find component with ID: %s\n", d.Id())
		return err
	}

	if component == nil {
		log.Printf("[INFO] Statuspage could not find component with ID: %s\n", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Statuspage read: %s\n", component.ID)

	d.Set("name", component.Name)
	d.Set("description", component.Description)
	d.Set("group_id", component.GroupID)
	d.Set("only_show_if_degraded", component.OnlyShowIfDegraded)
	d.Set("status", component.Status)
	d.Set("showcase", component.Showcase)
	d.Set("automation_email", component.AutomationEmail)

	return nil
}

func resourceComponentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)
	componentID := d.Id()

	_, err := sp.UpdateComponent(
		client,
		d.Get("page_id").(string),
		componentID,
		&sp.Component{
			Name:               d.Get("name").(string),
			Description:        d.Get("description").(string),
			OnlyShowIfDegraded: d.Get("only_show_if_degraded").(bool),
			Status:             d.Get("status").(string),
			Showcase:           d.Get("showcase").(bool),
		},
	)
	if err != nil {
		log.Printf("[WARN] Statuspage Failed creating component: %s\n", err)
		return err
	}

	d.SetId(componentID)

	return resourceComponentRead(d, m)
}

func resourceComponentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sp.Client)

	return sp.DeleteComponent(client, d.Get("page_id").(string), d.Id())
}

func resourceComponentImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	resourceId, err := parsePageResourceId(d.Id())
	if err != nil {
		return nil, errors.New("id is not formatted properly; id should be '$page_id/$component_id', but: " + d.Id())
	}
	client := m.(*sp.Client)

	component, err := sp.GetComponent(client, resourceId.pageId, resourceId.resourceId)
	if err != nil {
		log.Printf("[ERROR] Statuspage could not find component with ID: %s\n", d.Id())
		return nil, err
	}

	if component == nil {
		log.Printf("[ERROR] Statuspage returns null component with ID: %s\n", d.Id())
		return nil, errors.New("Statuspage could not find component with ID: " + d.Id())
	}

	d.SetId(component.ID)
	d.Set("page_id", component.PageID)
	d.Set("name", component.Name)
	d.Set("description", component.Description)
	d.Set("group_id", component.GroupID)
	d.Set("only_show_if_degraded", component.OnlyShowIfDegraded)
	d.Set("status", component.Status)
	d.Set("showcase", component.Showcase)
	d.Set("automation_email", component.AutomationEmail)

	log.Printf("[INFO] Statuspage imported component: %s\n", component.ID)
	return []*schema.ResourceData{d}, nil
}

func resourceComponent() *schema.Resource {
	return &schema.Resource{
		Create: resourceComponentCreate,
		Read:   resourceComponentRead,
		Update: resourceComponentUpdate,
		Delete: resourceComponentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComponentImport,
		},
		Schema: map[string]*schema.Schema{
			"page_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "the ID of the page this component belongs to",
				Required:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Display Name for the component",
				Required:    true,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "More detailed description for the component",
				Optional:    true,
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Status of component",
				Optional:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{"operational", "under_maintenance", "degraded_performance", "partial_outage", "major_outage", ""},
					false,
				),
				Default: "operational",
			},
			"only_show_if_degraded": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Should this component be shown component only if in degraded state",
				Optional:    true,
			},
			"showcase": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Should this component be showcased",
				Optional:    true,
				Default:     true,
			},
			"automation_email": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Email address to send automation events to",
				Computed:    true,
			},
		},
	}
}
