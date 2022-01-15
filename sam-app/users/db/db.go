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

	//TODO: Googleのuuidに変える
	"github.com/teris-io/shortid"
)

var (
	region    = os.Getenv("REGION")
	userTable = os.Getenv("USERS_TABLE_NAME")
)

type UserRepository struct {
	Instance *dynamodb.DynamoDB
	Context  context.Context
}

type User struct {
	ID   string `json:"user_id"`
	Name string `json:"user_name"`
}

func New() UserRepository {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)}),
	)
	dynamo := dynamodb.New(sess)
	xray.AWS(dynamo.Client)
	return UserRepository{Instance: dynamo}
}

func (d UserRepository) GetUser(userId string) (*User, error) {
	return d.doGetUser(userId, d.Context)
}

func (d UserRepository) doGetUser(userId string, ctx context.Context) (*User, error) {
	//Itemの取得（X-Rayトレース）
	result, err := d.Instance.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(userTable),
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

func (d UserRepository) PutUser(user *User) (*User, error) {
	return d.doPutUser(user, d.Context)
}

func (d UserRepository) doPutUser(user *User, ctx context.Context) (*User, error) {
	//ID採番
	userId := shortid.MustGenerate()
	user.ID = userId

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(userTable),
	}
	//Itemの登録（X-Rayトレース）
	_, err = d.Instance.PutItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return user, nil
}
