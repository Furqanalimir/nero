package db

import (
	"nero/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateUserTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("phone"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("phone"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Users"),
	}
	result, err := db.CreateTable(input)
	if err != nil {
		utils.LogError("db/db.go", err, "line-57, func-Init")
		return
	}
	utils.ColoredPrintln(result.GoString(), utils.CYellow)
}

func DeleteUserTable() {
	del := &dynamodb.DeleteTableInput{
		TableName: aws.String("Users"),
	}
	db.DeleteTable(del)
	utils.ColoredPrintln("User table deleted", utils.CRed)
}
