package driven_infra_aws

import (
	"app/pkg/crosscutting"
	"app/pkg/crosscutting/errors"
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient interface {
	GetItem(ctx context.Context, input dynamodb.GetItemInput) (dynamodb.GetItemOutput, error)
	Query(ctx context.Context, input dynamodb.QueryInput) (dynamodb.QueryOutput, error)
	PutItem(ctx context.Context, input dynamodb.PutItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.PutItemOutput, error)
	UpdateItem(ctx context.Context, input dynamodb.UpdateItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, input dynamodb.DeleteItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.DeleteItemOutput, error)
	IssueSequentialID(ctx context.Context, key string, ts int64) (int64, error)
}

func NewDynamoDBClient(ctx context.Context) (DynamoDBClient, error) {
	cnf, err := NewAWSConfig(ctx)
	if err != nil {
		return nil, errors.LiftWithCtx(err, ctx)
	}
	return &dynamoDBClient{
		dynamo: dynamodb.NewFromConfig(cnf),
	}, nil
}

type dynamoDBClient struct {
	dynamo *dynamodb.Client
}

func (d *dynamoDBClient) GetItem(ctx context.Context, input dynamodb.GetItemInput) (dynamodb.GetItemOutput, error) {
	output, err := d.dynamo.GetItem(ctx, &input)
	if err != nil {
		return dynamodb.GetItemOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.GetItemOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) Query(ctx context.Context, input dynamodb.QueryInput) (dynamodb.QueryOutput, error) {
	output, err := d.dynamo.Query(ctx, &input)
	if err != nil {
		return dynamodb.QueryOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.QueryOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) PutItem(ctx context.Context, input dynamodb.PutItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.PutItemOutput, error) {
	output, err := d.dynamo.PutItem(ctx, &input, optFns...)
	if err != nil {
		return dynamodb.PutItemOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.PutItemOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) UpdateItem(ctx context.Context, input dynamodb.UpdateItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.UpdateItemOutput, error) {
	output, err := d.dynamo.UpdateItem(ctx, &input, optFns...)
	if err != nil {
		return dynamodb.UpdateItemOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.UpdateItemOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) DeleteItem(ctx context.Context, input dynamodb.DeleteItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.DeleteItemOutput, error) {
	output, err := d.dynamo.DeleteItem(ctx, &input, optFns...)
	if err != nil {
		return dynamodb.DeleteItemOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.DeleteItemOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) IssueSequentialID(ctx context.Context, key string, ts int64) (int64, error) {
	_, err := d.getSequentialIDItem(ctx, key)
	isNotFound := errors.IsDataNotFoundError(err)
	if err != nil && isNotFound == false {
		return 0, errors.LiftWithCtx(err, ctx)
	}
	if isNotFound {
		err = d.initSequentialIDKey(ctx, key, ts)
		if err != nil {
			return 0, errors.LiftWithCtx(err, ctx)
		}
	}
	newID, err := d.countUpSequentialID(ctx, key, 1, ts)
	if err != nil {
		return 0, errors.LiftWithCtx(err, ctx)
	}
	return newID, nil
}

func (d *dynamoDBClient) getSequentialIDItem(ctx context.Context, key string) (sequentialIDItem, error) {
	item := sequentialIDItem{}
	input := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: key},
		},
		TableName:       crosscutting.Ptr(TableSequentialID),
		AttributesToGet: nil,
		ConsistentRead:  crosscutting.Ptr(true),
	}
	getRes, err := d.GetItem(ctx, input)
	if err != nil {
		return item, errors.LiftWithCtx(err, ctx)
	}

	err = item.SetByDynamoDBAttrs(getRes.Item)
	if err != nil {
		return item, errors.LiftWithCtx(err, ctx)
	}
	return item, nil
}

func (d *dynamoDBClient) initSequentialIDKey(ctx context.Context, key string, ts int64) error {
	tsStr := strconv.FormatInt(ts, 10)
	input := dynamodb.PutItemInput{
		TableName: crosscutting.Ptr(TableSequentialID),
		Item: map[string]types.AttributeValue{
			"id":         &types.AttributeValueMemberS{Value: key},
			"current_id": &types.AttributeValueMemberN{Value: "0"},
			"created_at": &types.AttributeValueMemberN{Value: tsStr},
			"updated_at": &types.AttributeValueMemberN{Value: tsStr},
		},
		ConditionExpression: crosscutting.Ptr("attribute_not_exists(id)"),
	}
	_, err := d.PutItem(ctx, input)
	if err != nil {
		return errors.LiftWithCtx(err, ctx)
	}
	return nil
}

func (d *dynamoDBClient) countUpSequentialID(ctx context.Context, key string, cnt int, ts int64) (int64, error) {
	tsStr := strconv.FormatInt(ts, 10)
	updateReq := dynamodb.UpdateItemInput{
		TableName: crosscutting.Ptr(TableSequentialID),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: key},
		},
		UpdateExpression: crosscutting.Ptr("ADD current_id :incr, updated_at = :updated_at"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":incr":       &types.AttributeValueMemberN{Value: strconv.Itoa(cnt)},
			":updated_at": &types.AttributeValueMemberN{Value: tsStr},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}
	output, err := d.UpdateItem(ctx, updateReq)
	if err != nil {
		return 0, errors.LiftWithCtx(err, ctx)
	}
	updatedItem := sequentialIDItem{}
	err = updatedItem.SetByDynamoDBAttrs(output.Attributes)
	if err != nil {
		return 0, errors.LiftWithCtx(err, ctx)
	}
	return updatedItem.CurrentID, nil
}
