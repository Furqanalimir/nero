package db

import (
	"nero/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateOrderTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("order_id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("order_id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Orders"),
	}
	result, err := db.CreateTable(input)
	if err != nil {
		utils.LogError("db/db.go", err, "line-40, func-Init")
		return
	}
	utils.ColoredPrintln(result.GoString(), utils.CYellow)
}

func DeleteOrderTable() {
	del := &dynamodb.DeleteTableInput{
		TableName: aws.String("Orders"),
	}
	db.DeleteTable(del)
	utils.ColoredPrintln("Orders table deleted", utils.CRed)
}
