package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceAwsDxConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsDxConnectionRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jumbo_frame_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_logical_redundancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_device": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"tags": tagsSchemaComputed(),
		},
	}
}

func dataSourceAwsDxConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).dxconn

	req := &directconnect.DescribeConnectionsInput{}

	if v, ok := d.GetOk("id"); ok {
		req.ConnectionId = aws.String(v.(string))
	}

	resp, err := conn.DescribeConnections(req)

	if err != nil {
		return err
	}

	connections := resp.Connections[:]

	if v, ok := d.GetOk("name"); ok {
		log.Printf("[DEBUG] Filtering DX Connections by Name: %s", v)
		nameFiltered := make([]*directconnect.Connection, 0)
		for _, c := range connections {
			if *(c.ConnectionName) == v.(string) {
				nameFiltered = append(nameFiltered, c)
			}
		}
		connections = nameFiltered
	}

	if len(connections) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}
	if len(connections) > 1 {
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria.")
	}

	connection := connections[0]

	d.SetId(aws.StringValue(connection.ConnectionId))

	arn := arn.ARN{
		Partition: meta.(*AWSClient).partition,
		Region:    meta.(*AWSClient).region,
		Service:   "directconnect",
		AccountID: meta.(*AWSClient).accountid,
		Resource:  fmt.Sprintf("dxcon/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("state", connection.ConnectionState)
	d.Set("location", connection.Location)
	d.Set("bandwidth", connection.Bandwidth)
	d.Set("jumbo_frame_capable", connection.JumboFrameCapable)
	d.Set("has_logical_redundancy", connection.HasLogicalRedundancy)
	d.Set("name", connection.ConnectionName)
	d.Set("aws_device", connection.AwsDeviceV2)
	d.Set("vlan", connection.Vlan)
	return nil
}
