package db

import (
	"nero/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateOrderTable() {
	tableName := "Orders"
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("user_id"),
				AttributeType: aws.String("S"),
			},
			// {
			// 	AttributeName: aws.String("phone"),
			// 	AttributeType: aws.String("N"),
			// },
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("user_id"),
				KeyType:       aws.String("HASH"), // HASH denotes the partition key
			},
			// {
			// 	AttributeName: aws.String("phone"),
			// 	KeyType:       aws.String("RANGE"),
			// },
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}
	_, err := db.CreateTable(input)
	if err != nil {
		utils.LogError("db/order_table.go", err, "line-40, func-CreateOrderTable")
		return
	}
	utils.ColoredPrintln("Table "+tableName+" created successfully", utils.CYellow)

}

func DeleteOrderTable() {
	del := &dynamodb.DeleteTableInput{
		TableName: aws.String("Orders"),
	}
	db.DeleteTable(del)
	utils.ColoredPrintln("Orders table deleted", utils.CRed)
}
