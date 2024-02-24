package models

import (
	"nero/db"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Info struct {
	ProductId   string `json:"product_id" validate:"required"`
	ProductName string `json:"product_name" validate:"required,min=3"`
	Quantity    int    `json:"quantity" validate:"numeric,gt=0"`
	Price       int    `json:"price" validate:"required,numeric,gt=1"`
}
type Order struct {
	ID        string    `json:"order_id"`
	UserId    string    `json:"user_id" validate:"required"`
	Total     int       `json:"total" validate:"required,numeric"`
	Tax       float32   `json:"tax" validate:"required"`
	Currency  string    `json:"currency" validate:"required"`
	Info      *[]Info   `json:"info"`
	CreatedAt time.Time `json:"createdAt"`
	updatedAt time.Time `json:"-"`
}

func (o *Order) CreateOrder() (*Order, error) {
	db := db.GetDB()

	o.ID = uuid.New().String()
	o.CreatedAt = time.Now()
	o.updatedAt = time.Now()

	item, err := dynamodbattribute.MarshalMap(o)
	if err != nil {
		return nil, err
	}

	param := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("Orders"),
	}
	_, err = db.PutItem(param)
	if err != nil {
		return nil, err
	}
	return o, nil
}
func (o *Order) GetById(id string) (*Order, error) {
	db := db.GetDB()
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"order_id": {
				S: aws.String(id),
			},
		},
		TableName:      aws.String("Orders"),
		ConsistentRead: aws.Bool(true),
	}

	resp, err := db.GetItem(params)
	if err != nil {
		return nil, err
	}

	var order *Order
	if err := dynamodbattribute.UnmarshalMap(resp.Item, &order); err != nil {
		return nil, err
	}
	return order, nil
}
func (o *Order) Validate() error {
	validate := validator.New()
	return validate.Struct(o)
}
