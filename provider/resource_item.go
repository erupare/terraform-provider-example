package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spaceapegames/terraform-provider-blog/api/client"
	"github.com/spaceapegames/terraform-provider-blog/api/server"
	"strings"
)

func resourceItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The name of the resource, also acts as it's unique ID",
				ForceNew:         true,
			},
			"description": {
				Type: schema.TypeString,
				Required: true,
				Description: "A description of an item",
			},
			"tags": {
				Type: schema.TypeSet,
				Optional: true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
		},
		Create:             resourceCreateItem,
		Read:               resourceReadItem,
		Update:             resourceUpdateItem,
		Delete:             resourceDeleteItem,
		Exists:             resourceExistsItem,
		Importer:           &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateItem(d *schema.ResourceData, m interface{}) error  {
	apiClient := m.(*client.Client)

	tfTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(tfTags))
	for i, tfTag := range tfTags {
		tags[i] = tfTag.(string)
	}

	item := server.Item{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        tags,
	}

	err := apiClient.NewItem(&item)

	if err !=  nil {
		return err
	}
	d.SetId(item.Name)
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	itemId := d.Id()
	item, err := apiClient.GetItem(itemId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Item with ID %s", itemId)
		}
	}

	d.SetId(item.Name)
	d.Set("name", item.Name)
	d.Set("description", item.Description)
	d.Set("tags", item.Tags)
	return nil
}

func resourceUpdateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	tfTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(tfTags))
	for _, tfTag := range tfTags {
		tags = append(tags, tfTag.(string))
	}

	item := server.Item{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        tags,
	}

	err := apiClient.UpdateItem(&item)
	if err !=  nil {
		return err
	}
	return nil
}

func resourceDeleteItem(d *schema.ResourceData, m interface{}) error  {
	apiClient := m.(*client.Client)

	itemId := d.Id()

	err := apiClient.DeleteItem(itemId)
	if err != nil {
		return err
	}
	return nil
}

func resourceExistsItem(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	itemId := d.Id()
	_, err := apiClient.GetItem(itemId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}