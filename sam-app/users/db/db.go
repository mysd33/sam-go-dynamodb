package db

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

var (
	Region = os.Getenv("REGION")
)

type DB struct {
	Instance *dynamodb.DynamoDB
	//TODO: context使っていない？
	ctx context.Context
}

type User struct {
	ID   string `json:"user_id"`
	Name string `json:"user_name"`
}

func New() DB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(Region)}),
	)
	dynamo := dynamodb.New(sess)
	xray.AWS(dynamo.Client)
	return DB{Instance: dynamo}
}

func (d DB) GetUser(userId string, ctx context.Context) (*User, error) {
	//Itemの取得（X-Rayトレース）
	result, err := d.Instance.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		//TODO: テーブル名切り出し
		TableName: aws.String("users"),
		//TODO: map[string]*の意味わからず
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(userId),
			},
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get item")
	}
	if result.Item == nil {
		return nil, nil
	}
	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal item")
	}
	return &user, nil
}

func (d DB) PutUser(user *User, ctx context.Context) (*User, error) {
	//ID採番
	userId := shortid.MustGenerate()
	user.ID = userId

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}
	//TODO: テーブル名切り出し
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("users"),
	}
	//Itemの登録（X-Rayトレース）
	_, err = d.Instance.PutItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return user, nil
}
