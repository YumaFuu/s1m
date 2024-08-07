package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Update an SSM parameter
func (c Client) Put(
	name *string,
	ptype ParameterType,
	newValue *string,
) error {
	ctx := context.TODO()

	_, err := c.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      name,
		Type:      ptype,
		Value:     aws.String(*newValue),
		Overwrite: aws.Bool(true),
	})
	return err
}
