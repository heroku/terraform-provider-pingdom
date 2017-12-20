package pingdom

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heroku/go-pingdom/pingdom"
)

func resourcePingdomCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomCheckCreate,
		Read:   resourcePingdomCheckRead,
		Update: resourcePingdomCheckUpdate,
		Delete: resourcePingdomCheckDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePingdomCheckImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resolution": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"sendnotificationwhendown": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"notifyagainevery": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},

			"notifywhenbackup": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"userids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"integrationids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"encryption": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},

			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "/",
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"shouldcontain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"shouldnotcontain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"postdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"requestheaders": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"probefilters": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func checkForResource(d *schema.ResourceData) (pingdom.Check, error) {
	b := pingdom.BaseCheck{}

	// required
	if v, ok := d.GetOk("name"); ok {
		b.Name = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		b.Host = v.(string)
	}

	// optional
	if v, ok := d.GetOk("resolution"); ok {
		b.Resolution = pingdom.OptInt(v.(int))
	}

	if v, ok := d.GetOk("sendnotificationwhendown"); ok {
		b.SendNotificationWhenDown = pingdom.OptInt(v.(int))
	}

	if v, ok := d.GetOk("notifyagainevery"); ok {
		b.NotifyAgainEvery = pingdom.OptInt(v.(int))
	}

	if v, ok := d.GetOk("notifywhenbackup"); ok {
		b.NotifyWhenBackup = pingdom.OptBool(v.(bool))
	}

	if v, ok := d.GetOk("userids"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var intSlice []int
		for i := range interfaceSlice {
			intSlice = append(intSlice, interfaceSlice[i].(int))
		}
		b.UserIds = &intSlice
	}

	if v, ok := d.GetOk("integrationids"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var intSlice []int
		for i := range interfaceSlice {
			intSlice = append(intSlice, interfaceSlice[i].(int))
		}
		b.IntegrationIds = &intSlice
	}

	if v, ok := d.GetOk("tags"); ok {
		b.Tags = pingdom.OptStr(v.(string))
	}

	if v, ok := d.GetOk("probefilters"); ok {
		b.ProbeFilters = pingdom.OptStr(v.(string))
	}

	checkType := d.Get("type")
	switch checkType {
	case "http":
		h := pingdom.HttpCheck{BaseCheck: b}

		if v, ok := d.GetOk("url"); ok {
			h.Url = pingdom.OptStr(v.(string))
		}

		if v, ok := d.GetOk("encryption"); ok {
			h.Encryption = pingdom.OptBool(v.(bool))
		}

		if v, ok := d.GetOk("port"); ok {
			h.Port = pingdom.OptInt(v.(int))
		}

		if v, ok := d.GetOk("username"); ok {
			h.Username = pingdom.OptStr(v.(string))
		}

		if v, ok := d.GetOk("password"); ok {
			h.Password = pingdom.OptStr(v.(string))
		}

		if v, ok := d.GetOk("shouldcontain"); ok {
			h.ShouldContain = pingdom.OptStr(v.(string))
		}

		if v, ok := d.GetOk("shouldnotcontain"); ok {
			h.ShouldNotContain = pingdom.OptStr(v.(string))
		}

		if v, ok := d.GetOk("postdata"); ok {
			h.PostData = pingdom.OptStr(v.(string))
		}

		if m, ok := d.GetOk("requestheaders"); ok {
			h.RequestHeaders = make(map[string]string)
			for k, v := range m.(map[string]interface{}) {
				h.RequestHeaders[k] = v.(string)
			}
		}
		return &h, nil
	case "ping":
		return pingdom.PingCheck{BaseCheck: b}, nil
	default:
		return nil, fmt.Errorf("unknown type for check '%v'", checkType)
	}
}

func resourcePingdomCheckCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	check, err := checkForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Check create configuration: %#v, %#v", d.Get("name"), d.Get("hostname"))

	ck, err := client.Checks.Create(check)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(ck.ID))

	return nil
}

func resourcePingdomCheckRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	cl, err := client.Checks.List()
	if err != nil {
		return fmt.Errorf("Error retrieving list of checks: %s", err)
	}
	exists := false
	for _, ckid := range cl {
		if ckid.ID == id {
			exists = true
			break
		}
	}
	if !exists {
		d.SetId("")
		return nil
	}
	ck, err := client.Checks.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving check: %s", err)
	}

	d.Set("host", ck.Hostname)
	d.Set("name", ck.Name)
	d.Set("resolution", ck.Resolution)
	d.Set("sendnotificationwhendown", ck.SendNotificationWhenDown)
	d.Set("notifyagainevery", ck.NotifyAgainEvery)
	d.Set("notifywhenbackup", ck.NotifyWhenBackup)
	cids := schema.NewSet(
		func(contactId interface{}) int { return contactId.(int) },
		[]interface{}{},
	)
	for _, contactId := range ck.ContactIds {
		cids.Add(contactId)
	}
	d.Set("userids", cids)

	integids := schema.NewSet(
		func(integrationId interface{}) int { return integrationId.(int) },
		[]interface{}{},
	)
	for _, integrationId := range ck.IntegrationIds {
		integids.Add(integrationId)
	}
	d.Set("integrationids", integids)

	if ck.Type.HTTP == nil {
		ck.Type.HTTP = &pingdom.CheckResponseHTTPDetails{}
	}
	d.Set("url", ck.Type.HTTP.Url)
	d.Set("encryption", ck.Type.HTTP.Encryption)
	d.Set("port", ck.Type.HTTP.Port)
	d.Set("username", ck.Type.HTTP.Username)
	d.Set("password", ck.Type.HTTP.Password)
	d.Set("shouldcontain", ck.Type.HTTP.ShouldContain)
	d.Set("shouldnotcontain", ck.Type.HTTP.ShouldNotContain)
	d.Set("postdata", ck.Type.HTTP.PostData)

	if v, ok := ck.Type.HTTP.RequestHeaders["User-Agent"]; ok {
		if strings.HasPrefix(v, "Pingdom.com_bot_version_") {
			delete(ck.Type.HTTP.RequestHeaders, "User-Agent")
		}
	}
	d.Set("requestheaders", ck.Type.HTTP.RequestHeaders)

	return nil
}

func resourcePingdomCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	check, err := checkForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Check update configuration: %#v, %#v", d.Get("name"), d.Get("hostname"))

	_, err = client.Checks.Update(id, check)
	if err != nil {
		return fmt.Errorf("Error updating check: %s", err)
	}

	return nil
}

func resourcePingdomCheckDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[INFO] Deleting Check: %v", id)

	_, err = client.Checks.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting check: %s", err)
	}

	return nil
}

func resourcePingdomCheckImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return nil, fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[INFO] Importing key using ADDR ID %s", id)

	return []*schema.ResourceData{d}, nil
}
