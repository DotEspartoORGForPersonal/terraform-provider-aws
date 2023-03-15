package lightsail

import (
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_lightsail_lb_attachment")
func ResourceLoadBalancerAttachment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLoadBalancerAttachmentCreate,
		ReadWithoutTimeout:   resourceLoadBalancerAttachmentRead,
		DeleteWithoutTimeout: resourceLoadBalancerAttachmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"lb_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(2, 255),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]`), "must begin with an alphabetic character"),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_\-.]+[^._\-]$`), "must contain only alphanumeric characters, underscores, hyphens, and dots"),
				),
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLoadBalancerAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailConn()
	lbName := d.Get("lb_name").(string)
	req := lightsail.AttachInstancesToLoadBalancerInput{
		LoadBalancerName: aws.String(lbName),
		InstanceNames:    aws.StringSlice([]string{d.Get("instance_name").(string)}),
	}

	out, err := conn.AttachInstancesToLoadBalancerWithContext(ctx, &req)

	if err != nil {
		return create.DiagError(names.Lightsail, lightsail.OperationTypeAttachInstancesToLoadBalancer, ResLoadBalancerAttachment, lbName, err)
	}

	diag := expandOperations(ctx, conn, out.Operations, lightsail.OperationTypeAttachInstancesToLoadBalancer, ResLoadBalancerAttachment, lbName)

	if diag != nil {
		return diag
	}

	// Generate an ID
	vars := []string{
		d.Get("lb_name").(string),
		d.Get("instance_name").(string),
	}

	d.SetId(strings.Join(vars, ","))

	return resourceLoadBalancerAttachmentRead(ctx, d, meta)
}

func resourceLoadBalancerAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailConn()

	out, err := FindLoadBalancerAttachmentById(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		create.LogNotFoundRemoveState(names.Lightsail, create.ErrActionReading, ResLoadBalancerAttachment, d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.Lightsail, create.ErrActionReading, ResLoadBalancerAttachment, d.Id(), err)
	}

	d.Set("instance_name", out)
	d.Set("lb_name", expandLoadBalancerNameFromId(d.Id()))

	return nil
}

func resourceLoadBalancerAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailConn()

	id_parts := strings.SplitN(d.Id(), ",", -1)
	if len(id_parts) != 2 {
		return nil
	}

	lbName := id_parts[0]
	iName := id_parts[1]

	in := lightsail.DetachInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(lbName),
		InstanceNames:    aws.StringSlice([]string{iName}),
	}

	out, err := conn.DetachInstancesFromLoadBalancerWithContext(ctx, &in)

	if err != nil {
		return create.DiagError(names.Lightsail, lightsail.OperationTypeDetachInstancesFromLoadBalancer, ResLoadBalancerAttachment, lbName, err)
	}

	diag := expandOperations(ctx, conn, out.Operations, lightsail.OperationTypeAttachInstancesToLoadBalancer, ResLoadBalancerAttachment, lbName)

	if diag != nil {
		return diag
	}

	return nil
}

func expandLoadBalancerNameFromId(id string) string {
	id_parts := strings.SplitN(id, ",", -1)
	lbName := id_parts[0]

	return lbName
}
