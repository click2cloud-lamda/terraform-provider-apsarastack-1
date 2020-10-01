package apsarastack

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackDBDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDBDatabaseCreate,
		Read:   resourceApsaraStackDBDatabaseRead,
		Update: resourceApsaraStackDBDatabaseUpdate,
		Delete: resourceApsaraStackDBDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"character_set": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "utf8",
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	request := rds.CreateCreateDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = d.Get("instance_id").(string)
	request.DBName = d.Get("name").(string)
	request.CharacterSetName = d.Get("character_set").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.CreateDatabase(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBInstanceId, COLON_SEPARATED, request.DBName))

	return resourceApsaraStackDBDatabaseRead(d, meta)
}

func resourceApsaraStackDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	rsdService := RdsService{client}
	object, err := rsdService.DescribeDBDatabase(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.DBInstanceId)
	d.Set("name", object.DBName)
	d.Set("character_set", object.CharacterSetName)
	d.Set("description", object.DBDescription)

	return nil
}

func resourceApsaraStackDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	if d.HasChange("description") && !d.IsNewResource() {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		request := rds.CreateModifyDBDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBInstanceId = parts[0]
		request.DBName = parts[1]
		request.DBDescription = d.Get("description").(string)
		var raw interface{}
		raw, err = client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyDBDescription(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceApsaraStackDBDatabaseRead(d, meta)
}

func resourceApsaraStackDBDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := rds.CreateDeleteDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBInstanceId = parts[0]
	request.DBName = parts[1]
	// wait instance status is running before deleting database
	if err := rdsService.WaitForDBInstance(parts[0], Running, 1800); err != nil {
		return WrapError(err)
	}
	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DeleteDatabase(request)
	})
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(rdsService.WaitForDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
